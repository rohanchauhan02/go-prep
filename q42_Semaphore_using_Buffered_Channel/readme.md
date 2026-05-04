# Q42: Semaphore using Buffered Channel

## Problem Statement
Implement a semaphore primitive using a buffered channel to limit concurrent access. Build a connection pool using the semaphore. Test with more goroutines than pool size.

## Key Concepts
```
make(chan struct{}, n), acquire=send, release=receive
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
