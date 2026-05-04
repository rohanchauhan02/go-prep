package main

import (
	"fmt"
	"io"
)

type MyWriter struct{ buf []byte }

// TODO: implement io.Writer
func (w *MyWriter) Write(p []byte) (int, error) {
	// TODO: append p to w.buf
	return 0, nil
}

type MyReadWriter struct {
	MyWriter
	data []byte
	pos  int
}

// TODO: implement io.Reader
func (rw *MyReadWriter) Read(p []byte) (int, error) {
	// TODO: copy data[pos:] into p, advance pos, return io.EOF when done
	return 0, io.EOF
}

// TODO: implement io.Closer
func (rw *MyReadWriter) Close() error {
	// TODO: reset buf and pos
	return nil
}

// Compile-time interface satisfaction — these will fail to compile if methods are missing
var _ io.Writer          = (*MyWriter)(nil)
var _ io.Reader          = (*MyReadWriter)(nil)
var _ io.ReadWriteCloser = (*MyReadWriter)(nil)

func main() {
	w := &MyWriter{}
	fmt.Fprintf(w, "hello %s", "world")
	fmt.Println("buf:", string(w.buf)) // Expected: hello world

	rw := &MyReadWriter{data: []byte("Go is great")}
	buf := make([]byte, 5)
	n, _ := rw.Read(buf)
	fmt.Println("read:", string(buf[:n])) // Expected: Go is
	rw.Close()
	fmt.Println("pos after close:", rw.pos) // Expected: 0
}
