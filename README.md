# Go Interview Prep — 50 Concepts Cheatsheet

> Each concept has a minimal runnable example. Use this as a quick reference while solving the questions.

---

## 1. Goroutine Leak
A goroutine that blocks forever with no way to exit.
```go
// LEAK — blocks forever
go func() { val := <-ch }()

// FIX — use context
go func() {
    select {
    case val := <-ch: fmt.Println(val)
    case <-ctx.Done(): return  // exits cleanly
    }
}()
```

---

## 2. Deadlock
Two goroutines each waiting for the other's lock.
```go
// FIX: always acquire locks in the same order
var muA, muB sync.Mutex
muA.Lock(); muB.Lock()   // goroutine 1
muA.Lock(); muB.Lock()   // goroutine 2  ✅ consistent = no deadlock
```

---

## 3. Directional Channels
```go
func produce(out chan<- int) { out <- 42 }   // send-only
func consume(in <-chan int)  { fmt.Println(<-in) }  // receive-only
```

---

## 4. Select + Timeout
```go
select {
case res := <-svcA: fmt.Println(res)
case res := <-svcB: fmt.Println(res)
case <-time.After(2 * time.Second): fmt.Println("timeout")
}
```

---

## 5. WaitGroup vs ErrGroup
```go
// WaitGroup — no error handling
var wg sync.WaitGroup
wg.Add(1); go func() { defer wg.Done(); work() }()
wg.Wait()

// ErrGroup — stops all on first error
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return doWork(ctx) })
err := g.Wait()
```

---

## 6. RWMutex
```go
var mu sync.RWMutex
// Multiple readers at once:
mu.RLock(); val := cache[key]; mu.RUnlock()
// Only one writer at a time:
mu.Lock(); cache[key] = val; mu.Unlock()
```

---

## 7. sync.Once (Singleton)
```go
var instance *DB
var once sync.Once

func GetDB() *DB {
    once.Do(func() { instance = &DB{} }) // runs exactly once
    return instance
}
```

---

## 8. Atomic Operations
```go
var count atomic.Int64
count.Add(1)          // thread-safe, no mutex
fmt.Println(count.Load())
```

---

## 9. Nil Interface Trap ⚠️
```go
// BUG: typed nil — interface is NOT nil
func getErr() error {
    var e *MyError = nil
    return e            // interface{type=*MyError, value=nil} ≠ nil!
}

// FIX: return untyped nil
func getErr() error {
    return nil          // interface{type=nil, value=nil} == nil ✅
}
```

---

## 10. Struct Embedding
```go
type Animal struct{ Name string }
func (a Animal) Speak() string { return a.Name + " speaks" }

type Dog struct{ Animal }  // promotes Speak()

d := Dog{Animal{"Rex"}}
d.Speak()   // "Rex speaks" — promoted from Animal
d.Name      // "Rex"        — promoted field
```

---

## 11. Closure Variable Capture
```go
// BUG: all closures capture same i
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()  // prints 3,3,3
}

// FIX 1: shadow
for i := 0; i < 3; i++ {
    i := i
    go func() { fmt.Println(i) }()  // prints 0,1,2
}

// FIX 2: parameter
for i := 0; i < 3; i++ {
    go func(n int) { fmt.Println(n) }(i)
}
```

---

## 12. Defer + Panic + Recover
```go
// Defer runs LIFO (last in, first out)
defer fmt.Println("1st defer — runs last")
defer fmt.Println("2nd defer — runs first")

// Recover panic inside defer
defer func() {
    if r := recover(); r != nil {
        fmt.Println("caught:", r)
    }
}()
panic("boom")

// Named return modified by defer
func double() (result int) {
    defer func() { result *= 2 }()
    result = 5
    return  // returns 10, not 5
}
```

---

## 13. Slice Internals
```go
// Slice header: [pointer | len | cap]
orig := []int{1, 2, 3, 4, 5}
sub  := orig[1:3]      // shares same backing array!
sub[0] = 99            // orig[1] is now 99 — BUG

// Fix with copy
sub2 := make([]int, 2)
copy(sub2, orig[1:3])  // independent copy
```

---

## 14. Concurrent Map
```go
// PANIC: concurrent map write without lock
go func() { m["k"] = 1 }()  // ← data race!

// FIX 1: sync.Mutex
mu.Lock(); m["k"] = 1; mu.Unlock()

// FIX 2: sync.Map
var sm sync.Map
sm.Store("k", 1)
v, ok := sm.Load("k")
sm.Range(func(k, v any) bool { fmt.Println(k, v); return true })
```

---

## 15. String vs Rune
```go
s := "Hello, 世界"
len(s)                          // 13 (bytes)
utf8.RuneCountInString(s)       // 9  (runes/characters)

// Range gives runes:
for idx, ch := range s { fmt.Println(idx, ch) }
// Indexing gives bytes:
fmt.Println(s[7])  // 228 (first byte of 世, not the character)

// Reverse correctly:
r := []rune(s)
// reverse r, then string(r)
```

