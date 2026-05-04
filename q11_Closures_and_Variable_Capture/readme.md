# Q11: Closures and Variable Capture

## Problem Statement
Show the classic for-loop goroutine closure bug where all goroutines print the same value. Fix it two ways: (1) pass the variable as a function argument, (2) shadow with a new variable inside the loop.

## Key Concepts
```
closure variable capture, goroutine scheduling, loop variable
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
