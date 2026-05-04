# Q04: Select Statement and Timeout

## Problem Statement
Write a function that queries two mock services concurrently and returns whichever responds first. If neither responds within 2 seconds, return a timeout error. Use select with a time.After case.

## Key Concepts
```
select, time.After, goroutines, channels
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