---

## 16. GOMAXPROCS
```go
runtime.GOMAXPROCS(1)              // single-threaded (concurrency, no parallelism)
runtime.GOMAXPROCS(runtime.NumCPU()) // full parallelism
```

---

## 17. Context
```go
ctx := context.WithValue(ctx, "userID", "42")
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// Read value
uid := ctx.Value("userID").(string)

// Check cancellation
select {
case <-ctx.Done(): return ctx.Err()
default: doWork()
}
```

---

## 18. Custom Errors + Wrapping
```go
type ValidationError struct{ Code int; Msg string }
func (e *ValidationError) Error() string { return e.Msg }

// Wrap
err := fmt.Errorf("process failed: %w", &ValidationError{400, "bad input"})

// Unwrap
var ve *ValidationError
errors.As(err, &ve)   // true, ve.Code == 400
errors.Is(err, ErrNotFound)  // checks chain
```

---

## 19. Generics
```go
// Type parameter with constraint
func Map[T, U any](s []T, fn func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s { result[i] = fn(v) }
    return result
}

Map([]int{1,2,3}, func(n int) string { return fmt.Sprint(n) })
// → ["1","2","3"]

// Number constraint
type Number interface{ ~int | ~float64 }
func Sum[T Number](nums []T) T { ... }
```

---

## 20. Reflection
```go
t := reflect.TypeOf(myStruct)
v := reflect.ValueOf(myStruct)

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    tag   := field.Tag.Get("json")
    value := v.Field(i).Interface()
}
```

---

## 21. Escape Analysis
```go
func stack() int  { x := 42; return x  }   // x stays on stack
func heap()  *int { x := 42; return &x }   // x escapes to heap

// Diagnose:
// go build -gcflags="-m" main.go
```

---

## 22. Rate Limiter (Token Bucket)
```go
// tokens refill at `rate` per second, max `capacity`
func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock(); defer tb.mu.Unlock()
    tb.refill()
    if tb.tokens >= 1 { tb.tokens--; return true }
    return false
}
```

---

## 23. Worker Pool
```go
jobs := make(chan Job, 100)
for i := 0; i < numWorkers; i++ {
    go func() {
        for job := range jobs {
            process(job)
        }
    }()
}
// Send jobs; close(jobs) to stop workers
```

---

## 24. Fan-Out / Fan-In
```go
// Fan-out: one source → multiple workers
for i := 0; i < 3; i++ {
    go worker(ctx, i, src)  // all read from same src
}

// Fan-in: merge multiple channels into one
func merge(channels ...<-chan string) <-chan string {
    out := make(chan string)
    for _, ch := range channels {
        go func(c <-chan string) { for v := range c { out <- v } }(ch)
    }
    return out
}
```

---

## 25. Pub/Sub Event Bus
```go
type EventBus struct {
    mu   sync.RWMutex
    subs map[string][]chan string
}

func (eb *EventBus) Publish(topic, msg string) {
    for _, ch := range eb.subs[topic] {
        select {
        case ch <- msg:   // deliver
        default:          // skip slow subscriber (non-blocking)
        }
    }
}
```

---

## 26. LRU Cache
```go
// O(1) Get and Put using doubly-linked list + hashmap
type LRU struct {
    cap   int
    list  *list.List
    items map[string]*list.Element
}
// Get → MoveToFront
// Put → PushFront; if full → Remove(list.Back())
```

---

## 27. Trie
```go
type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool
}
// Insert: walk char by char, create nodes, set isEnd=true
// Search: walk and check isEnd
// Autocomplete: DFS from prefix node collecting isEnd words
```

---

## 28. Binary Search Variants
```go
// First occurrence: when found, save result, go left (hi = mid-1)
// Last occurrence:  when found, save result, go right (lo = mid+1)
// Rotated array:    one half is always sorted — check which, then narrow
mid := lo + (hi-lo)/2  // avoids int overflow vs (lo+hi)/2
```

---

## 29. Graph Traversal
```go
// BFS — queue
queue := []int{start}
for len(queue) > 0 {
    node := queue[0]; queue = queue[1:]
    for _, nb := range adj[node] { queue = append(queue, nb) }
}

// Cycle detection — visited + recursion stack
// Topo sort — Kahn's: process nodes with in-degree 0
```

---

## 30. Dynamic Programming
```go
// Coin Change (bottom-up)
dp := make([]int, amount+1)
for i := 1; i <= amount; i++ {
    for _, c := range coins {
        if c <= i { dp[i] = min(dp[i], dp[i-c]+1) }
    }
}

// LCS (2D table)
if a[i-1] == b[j-1] { dp[i][j] = dp[i-1][j-1] + 1 }
else { dp[i][j] = max(dp[i-1][j], dp[i][j-1]) }
```

