# Q26: LRU Cache with Goroutine Safety

## Problem Statement
Implement a thread-safe LRU cache with O(1) Get and Put using a doubly linked list + hashmap. Wrap with sync.RWMutex. Write tests including concurrent access tests.

## Key Concepts
```
container/list, map[key]*list.Element, sync.Mutex, eviction
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
