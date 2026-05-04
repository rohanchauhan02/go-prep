# Q16: GOMAXPROCS and Parallelism

## Problem Statement
Write a program demonstrating the effect of GOMAXPROCS on parallel execution. Run CPU-bound work with GOMAXPROCS=1 vs runtime.NumCPU(). Measure execution time difference.

## Key Concepts
```
runtime.GOMAXPROCS, runtime.NumCPU, CPU-bound goroutines
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
