# Q17: Context Propagation

## Problem Statement
Build a simulated HTTP middleware chain: AuthMiddleware → LogMiddleware → Handler. Pass request-scoped values (userID, requestID) via context.WithValue. Respect cancellation at each layer.

## Key Concepts
```
context.WithValue, context.WithTimeout, ctx.Value, ctx.Err()
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
