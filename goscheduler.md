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

---

---

# Plain English Explanations

> Everything above in human language. Read this if the code blocks feel dense. Then go back and the code will click immediately.

---

## 1. Thread vs Goroutine — Explained

An **OS thread** is created by the operating system kernel. Every thread gets a fixed block of memory (usually 1–8 MB) reserved for its stack, even if it only uses a fraction of it. Creating a thread involves a kernel syscall which takes ~1 millisecond. Switching between threads (context switching) forces the CPU into kernel mode and takes 1–10 microseconds. Because of this cost, you can only realistically run thousands of threads before the machine runs out of memory and CPU time.

A **goroutine** is not a thread. It is a lightweight unit of work managed entirely by the Go runtime in user space — the kernel knows nothing about it. A goroutine starts with just 2 KB of stack and grows on demand. Creating a goroutine takes ~0.3 microseconds (3000x faster than a thread) because no kernel call is needed. Switching between goroutines takes ~100 nanoseconds because the runtime does it entirely in user space without involving the OS.

This means a Go program can comfortably run **millions of goroutines** on just a handful of OS threads — the Go runtime multiplexes them efficiently.

The analogy: OS threads are like trucks — powerful but expensive and slow to deploy. Goroutines are like motorcycles — cheap, fast, and you can have thousands of them weaving through the same lanes (OS threads).

---

## 2. GMP Model — Explained

The GMP model is the heart of Go's scheduler. It has three building blocks:

**G — Goroutine.** This represents a unit of concurrent work. Every time you write `go func()`, the runtime creates a G. A G holds the function to run, its current stack, saved CPU registers (so it can be paused and resumed), and a status flag (runnable, running, waiting, dead). A G is just a struct in memory — very cheap.

**M — Machine (OS Thread).** An M is an actual OS thread. It is the thing that physically executes code on a CPU core. An M can only run one G at a time. The number of Ms in a program can grow beyond GOMAXPROCS — extra Ms are created when goroutines block on system calls (so other Ps can keep working).

**P — Processor (Logical CPU).** A P is a logical processor — think of it as a "slot" that connects an M to a run queue of Gs. The number of Ps is controlled by `GOMAXPROCS` (default = number of CPU cores). Each P has its own private **run queue** of up to 256 goroutines. Crucially, each P also has its own **memory cache (mcache)** so small allocations don't need any locks.

**How they work together:** An M must hold a P to run goroutines. When M picks up a P, it looks in P's local run queue, picks the next G, and starts executing it. If P's queue is empty, the M steals work from other Ps. If M blocks on a system call (like reading a file), the P immediately detaches from M and attaches to a different idle M — so the P keeps running other goroutines uninterrupted.

The ratio is roughly: **N Ps** (one per CPU core), **N or more Ms** (can temporarily exceed P count during syscalls), **millions of Gs** distributed across Ps.

---

## 3. Goroutine Lifecycle — Explained

A goroutine moves through five states during its lifetime:

