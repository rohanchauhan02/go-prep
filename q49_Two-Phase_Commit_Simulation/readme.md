# Q49: Two-Phase Commit Simulation

## Problem Statement
Simulate a distributed two-phase commit with a Coordinator and 3 Participants. Each participant votes Commit or Abort. Coordinator decides based on unanimous votes. Use channels for messaging.

## Key Concepts
```
channels as message passing, goroutines as nodes, coordinator pattern
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
