# Q43: Functional Options Pattern

## Problem Statement
Implement a Server struct configurable via functional options: WithTimeout, WithMaxConns, WithTLS. Show why this pattern is preferred over large config structs or builder pattern in Go.

## Key Concepts
```
type Option func(*Server), variadic opts, sane defaults
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