---

## 31. Stringer Interface
```go
type Color struct{ R, G, B uint8 }

func (c Color) String() string {
    return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
}
// Now fmt.Println(c) → "rgb(255,0,0)"
// %#v calls GoString() instead
```

---

## 32. Custom JSON Marshal
```go
type UnixTime struct{ time.Time }

func (u UnixTime) MarshalJSON() ([]byte, error) {
    return json.Marshal(u.Unix())     // → int64
}
func (u *UnixTime) UnmarshalJSON(b []byte) error {
    var ts int64
    json.Unmarshal(b, &ts)
    u.Time = time.Unix(ts, 0)
    return nil
}
```

---

## 33. Graceful Shutdown
```go
srv := &http.Server{Addr: ":8080"}
go srv.ListenAndServe()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit  // block until signal

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)  // stops accepting, waits for in-flight
```

---

## 34. Middleware Chain
```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
    for i := len(mws)-1; i >= 0; i-- { h = mws[i](h) }
    return h
}
// Usage: Chain(handler, Logger, Auth, RateLimit)
// Executes: Logger → Auth → RateLimit → handler
```

---

## 35. Interceptor Pattern
```go
type Interceptor func(ctx context.Context, req Request, next Handler) (Response, error)

// Chain interceptors like middleware
func ChainInterceptors(h Handler, ics ...Interceptor) Handler {
    for i := len(ics)-1; i >= 0; i-- {
        ic, next := ics[i], h
        h = func(ctx context.Context, req Request) (Response, error) {
            return ic(ctx, req, next)
        }
    }
    return h
}
```

---

## 36. Heap (Priority Queue)
```go
type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }  // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any          { old:=*h; x:=old[len(old)-1]; *h=old[:len(old)-1]; return x }

heap.Init(&h); heap.Push(&h, 5); heap.Pop(&h)
```

---

## 37. Sliding Window
```go
// Two-pointer pattern
left := 0
for right := 0; right < len(s); right++ {
    window[s[right]]++
    for window[s[right]] > 1 {  // shrink
        window[s[left]]--; left++
    }
    maxLen = max(maxLen, right-left+1)
}
```

---

## 38. Singleton Logger
```go
var once sync.Once
var logger *Logger

func GetLogger() *Logger {
    once.Do(func() {
        f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        out := io.MultiWriter(os.Stdout, f)
        logger = &Logger{logger: log.New(out, "", log.LstdFlags)}
    })
    return logger
}
```

---

## 39. Interface Composition + Adapter
```go
// Compose small interfaces
type ReadWriteCloser interface{ io.Reader; io.Writer; io.Closer }

// Adapter: make bytes.Buffer satisfy ReadWriteCloser
type BufferRWC struct{ bytes.Buffer }
func (b *BufferRWC) Close() error { b.Reset(); return nil }
```

---

## 40. Type Switch
```go
switch v := i.(type) {
case int:    fmt.Println("int:", v)
case string: fmt.Println("str:", v)
case nil:    fmt.Println("nil")
default:     fmt.Printf("unknown: %T\n", v)
}

// Safe type assertion (comma-ok)
s, ok := i.(string)  // ok=false if wrong type, no panic
```

---

## 41. Recursive Descent Parser
```go
// Grammar: expr = term ((+|-) term)*
//          term = factor ((*|/) factor)*
//          factor = NUMBER | '(' expr ')'
//
// Operator precedence enforced by grammar levels:
// lower in grammar = lower precedence = evaluated last
```

---

## 42. Semaphore (Buffered Channel)
```go
sem := make(chan struct{}, 3)  // max 3 concurrent

acquire := func() { sem <- struct{}{} }
release := func() { <-sem }

go func() {
    acquire(); defer release()
    doWork()
}()
```

---

## 43. Functional Options
```go
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(opts ...Option) *Server {
    s := &Server{timeout: 30*time.Second}  // defaults
    for _, o := range opts { o(s) }
    return s
}

NewServer(WithTimeout(10*time.Second), WithMaxConns(100))
```

---

## 44. Compile-Time Interface Check
```go
// If MyWriter doesn't implement io.Writer, compile fails here:
var _ io.Writer = (*MyWriter)(nil)
var _ io.ReadWriteCloser = (*MyReadWriter)(nil)
```

---

## 45. Barrier (sync.Cond)
```go
cond := sync.NewCond(&mu)

// Goroutine waits at barrier
mu.Lock()
count++
if count == total { cond.Broadcast() }  // release all
else { cond.Wait() }                    // block until Broadcast
mu.Unlock()
```

---

