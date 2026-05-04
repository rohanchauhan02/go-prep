# Q21: Memory Escape Analysis

## Problem Statement
Write two versions of a function — one where the variable escapes to the heap, one that stays on the stack. Use go build -gcflags='-m' to verify. Explain the GC pressure implications.

## Key Concepts
```
go build -gcflags='-m', heap vs stack, &local variable
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