**Runnable:** The goroutine is ready to run but is waiting for a P to pick it up. It sits in a run queue (either a P's local queue or the global queue).

**Running:** The goroutine is actively executing on an M. Only one goroutine can be in the Running state per P at any time.

**Waiting:** The goroutine is blocked — it's waiting for something. This might be a channel receive/send, a mutex lock, a timer, a network read, or a GC pause. While waiting, it does NOT occupy an M or P. The M and P are free to run other goroutines. When the awaited event happens (e.g. someone sends to the channel), the goroutine moves back to Runnable.

**Syscall:** The goroutine is inside an OS system call. The M is blocked in kernel mode. When this happens, the P immediately detaches and either finds an idle M or creates a new one. This is how Go avoids stalling other goroutines when one does file IO.

**Dead:** The goroutine has finished. The runtime doesn't immediately discard the G struct — it reuses it for the next `go` call to avoid allocation overhead.

The key insight: **blocked goroutines don't waste CPU**. They park themselves and give the M+P back to productive work.

---

## 4. Work Stealing — Explained

Work stealing solves the load balancing problem: what should a P do when its run queue is empty while other Ps are overloaded?

The naive answer — a global job queue — creates a bottleneck because every P would need to lock it to read or write. Go avoids this with per-P local queues and work stealing.

**Here's the stealing algorithm in order:**
1. Check if there's a goroutine in `P.runnext` (the highest-priority slot, set when a goroutine is just unblocked).
2. Check P's own local run queue (circular buffer of 256 slots, lock-free using atomic operations).
3. Every 61 scheduling ticks, check the global run queue. This prevents goroutines in the global queue from starving if local queues stay full.
4. Check the network poller for goroutines that are now ready (IO completed).
5. **Steal** from another randomly chosen P — take exactly half their local queue.

The steal operation uses **atomic compare-and-swap (CAS)** instructions — no locks, no kernel involvement. This is why Go's scheduler scales well even with hundreds of Ps.

If after all this there's still no work, the M marks itself as "spinning" (actively looking for work) and eventually parks itself (goes to sleep until signalled by a `go` call).

**Why stealing half?** Taking half distributes work evenly. Taking one at a time would cause too many steal attempts; taking all would create the same imbalance on the other side.

---

## 5. Preemption — Explained

**The problem:** What if a goroutine runs a tight CPU loop and never calls any Go functions or does any IO? Early versions of Go's scheduler relied on *cooperative preemption* — goroutines had to voluntarily yield by calling any function (which includes a stack overflow check). A tight loop with no function calls could monopolize a P forever, starving other goroutines.

**Old solution (Go < 1.14) — Cooperative:** The Go compiler inserts a preemption check at the beginning of every function call. When `sysmon` (a background goroutine) detects a goroutine has been running for more than 10ms, it sets a `preempt` flag on that goroutine's G struct. The next time the goroutine calls any function, the check fires, and the goroutine is paused.

**New solution (Go 1.14+) — Asynchronous:** The runtime sends a `SIGURG` Unix signal to the OS thread running the goroutine. The signal handler interrupts execution at the next "safe point" (a point where the GC can safely inspect the stack). This works even in tight loops with no function calls. The goroutine is paused and another one runs.

`runtime.Gosched()` is a manual yield — you call it explicitly to tell the scheduler "I'm done with my time slice, run someone else." Useful in CPU-bound goroutines that want to be cooperative.

The `sysmon` goroutine (system monitor) runs in a dedicated OS thread and wakes up every 20µs to 10ms, handling preemption, retaking Ps from goroutines stuck in long system calls, and triggering GC.

---

## 6. Goroutine Stack Management — Explained

Every goroutine starts with a tiny 2 KB stack (down from 8 KB in early Go versions). This is so cheap that you can create millions of goroutines without running out of memory. But what if a goroutine needs more stack — deep recursive calls, large local arrays?

**Stack growth:** The Go compiler inserts a stack overflow check at the start of every function. If the current stack pointer (SP) is close to the bottom of the allocated stack, Go calls `morestack()`. This function allocates a **new, larger stack (2x the current size)**, copies all existing stack frames to the new stack, updates every pointer on the stack that points to another stack location (these would be dangling after the copy), and then resumes the function. The old stack is freed.

This is called a **contiguous stack** (vs Go's old "segmented stack" model which had the "hot-split" performance problem).

**Stack shrinking:** During GC, the runtime scans all goroutine stacks. If a stack's actual usage has fallen below 25% of its allocated size (e.g. a recursive function returned and unwound deeply), Go shrinks the stack to half its current size. This reclaims memory from goroutines that had a burst of deep calls but are now mostly idle.

**Why this matters for you:** Stack-allocated variables in Go are safe to point to even after the function that created them returns — because if the stack grows, all pointers are updated. Go's escape analysis decides at compile time whether a variable should live on the stack or the heap, and it makes the right choice automatically.

---

## 7. Garbage Collector — Explained

Go uses a **concurrent, tri-color, mark-and-sweep garbage collector**. Let's break each word down:

**Mark-and-sweep:** GC happens in two phases. First, *mark* — traverse the object graph from all roots (global variables, goroutine stacks, CPU registers) and mark every reachable object as "alive". Second, *sweep* — reclaim all memory occupied by objects that were *not* marked (they're unreachable and therefore garbage).

**Tri-color:** Objects are assigned one of three colors at any moment. **White** means "not yet visited" (potentially garbage). **Grey** means "visited but children not yet scanned". **Black** means "fully scanned, keep this". The GC starts by making everything white, marks roots grey, and then processes grey objects until none remain — everything still white is garbage.

**Concurrent:** The mark phase runs *alongside* your application code (the "mutator"). Your goroutines keep running while GC marks objects. This is why Go's GC pauses are sub-millisecond — you only stop the world briefly at the start and end of the mark phase (~100µs each), not during the entire collection.

**Write barrier:** The challenge of concurrent marking is that your program might modify pointers while GC is scanning. For example, a goroutine might take a white (unmarked) object and store it inside a black (fully scanned) object — now GC would never find it and would incorrectly free it. The *write barrier* prevents this: every pointer store in your program includes a tiny runtime snippet that "shades" the old and new pointer values grey if they're white. This keeps the invariant "no black object points directly to a white object" intact throughout concurrent marking.

**GOGC:** This is the knob that controls how often GC runs. `GOGC=100` (default) means GC runs when the heap has grown 100% from the size after the last collection (i.e. doubles). Set `GOGC=200` to reduce GC frequency at the cost of using more memory. Set `GOMEMLIMIT` (Go 1.19+) as an absolute memory ceiling.

**Memory Allocator:** Go's allocator is inspired by TCMalloc. There are 67 size classes from 8 bytes to 32 KB. Small allocations (< 32 KB) go through `mcache` — a per-P cache with no locking required (since only one goroutine runs per P at a time). When the mcache is empty, it refills from `mcentral` (shared, locked). Large allocations (> 32 KB) go directly to `mheap`. This tiered system is why Go allocation is extremely fast for small objects.

---

## 8. Memory Model & Happens-Before — Explained

The **Go Memory Model** is a specification that defines when one goroutine is *guaranteed* to see the memory writes of another goroutine.

**The problem:** Modern CPUs reorder instructions and cache writes in registers. Without explicit synchronization, there is no guarantee that a write in goroutine A is visible to goroutine B, even if A wrote first in wall-clock time.

**Happens-before** is the formal way to express visibility guarantees. If operation A *happens-before* operation B, then B is guaranteed to see all memory effects of A.

Go's rules for happens-before relationships:
- A **channel send** happens-before the **corresponding receive** completes. So if you write a value then send on a channel, the receiver is guaranteed to see that write.
- A **mutex Unlock** happens-before the next **Lock** on the same mutex. So all writes done before unlocking are visible to the next locker.
- `sync.Once.Do(f)` — the completion of `f` happens-before the return of any other `Do` call. Singleton initialization is safe.
- A **goroutine's start** happens-before any code in that goroutine.

**Data race:** A data race occurs when two goroutines access the same variable concurrently, at least one access is a write, and there is no synchronization (no mutex, channel, atomic). The Go Memory Model says the behavior is *undefined* — you might read stale data, corrupt data, or trigger a crash. Always use `-race` flag in tests: `go test -race ./...`. The race detector adds ~5-10x overhead but catches races reliably.

---

## 9. Runtime Introspection — Explained

Go gives you first-class tools to observe the runtime from inside your code and from external tools.

**`runtime.NumGoroutine()`** returns the number of live goroutines right now. Call it before and after suspected leak areas. A continuously growing number means goroutines are leaking.

**`runtime.ReadMemStats(&ms)`** fills a `MemStats` struct with everything about heap usage, GC cycles, pause times, and allocations. Key fields: `HeapAlloc` (live bytes), `NumGC` (total GC cycles), `PauseTotalNs` (total STW time), `Mallocs` and `Frees` (allocation rate). Force a GC first with `runtime.GC()` to get up-to-date heap stats.

**`runtime.Stack(buf, true)`** dumps the stack trace of all goroutines into a buffer. The `true` argument means "include all goroutines, not just the current one". This is invaluable for diagnosing what every goroutine is waiting on. It's the same output you see on a panic.

**`pprof`** is Go's profiling framework. It captures:
- **CPU profile** — where is your program spending CPU time? Samples the call stack every 10ms.
- **Heap profile** — what is allocated on the heap and where was it allocated?
- **Goroutine profile** — stack traces of all current goroutines (great for leak detection).
- **Block profile** — where are goroutines blocking on channels and mutexes?
- **Mutex profile** — which mutexes are contended?

Importing `_ "net/http/pprof"` automatically registers HTTP endpoints under `/debug/pprof/` so you can profile a live production service without restarting.

**`go tool trace`** gives you a visual timeline of scheduler events — which goroutine ran on which P, when GC paused, when goroutines blocked and woke up. It's the most detailed view of Go runtime behavior.

---

## 10. Common Pitfalls & Tuning — Explained

**Goroutine leaks** are the most common production issue in Go services. Every blocked goroutine consumes memory (its stack) and may hold references that prevent GC from collecting objects. The fix is always the same: use context cancellation so goroutines know when to exit. Monitor `runtime.NumGoroutine()` in your health endpoint.

**False sharing** is a subtle CPU cache performance issue. CPU caches work in 64-byte "cache lines". If two goroutines on different CPU cores frequently write to different fields of the same struct, they may land on the same cache line. Every write forces the other core to invalidate its cache copy — causing dramatic slowdowns despite having no logical conflict. Fix: add padding between hot fields to put them on separate cache lines.

**Spinning vs Parking:** If you implement a wait loop with `runtime.Gosched()`, you're wasting CPU — the goroutine loops constantly doing nothing useful. Prefer `sync.Cond.Wait()`, channel receives, or `time.Sleep()` which properly *park* the goroutine (remove it from the scheduler entirely until woken up).

**GOMAXPROCS in containers** is a classic trap. Go sets GOMAXPROCS to `runtime.NumCPU()` at startup. But in a Docker/Kubernetes container with CPU limits, `NumCPU()` returns the *host* machine's core count, not your container's quota (e.g. 0.5 CPU). Having 32 Ps competing for 0.5 CPU causes excessive context switching and high latency. Use the `automaxprocs` library or set GOMAXPROCS manually to match your CPU limit.

**`sync.Pool`** is a cache for reusable objects. Instead of allocating a new 4 KB buffer for every HTTP request, you `Get()` one from the pool, use it, and `Put()` it back. This dramatically reduces GC pressure in high-throughput services. Important: Pool objects are cleared at every GC cycle, so don't use Pool to store persistent state — only scratch buffers and temporary objects.

**GOGC and GOMEMLIMIT tuning:**
- Latency-sensitive services (APIs, game servers): raise GOGC to 200–400 so GC runs less frequently. Accept higher memory usage in exchange for fewer GC pauses.
- Memory-constrained environments (containers): set GOMEMLIMIT to your container's memory limit minus headroom. This prevents OOM kills by triggering GC earlier when memory is tight.
- Batch jobs: set GOGC=off and trigger `runtime.GC()` manually at checkpoints for maximum throughput.

---

## The Big Picture — Everything Connected

Here is how all the pieces fit together when you write `go doWork()`:

1. The compiler generates a call to `runtime.newproc`, which allocates a **G** struct (2 KB stack), sets it to `_Grunnable`, and places it in the current **P**'s `runnext` slot.

2. The P's **M** finishes the current G (or the G blocks), picks up the new G, and starts executing it. This is purely user-space — no kernel involved.

3. If `doWork` blocks on a channel: the G transitions to `_Gwaiting`. The M immediately picks the next G from the local run queue. The blocked G is stored inside the channel's wait queue. Zero CPU wasted.

4. If `doWork` calls a blocking system call (like `read`): the M enters kernel mode. Go's runtime detects this (via `entersyscall`) and **detaches the P** from the M. The P finds or creates another M and keeps running other goroutines. When the syscall returns, the original M tries to reclaim a P — if none available, it parks itself and the G goes to the global run queue.

5. While all this happens, the **GC** runs concurrently on a separate set of goroutines, marking live objects. It uses write barriers to stay consistent with your running program. Every 100µs or so you get a sub-millisecond STW pause for coordination.

6. **sysmon** runs in the background watching for goroutines that have been running too long (> 10ms) and sends `SIGURG` to preempt them. It also rebalances Ps stuck in system calls and polls the network for IO-ready goroutines.

7. **Work stealing** ensures that if one P finishes all its work while another P has 200 goroutines queued, the idle P takes 100 of them — keeping all CPU cores productive.

The result: **millions of goroutines, microsecond-level scheduling, sub-millisecond GC pauses, near-linear CPU scaling** — all built into a language with no manual memory management and simple concurrency primitives.

