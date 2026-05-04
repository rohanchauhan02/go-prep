# Q03: Channel Direction and Pipeline

## Problem Statement
Build a 3-stage pipeline: generator → square → print, using directional channels (chan<- and <-chan). Each stage runs in its own goroutine. The pipeline must be cancellable via context.

## Key Concepts
```
chan<-, <-chan, context.WithCancel, range over channel
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
