package main

import (
	"bytes"
	"fmt"
	"io"
)

// Fine-grained interfaces
type Reader interface{ Read(p []byte) (int, error) }
type Writer interface{ Write(p []byte) (int, error) }
type Closer interface{ Close() error }

// Composed interface
type ReadWriteCloser interface{ Reader; Writer; Closer }

// TODO: BufferRWC adapts bytes.Buffer to satisfy ReadWriteCloser
// Close() should reset the buffer
type BufferRWC struct {
	buf bytes.Buffer
}

func (b *BufferRWC) Write(p []byte) (int, error) {
	// TODO: delegate to b.buf
	return 0, nil
}
func (b *BufferRWC) Read(p []byte) (int, error) {
	// TODO: delegate to b.buf
	return 0, nil
}
func (b *BufferRWC) Close() error {
	// TODO: reset buffer
	return nil
}

// Compile-time check
var _ ReadWriteCloser = (*BufferRWC)(nil)
var _ io.ReadWriteCloser = (*BufferRWC)(nil)

func process(rwc ReadWriteCloser) {
	defer rwc.Close()
	n, _ := rwc.Write([]byte("Hello, Adapter!"))
	fmt.Println("wrote:", n, "bytes")
	buf := make([]byte, 64)
	n, _ = rwc.Read(buf)
	fmt.Println("read:", string(buf[:n]))
}

func main() {
	process(&BufferRWC{})
}
