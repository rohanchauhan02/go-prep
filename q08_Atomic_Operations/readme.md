# Q08: Atomic Operations

## Problem Statement
Build a high-performance concurrent counter using sync/atomic instead of a mutex. Compare throughput with a mutex-based counter using benchmarks. Use atomic.Int64 (Go 1.19+).

## Key Concepts
```
sync/atomic, atomic.Int64, LoadInt64, AddInt64
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
