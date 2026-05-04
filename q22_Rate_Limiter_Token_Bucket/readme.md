# Q22: Rate Limiter Token Bucket

## Problem Statement
Implement a thread-safe token bucket rate limiter from scratch. Support Allow() bool and Wait(ctx) error methods. Then solve the same using golang.org/x/time/rate.

## Key Concepts
```
time.Ticker, sync.Mutex, rate.Limiter, context timeout
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
