# Q18: Custom Error Types and Wrapping

## Problem Statement
Define a custom ValidationError type with a Code field. Wrap it using fmt.Errorf with %w. Use errors.Is and errors.As to unwrap and inspect errors at the call site. Show the error chain.

## Key Concepts
```
errors.Is, errors.As, fmt.Errorf %w, Unwrap() method
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