## 46. String Builder Performance
```go
// Slowest — O(n²) allocations
s := ""
for i := 0; i < n; i++ { s += "a" }

// Fast — single allocation
var sb strings.Builder
sb.Grow(n)  // pre-allocate
for i := 0; i < n; i++ { sb.WriteByte('a') }
s := sb.String()
```

---

## 47. Struct Memory Alignment
```go
// BAD — 32 bytes (padding waste)
type Bad struct {
    a bool    // 1 + 7 pad
    b float64 // 8
    c bool    // 1 + 7 pad
    d int32   // 4 + 4 pad
}

// GOOD — 16 bytes (largest fields first)
type Good struct {
    b float64 // 8
    d int32   // 4
    a bool    // 1
    c bool    // 1 + 2 pad
}
// Rule: sort fields largest alignment → smallest
```

---

## 48. Tree Serialization
```go
// Preorder: root → left → right, "#" for nil
// "1,2,4,#,#,5,#,#,3,#,6,#,#"
func serialize(node *TreeNode) string {
    if node == nil { return "#" }
    return fmt.Sprintf("%d,%s,%s", node.Val,
        serialize(node.Left), serialize(node.Right))
}
```

---

## 49. Two-Phase Commit
```go
// Phase 1: Coordinator asks all participants to vote
// Phase 2: if ALL vote COMMIT → broadcast COMMIT
//          if ANY votes ABORT → broadcast ABORT
//
// Use channels as message-passing between goroutines
voteCh := make(chan Vote, 1)     // participant → coordinator
decisionCh := make(chan Vote, 1) // coordinator → participant
```

---

## 50. Mini ORM with Reflection
```go
// Generate SQL from struct tags
t := reflect.TypeOf(v)
for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    tag   := field.Tag.Get("db")  // "id", "name", etc.
    kind  := field.Type.Kind()    // reflect.Int, reflect.String...
}

// Scan rows into structs
rv := reflect.ValueOf(target).Elem()          // *[]User → []User
elem := reflect.New(elemType).Elem()          // new User{}
elem.Field(i).Set(reflect.ValueOf(rowVal))    // set field value
```

---

## Quick Reference

| # | Concept | Key API |
|---|---------|---------|
| 1 | Goroutine Leak | `ctx.Done()`, `select` |
| 2 | Deadlock | consistent lock order |
| 3 | Channel Direction | `chan<-`, `<-chan` |
| 4 | Select+Timeout | `time.After()` |
| 5 | ErrGroup | `errgroup.WithContext` |
| 6 | RWMutex | `RLock/RUnlock` |
| 7 | Once | `sync.Once.Do()` |
| 8 | Atomic | `atomic.Int64.Add()` |
| 9 | Nil Interface | typed nil trap |
| 10 | Embedding | method promotion |
| 11 | Closures | `i := i` shadow fix |
| 12 | Defer/Panic | `recover()`, LIFO |
| 13 | Slice | backing array, `copy()` |
| 14 | Sync Map | `sync.Map` |
| 15 | Runes | `[]rune`, `range` |
| 16 | GOMAXPROCS | `runtime.GOMAXPROCS` |
| 17 | Context | `WithValue/Timeout/Cancel` |
| 18 | Errors | `errors.As/Is`, `%w` |
| 19 | Generics | `[T any]`, constraints |
| 20 | Reflection | `reflect.TypeOf/ValueOf` |
| 21 | Escape | `-gcflags="-m"` |
| 22 | Rate Limit | token bucket |
| 23 | Worker Pool | buffered chan |
| 24 | Fan-Out/In | merge channels |
| 25 | Pub/Sub | non-blocking send |
| 26 | LRU | `container/list` |
| 27 | Trie | prefix tree, DFS |
| 28 | Binary Search | lo+hi/2 overflow |
| 29 | Graph | BFS queue, DFS stack |
| 30 | DP | tabulation, memoize |
| 31 | Stringer | `String() string` |
| 32 | JSON Marshal | `MarshalJSON` |
| 33 | Graceful Shutdown | `srv.Shutdown(ctx)` |
| 34 | Middleware | `func(Handler) Handler` |
| 35 | Interceptor | closure chain |
| 36 | Heap | `container/heap` |
| 37 | Sliding Window | two pointers |
| 38 | Logger | `sync.Once + MultiWriter` |
| 39 | Adapter | interface composition |
| 40 | Type Switch | `v := i.(type)` |
| 41 | Parser | recursive descent |
| 42 | Semaphore | buffered channel |
| 43 | Options | `func(*Config)` |
| 44 | Interface Check | `var _ I = (*T)(nil)` |
| 45 | Barrier | `sync.Cond.Broadcast` |
| 46 | Builder | `strings.Builder.Grow` |
| 47 | Alignment | `unsafe.Sizeof` |
| 48 | Serialize | preorder + `#` |
| 49 | 2PC | channel messaging |
| 50 | ORM | `reflect` + struct tags |

