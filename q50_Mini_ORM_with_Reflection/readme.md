# Q50: Mini ORM with Reflection

## Problem Statement
Build a minimal ORM using reflection to: (1) generate CREATE TABLE SQL from a struct with db tags, (2) generate INSERT SQL with values, (3) scan sql.Rows into a struct slice. No external packages.

## Key Concepts
```
reflect.TypeOf, StructTag.Get("db"), sql.Rows.Scan, variadic args
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
