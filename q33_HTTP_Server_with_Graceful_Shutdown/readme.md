# Q33: HTTP Server with Graceful Shutdown

## Problem Statement
Build an HTTP server with /health and /process endpoints (simulates long work). Implement graceful shutdown: on SIGINT/SIGTERM stop accepting new requests and wait for in-flight ones.

## Key Concepts
```
http.Server.Shutdown, signal.NotifyContext, sync.WaitGroup
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
