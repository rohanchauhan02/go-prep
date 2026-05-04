# Q20: Reflection and Struct Tag Parser

## Problem Statement
Write a mini JSON-like serializer using the reflect package that reads struct field tags (json:"name") and serializes a struct to map[string]any. Handle nested structs and pointer fields.

## Key Concepts
```
reflect.TypeOf, reflect.ValueOf, StructTag.Get, Kind()
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
