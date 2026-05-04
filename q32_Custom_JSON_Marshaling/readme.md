# Q32: Custom JSON Marshaling

## Problem Statement
Define a type wrapping time.Time that marshals to/from Unix timestamp (int64) in JSON. Implement MarshalJSON and UnmarshalJSON. Handle null/zero time edge cases.

## Key Concepts
```
json.Marshaler, json.Unmarshaler, json.RawMessage, omitempty
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
