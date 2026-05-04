# Q39: Interface Composition and Adapter Pattern

## Problem Statement
Define fine-grained interfaces (Reader, Writer, Closer). Compose them into ReadWriteCloser. Implement an Adapter wrapping bytes.Buffer to satisfy ReadWriteCloser. Demonstrate the pattern.

## Key Concepts
```
interface embedding, adapter struct, method delegation
```

## Requirements
- Solve the problem in `main.go`
- Include example output in comments
- Add unit tests where applicable

## How to Run
```bash
go run main.go
```
