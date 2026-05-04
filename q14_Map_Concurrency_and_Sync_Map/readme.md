# Q14: Map Concurrency and Sync Map

## Problem Statement
Show the concurrent map write panic. Fix it first with sync.Mutex, then with sync.Map. Benchmark both. Explain when sync.Map outperforms mutex-protected maps.

## Key Concepts
```
sync.Map, Store/Load/Delete/Range, concurrent map panic
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
