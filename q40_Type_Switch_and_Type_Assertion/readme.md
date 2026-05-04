# Q40: Type Switch and Type Assertion

## Problem Statement
Write a polymorphic Evaluate function accepting interface{} using type switch to handle int, float64, string, []int, map[string]int. Show panic from failed type assertion vs comma-ok idiom.

## Key Concepts
```
switch v := i.(type), x.(T), x ok := i.(T), comma-ok
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
