# Q25: Publish-Subscribe Event Bus

## Problem Statement
Build an in-memory pub/sub event bus supporting Subscribe(topic), Publish(topic, msg), Unsubscribe. Use channels per subscriber. Handle slow subscribers without blocking publishers.

## Key Concepts
```
map[string][]chan, sync.RWMutex, non-blocking send, select default
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
