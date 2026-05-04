# Q24: Fan-Out Fan-In Pattern

## Problem Statement
Implement fan-out/fan-in: one input channel fans out to 5 worker goroutines, results fan back into a single output channel. Ensure proper cleanup when context is cancelled.

## Key Concepts
```
fan-out, fan-in, merge channels, sync.WaitGroup
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
