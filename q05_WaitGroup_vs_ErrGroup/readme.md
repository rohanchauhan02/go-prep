# Q05: WaitGroup vs ErrGroup

## Problem Statement
Download 5 URLs concurrently (mock with time.Sleep). First version uses sync.WaitGroup. Second uses golang.org/x/sync/errgroup and stops all goroutines on the first error. Compare both.

## Key Concepts
```
sync.WaitGroup, errgroup.WithContext, early cancellation
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
