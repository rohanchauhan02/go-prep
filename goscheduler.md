# Go Runtime Internals — Expert Level Deep Dive

> For senior/expert engineers. Covers the GMP scheduler, garbage collector, work stealing, goroutine lifecycle, memory model, and runtime internals with code examples and diagrams.

---

## Table of Contents
1. [Thread vs Goroutine](#1-thread-vs-goroutine)
2. [GMP Model](#2-gmp-model)
3. [Goroutine Lifecycle](#3-goroutine-lifecycle)
4. [Work Stealing](#4-work-stealing)
5. [Preemption](#5-preemption)
6. [Goroutine Stack Management](#6-goroutine-stack-management)
7. [Garbage Collector](#7-garbage-collector)
8. [Memory Model & Happens-Before](#8-memory-model--happens-before)
9. [Runtime Introspection](#9-runtime-introspection)
10. [Common Pitfalls & Tuning](#10-common-pitfalls--tuning)

---

## 1. Thread vs Goroutine

| Property | OS Thread | Goroutine |
|----------|-----------|-----------|
| Stack size | 1–8 MB (fixed) | 2–8 KB (grows dynamically) |
| Creation cost | ~1ms, kernel syscall | ~0.3µs, user-space |
| Context switch | ~1–10µs (kernel mode) | ~100ns (user-space) |
| Max count | ~thousands | ~millions |
| Managed by | OS kernel | Go runtime |
| Communication | shared memory + locks | channels (CSP model) |

```
OS Thread Model:               Go Goroutine Model:
┌──────────────┐               ┌──────────────────────────────┐
│  OS Kernel   │               │        Go Runtime            │
│  Scheduler   │               │  ┌────┐  ┌────┐  ┌────┐    │
│              │               │  │ G  │  │ G  │  │ G  │    │
│  T1  T2  T3  │               │  └────┘  └────┘  └────┘    │
│  │   │   │  │               │      Goroutine Queue         │
│  each = 1MB+ │               │         ↓                    │
└──────────────┘               │  ┌────┐  ┌────┐             │
                               │  │ P  │  │ P  │  Processors  │
                               │  └──┬─┘  └──┬─┘             │
                               │     │        │               │
                               │  ┌──▼─┐  ┌──▼─┐             │
                               │  │ M  │  │ M  │  OS Threads  │
                               │  └────┘  └────┘             │
                               └──────────────────────────────┘
```

---

## 2. GMP Model

The Go scheduler uses three entities: **G** (goroutine), **M** (machine/OS thread), **P** (processor/logical CPU).

### G — Goroutine
```go
// Internal runtime structure (simplified from src/runtime/runtime2.go)
type g struct {
    stack       stack       // [lo, hi] — current stack bounds
    stackguard0 uintptr     // stack overflow check
    m           *m          // current M executing this G (nil if not running)
    sched       gobuf       // saved registers (SP, PC) for context switch
    atomicstatus atomic.Uint32 // _Gidle | _Grunnable | _Grunning | _Gwaiting | _Gdead
    goid        int64       // unique goroutine ID
    waitsince   int64       // time G started waiting (for deadlock detection)
    preempt     bool        // preemption signal
}

// Goroutine states:
// _Gidle      = just allocated, not initialized
// _Grunnable  = in run queue, ready to run
// _Grunning   = executing on an M
// _Gsyscall   = in a system call
// _Gwaiting   = blocked (channel, mutex, timer, GC, etc.)
// _Gdead      = finished, can be reused
```

### M — Machine (OS Thread)
```go
type m struct {
    g0       *g     // goroutine with scheduling stack (special)
    curg     *g     // current running goroutine
    p        puintptr // attached P (nil if idle)
    nextp    puintptr // P to attach after syscall returns
    spinning bool   // m is looking for work, no G attached
    blocked  bool   // m is blocked on a note
    // ...
}
// g0 is special: it has a larger stack and runs scheduler code
// All scheduling decisions happen on g0's stack
```

### P — Processor (Logical CPU)
```go
type p struct {
    id          int32
    status      uint32   // _Pidle | _Prunning | _Psyscall | _Pgcstop
    m           muintptr // back-link to M running this P
    runq        [256]guintptr  // circular run queue (lock-free!)
    runqhead    atomic.Uint32
    runqtail    atomic.Uint32
    runnext     guintptr  // next G to run (bypass queue, highest priority)
    // mcache (memory allocator cache — per P, lock-free allocation)
    mcache      *mcache
}

// Number of Ps = GOMAXPROCS (default = runtime.NumCPU())
```

### Relationship Diagram
```
GOMAXPROCS = 4 → 4 Ps

P0 ──── M0 ──── G (running)
│               Local RunQ: [G1, G2, G3, ...]
│
P1 ──── M1 ──── G (running)
│               Local RunQ: [G4, G5]
│
P2 ──── M2 ──── (syscall) ─── G blocked in syscall
│               (P detached! P2 grabs idle M)
│
P3 ──── M3 ──── G (running)

         Global RunQ: [G10, G11, G12, ...]  ← overflow from local queues
         Network Poller: [G waiting on IO]
```

### How `go func()` Works Internally
```go
go someFunc()
// ↓ compiler rewrites to:
runtime.newproc(someFunc)

// runtime.newproc:
// 1. Allocates a G struct (from free list or new)
// 2. Sets up stack frame for someFunc
// 3. Sets G status to _Grunnable
// 4. Puts G in P's local runq (or runnext slot)
// 5. If there's an idle P, wake an M to run it
```

---

## 3. Goroutine Lifecycle

```
                    go func() {}
                         │
                         ▼
                    [_Grunnable] ──────────────────┐
                         │                         │
                    P picks it up                  │ preempted / yield
                         │                         │
                         ▼                         │
                    [_Grunning] ──────────────────►┘
                    │         │
               blocks         finishes
               (chan/mu)       │
                  │            ▼
                  ▼        [_Gdead] → recycled for next goroutine
             [_Gwaiting]
                  │
          event fires
          (chan send, mu unlock)
                  │
                  ▼
             [_Grunnable] → back to run queue
```

### Blocking Operations
```go
// These put G in _Gwaiting and release P to run other Gs:
ch <- val           // channel send (buffer full or no receiver)
<-ch                // channel receive (empty)
mu.Lock()           // mutex already locked
time.Sleep(d)       // timer
select {}           // forever block

// This does NOT block the M — Go wraps syscalls:
os.ReadFile(...)    // → netpoller or async IO on Darwin/Linux
// But raw syscalls DO block the M:
syscall.Read(fd, b, len(b))  // M is blocked, P detaches, finds new M
```

---

## 4. Work Stealing

When a P's local run queue is empty, it doesn't idle — it **steals** from other Ps.

```
Work Stealing Algorithm (src/runtime/proc.go: findRunnable):

1. Check P.runnext                    // highest priority
2. Check P.runq (local queue)         // fast, lock-free
3. Check global runq (1/61 of time)   // prevent starvation
4. Check network poller               // IO-ready goroutines
5. STEAL from another random P        // take half their runq
   └── steal(p2, &gp, max/2)
6. If still no work → M goes idle (spinning→parked)
```

```go
// Visualizing work stealing:
//
// P0.runq = [G1, G2, G3, G4]   P1.runq = []
//                │
//                └─── P1 steals half ───►  P1.runq = [G3, G4]
//
// P0.runq = [G1, G2]            P1.runq = [G3, G4]
// Both Ps now have work — maximum CPU utilization

// Key insight: local runq is a ring buffer of size 256
// Operations are lock-free using atomic CAS (compare-and-swap)
// Steal uses atomic operations — no global lock!
```

### Why Work Stealing Is Efficient
```
Traditional Thread Pool:         Go Work Stealing:
┌─────────────────────┐         ┌──────────────────────────────┐
│  Central Job Queue  │         │  P0:[G,G,G] P1:[G,G] P2:[]  │
│  [J,J,J,J,J,J,J]  │         │                    ↑steal     │
│         ↑           │         │            P2 steals from P0 │
│   global lock!      │         │  No global lock — CAS only   │
└─────────────────────┘         └──────────────────────────────┘
  Bottleneck under load            Scales to hundreds of Ps
```

---

## 5. Preemption

### Cooperative Preemption (Go < 1.14)
```go
// Goroutines only yielded at:
// - function calls (goroutine safe point check)
// - channel operations
// - time.Sleep
// - runtime.Gosched()

// Problem: tight loops could starve scheduler
func bad() {
    for { i++ }  // never yields! blocks entire P
}
```

### Asynchronous Preemption (Go 1.14+)
```go
// The runtime sends SIGURG signal to OS threads
// Signal handler: sets preempt flag on running G
// At next safe point: G is parked, scheduler runs

// sysmon goroutine checks every 10ms:
// - Any G running > 10ms? → send SIGURG to its M
// - This works even in tight loops

func fine() {
    for { i++ }  // Go 1.14+: preempted by SIGURG after 10ms
}

// Force a yield manually:
runtime.Gosched()  // yields processor, reschedules G
```

### Goroutine Parking & Unparking
```go
// park: G moves from _Grunning → _Gwaiting
// gopark(unlockf, lock, waitReason, ...)

// unpark: G moves from _Gwaiting → _Grunnable
// goready(g, traceskip)
// → puts G in run queue, may wake an idle M
```

---

## 6. Goroutine Stack Management

### Segmented → Contiguous Stacks
```
Go 1.2 and earlier: Segmented stacks (hot-split problem)
Go 1.3+:           Contiguous stacks (copy-on-grow)

Initial stack: 2KB (Go 1.4+, was 8KB before)
Max stack:     1GB (64-bit), 250MB (32-bit)
```

### Stack Growth
```go
// At each function call, compiler inserts a stack overflow check:
// if SP < stackguard0 { morestack() }

// morestack():
// 1. Allocate new stack 2x the current size
// 2. Copy all stack frames to new stack
// 3. Update all pointers (including heap pointers to stack)
// 4. Free old stack
// 5. Resume function

// This is why Go pointers to local vars work even after stack grows:
// All interior pointers are updated atomically
```

### Stack Shrinking
```go
// During GC: scan stacks for liveness
// If stack usage < 1/4 of stack size → shrink to 1/2
// Prevents memory waste from temporary deep call stacks
```

---

## 7. Garbage Collector

Go uses a **concurrent, tri-color, mark-and-sweep GC** with write barriers.

### Tri-Color Marking
```
Colors:
  White = not yet visited (candidates for collection)
  Grey  = visited, but children not yet scanned
  Black = visited, all children scanned (safe, keep)

Algorithm:
1. START: all objects are White
2. MARK roots as Grey (globals, stack vars, registers)
3. While Grey set non-empty:
   a. Pick a Grey object
   b. Mark all its White children as Grey
   c. Mark the object itself as Black
4. SWEEP: collect all remaining White objects

Invariant: A Black object never points to a White object
           (maintained by write barrier)
```

### Write Barrier
```go
// Without write barrier (not safe):
black.ptr = white  // black now points to white — violates invariant!

// Go's Hybrid Write Barrier (Go 1.17+ — Dijkstra+Yuasa):
// On pointer write: shade(old), shade(new)
// "shade" = mark grey if white

// This ensures: if black points to white, white gets greyed
// Implemented as a small snippet before every pointer store
```

### GC Phases
```
Phase 1: Mark Setup (STW — stop the world)
  - Enable write barriers
  - ~100µs pause

Phase 2: Concurrent Mark (runs alongside mutators)
  - Mark goroutine stacks
  - Scan heap — all objects reachable from roots
  - GC goroutines run at 25% CPU (GOGC target)

Phase 3: Mark Termination (STW)
  - Flush caches, disable write barriers
  - ~100µs pause

Phase 4: Concurrent Sweep
  - Reclaim white (dead) objects
  - Lazy: done during normal allocation
  - No STW needed
```

### GOGC and GC Tuning
```go
// GOGC=100 (default): trigger GC when heap doubles
// heap trigger = heap_live_after_last_GC * (1 + GOGC/100)

// Reduce GC pressure:
GOGC=200        // GC less often, more memory used
GOGC=off        // disable GC (dangerous — use for batch jobs)
GOMEMLIMIT=4GiB // Go 1.19+ — hard memory limit, triggers GC

// In code:
runtime.GC()               // force GC
runtime.ReadMemStats(&ms)  // inspect GC stats
debug.SetGCPercent(200)    // same as GOGC=200
debug.SetMemoryLimit(4 << 30) // 4GB hard limit

// Key stats:
ms.NumGC          // GC cycles count
ms.PauseTotalNs   // total STW time
ms.HeapAlloc      // live heap bytes
ms.GCSys          // GC metadata overhead
```

### Memory Allocator (TCMalloc-inspired)
```
Allocation tiers:

1. Tiny allocator (< 16 bytes, no pointers)
   → packed into 16-byte blocks, lock-free

2. Small objects (16B – 32KB)
   → mcache (per-P, lock-free) → mcentral (per class, locked) → mheap

3. Large objects (> 32KB)
   → directly from mheap (global, locked)

Size classes: 67 classes from 8B to 32KB
Each class has its own free list — avoids fragmentation

mcache: per-P cache → NO LOCK NEEDED for small allocs!
This is why Go allocation is so fast.
```

---

## 8. Memory Model & Happens-Before

### The Go Memory Model
```go
// Rule: if event A "happens before" event B,
//       then A's effects are visible to B.

// Channel send happens-before channel receive:
var x int
go func() {
    x = 42     // A
    ch <- 1    // send happens-before receive
}()
<-ch           // B
fmt.Println(x) // guaranteed to see 42

// Mutex unlock happens-before next lock:
mu.Unlock()  // A
// ...
mu.Lock()    // B — sees all writes before A's unlock

// sync.Once: first Do() happens-before all subsequent Do()s
```

### Data Race
```go
// DATA RACE: two goroutines access same memory,
// at least one write, no synchronization
var x int
go func() { x = 1 }()  // write
fmt.Println(x)          // read — RACE! undefined behavior

// Detect with: go run -race main.go
// go test -race ./...

// Fix options:
// 1. Channel: send/receive value
// 2. sync.Mutex / sync.RWMutex
// 3. sync/atomic for primitives
// 4. sync.Once for one-time init
```

---

## 9. Runtime Introspection

```go
import "runtime"
import "runtime/debug"

// Goroutine count (detect leaks)
fmt.Println(runtime.NumGoroutine())

// Force GC + print memory stats
runtime.GC()
var ms runtime.MemStats
runtime.ReadMemStats(&ms)
fmt.Printf("HeapAlloc=%dMB NumGC=%d PauseTotalNs=%dms\n",
    ms.HeapAlloc/1024/1024, ms.NumGC, ms.PauseTotalNs/1e6)

// Stack trace of all goroutines
buf := make([]byte, 1<<20)
n := runtime.Stack(buf, true) // true = all goroutines
fmt.Printf("%s\n", buf[:n])

// CPU profile
f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()

// Heap profile
f2, _ := os.Create("heap.prof")
pprof.WriteHeapProfile(f2)

// Live profiling endpoint
import _ "net/http/pprof"
go http.ListenAndServe(":6060", nil)
// Then: go tool pprof http://localhost:6060/debug/pprof/heap

// Goroutine dump
debug.PrintStack()           // current goroutine only
runtime.Stack(buf, true)     // all goroutines
```

### pprof Profiles
```bash
# CPU profile — where is time spent?
go tool pprof cpu.prof

# Heap profile — what's allocated?
go tool pprof -alloc_objects heap.prof

# Goroutine profile — what are goroutines doing?
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# Trace — scheduler events, GC, goroutine states
go test -trace trace.out ./...
go tool trace trace.out
```

---

## 10. Common Pitfalls & Tuning

### Pitfall 1: Goroutine Leak
```go
// Every blocked goroutine holds memory + stack
// runtime.NumGoroutine() should be stable over time

// Common leaks:
go func() { <-ch }()      // ch never receives
go func() { mu.Lock() }() // mu never unlocks
go func() {
    select {}              // forever
}()

// Fix: always use context cancellation
go func() {
    select {
    case <-ch:
    case <-ctx.Done(): return
    }
}()
```

### Pitfall 2: False Sharing
```go
// Two goroutines writing to different fields of same struct
// can still cause cache line contention (64 bytes = 1 cache line)

type Bad struct {
    a int64  // goroutine 1 writes
    b int64  // goroutine 2 writes — same cache line as a!
}

type Good struct {
    a   int64
    _   [56]byte // padding to separate cache lines
    b   int64
}
```

### Pitfall 3: Spinning vs Parking
```go
// BAD: busy-wait burns CPU
for !ready { runtime.Gosched() }

// GOOD: use sync primitives that park the goroutine
cond.Wait()        // parks G, wakes on Broadcast/Signal
<-readyCh          // parks G, wakes when value sent
```

### Pitfall 4: GOMAXPROCS in Containers
```go
// In containers, runtime.NumCPU() returns HOST CPU count!
// Container may only have 0.5 CPU quota → too many Ps → thrashing

import "go.uber.org/automaxprocs"
_ = automaxprocs.New() // auto-sets GOMAXPROCS to container quota

// Or manually:
runtime.GOMAXPROCS(2) // match container CPU limit
```

### Tuning Checklist
```
□ GOMAXPROCS matches container/host CPU quota
□ GOGC tuned for latency (200+) vs throughput (default 100)
□ GOMEMLIMIT set to prevent OOM (Go 1.19+)
□ No goroutine leaks (runtime.NumGoroutine() stable)
□ No data races (go test -race ./...)
□ Large allocations pooled with sync.Pool
□ Profiling done: pprof CPU + heap
□ False sharing eliminated in hot structs
□ Prealloc slices/maps where capacity known
□ Close channels to signal completion (not nil)
```

### sync.Pool — Reduce GC Pressure
```go
var pool = sync.Pool{
    New: func() any { return make([]byte, 4096) },
}

func handler(w http.ResponseWriter, r *http.Request) {
    buf := pool.Get().([]byte)    // reuse from pool
    defer pool.Put(buf)           // return to pool
    // use buf — avoids allocation per request
}
// Pool objects are cleared on each GC cycle
// Great for short-lived, frequently allocated objects (buffers, scratch space)
```

---

## Summary Diagram — Full Picture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Go Runtime                              │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    GMP Scheduler                         │  │
│  │                                                          │  │
│  │  Global RunQ        ┌──────┐  ┌──────┐  ┌──────┐       │  │
│  │  [G, G, G, ...]     │  P0  │  │  P1  │  │  P2  │       │  │
│  │                     │[G,G] │  │[G,G] │  │ [ ]  │       │  │
│  │  Network Poller     └──┬───┘  └──┬───┘  └──┬───┘       │  │
│  │  [G blocked on IO]     │         │    steal↗│           │  │
│  │                        M0        M1        M2           │  │
│  │                        │         │         │            │  │
│  └────────────────────────┼─────────┼─────────┼────────────┘  │
│                           │ syscall │         │               │
│  ┌────────────────────────▼─────────▼─────────▼────────────┐  │
│  │                    OS Kernel                             │  │
│  │             Threads: T0, T1, T2, ...                    │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌─────────────────────┐  ┌────────────────────────────────┐   │
│  │  Memory Allocator   │  │    Garbage Collector           │   │
│  │  mcache (per-P)     │  │    Tri-color mark-sweep        │   │
│  │  mcentral (shared)  │  │    Concurrent (25% CPU)        │   │
│  │  mheap (global)     │  │    STW pauses < 1ms            │   │
│  │  67 size classes    │  │    Write barrier                │   │
│  └─────────────────────┘  └────────────────────────────────┘   │
│                                                                 │
│  sysmon (background)                                           │
│  - Retract Ps from syscall-blocked Ms                         │
│  - Send SIGURG for preemption (every 10ms)                    │
│  - Force GC if idle too long                                   │
│  - Network poller check                                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## Key Takeaways

| Concept | Key Insight |
|---------|-------------|
| **G** | Goroutine = user-space thread, 2KB stack, millions possible |
| **M** | OS thread, 1:1 with kernel thread, blocked M loses its P |
| **P** | Logical processor, holds local run queue (256 slots), per-P mcache |
| **Work Stealing** | P steals half of another P's queue — lock-free CAS |
| **Preemption** | SIGURG signal forces yield at any point (Go 1.14+) |
| **Stack** | Grows 2x on overflow (copy), shrinks at GC |
| **GC** | Tri-color concurrent, write barrier, < 1ms STW pauses |
| **Allocator** | Per-P mcache = lock-free allocation for small objects |
| **GOGC** | Heap growth ratio before GC trigger; tune for latency vs memory |
| **Data Race** | Two concurrent accesses, one write, no sync = undefined behavior |

---

*Reference: [Go runtime source](https://github.com/golang/go/tree/master/src/runtime) — `proc.go`, `runtime2.go`, `malloc.go`, `mgc.go`*
