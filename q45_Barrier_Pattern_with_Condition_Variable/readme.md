# Q45: Barrier Pattern with Condition Variable

## Problem Statement
Implement a reusable barrier that blocks N goroutines until all have reached it, then releases all simultaneously. Implement using sync.WaitGroup first, then using sync.Cond.

## Key Concepts
```
sync.Cond, Broadcast(), Wait(), sync.WaitGroup
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
