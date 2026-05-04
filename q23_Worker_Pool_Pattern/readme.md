# Q23: Worker Pool Pattern

## Problem Statement
Build a worker pool with N workers processing jobs from a shared channel. Support graceful shutdown using context cancellation. Track completed and failed job counts atomically.

## Key Concepts
```
buffered channel, sync.WaitGroup, atomic counter, context
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
