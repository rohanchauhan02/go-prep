# Q34: Middleware Chain

## Problem Statement
Build a composable HTTP middleware system from scratch (no framework). Implement Logger, Auth, and RateLimit middlewares. Chain using functional pattern: type Middleware func(Handler) Handler.

## Key Concepts
```
http.Handler, http.HandlerFunc, adapter pattern, closure chain
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
