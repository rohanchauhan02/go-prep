# Q44: Compile-Time Interface Check

## Problem Statement
Show the Go idiom for compile-time interface satisfaction: var _ io.Writer = (*MyWriter)(nil). Implement a type satisfying multiple interfaces and add compile-time checks for each.

## Key Concepts
```
blank identifier, var _ Interface = (*Type)(nil)
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
