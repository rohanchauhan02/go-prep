# Q13: Slice Internals and Append Gotchas

## Problem Statement
Demonstrate how slice header (ptr, len, cap) works. Show the subtle bug where two slices share the same underlying array after slicing causing unexpected mutations. Fix with copy().

## Key Concepts
```
reflect.SliceHeader, append, copy, capacity growth
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
