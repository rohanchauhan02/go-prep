# Q06: Mutex vs RWMutex Benchmark

## Problem Statement
Implement a thread-safe cache using sync.Mutex in one version and sync.RWMutex in another. Benchmark both under high-read/low-write scenarios and explain the performance difference.

## Key Concepts
```
sync.RWMutex, RLock/RUnlock, testing.B, go test -bench
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
