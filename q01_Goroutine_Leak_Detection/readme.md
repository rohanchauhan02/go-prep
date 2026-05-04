# Q01: Goroutine Leak Detection

## Problem Statement
Write a program that deliberately creates a goroutine leak (a goroutine that never terminates), then fix it using context cancellation. Explain why the leak happens and how to detect it using runtime.NumGoroutine().

## Key Concepts
```
context.WithCancel, select with done channel, runtime.NumGoroutine()
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
