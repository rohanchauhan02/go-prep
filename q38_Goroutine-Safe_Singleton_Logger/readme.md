# Q38: Goroutine-Safe Singleton Logger

## Problem Statement
Build a structured logger singleton initialized with sync.Once. Support log levels (DEBUG/INFO/WARN/ERROR), output to stdout and file, include caller info using runtime.Caller.

## Key Concepts
```
sync.Once, log.New, runtime.Caller, io.MultiWriter
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
