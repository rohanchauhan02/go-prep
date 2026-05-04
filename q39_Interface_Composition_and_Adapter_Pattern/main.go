package main

import (
	"bytes"
	"fmt"
	"io"
)

// Fine-grained interfaces (like io package)
type Reader interface{ Read(p []byte) (int, error) }
type Writer interface{ Write(p []byte) (int, error) }
type Closer interface{ Close() error }

// Composed interface
type ReadWriteCloser interface {
	Reader; Writer; Closer
}

// Adapter: wraps bytes.Buffer and adds Close()
type BufferRWC struct{ buf bytes.Buffer }

func (b *BufferRWC) Write(p []byte) (int, error) { return b.buf.Write(p) }
func (b *BufferRWC) Read(p []byte) (int, error)  { return b.buf.Read(p) }
func (b *BufferRWC) Close() error                { b.buf.Reset(); return nil }

// Uses ReadWriteCloser interface
func process(rwc ReadWriteCloser) {
	defer rwc.Close()
	n, _ := rwc.Write([]byte("Hello, Adapter Pattern!"))
	fmt.Printf("Wrote %d bytes\n", n)

	buf := make([]byte, 64)
	n, _ = rwc.Read(buf)
	fmt.Printf("Read  %d bytes: %s\n", n, buf[:n])
}

// Verify at compile time
var _ ReadWriteCloser = (*BufferRWC)(nil)
var _ io.ReadWriteCloser = (*BufferRWC)(nil) // also satisfies stdlib

func main() {
	rwc := &BufferRWC{}
	process(rwc)
	fmt.Println("After Close, buffer is empty:", rwc.buf.Len() == 0)
}
