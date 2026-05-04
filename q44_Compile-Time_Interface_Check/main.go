package main

import (
	"fmt"
	"io"
)

type MyWriter struct{ buf []byte }

func (w *MyWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

type MyReadWriter struct{ MyWriter; data []byte; pos int }

func (rw *MyReadWriter) Read(p []byte) (int, error) {
	if rw.pos >= len(rw.data) { return 0, io.EOF }
	n := copy(p, rw.data[rw.pos:]); rw.pos += n
	return n, nil
}

func (rw *MyReadWriter) Close() error {
	rw.buf = nil; rw.pos = 0; return nil
}

// Compile-time interface satisfaction checks
var (
	_ io.Writer         = (*MyWriter)(nil)
	_ io.Reader         = (*MyReadWriter)(nil)
	_ io.ReadWriteCloser = (*MyReadWriter)(nil)
)

// If any of the above fail, you get a compile error like:
// cannot use (*MyWriter)(nil) as type io.Writer: missing method Write

func useWriter(w io.Writer, msg string) {
	fmt.Fprintf(w, msg)
}

func main() {
	w := &MyWriter{}
	useWriter(w, "Hello, compile-time check!")
	fmt.Println("Written:", string(w.buf))

	rw := &MyReadWriter{data: []byte("Go interfaces are powerful")}
	buf := make([]byte, 10)
	n, _ := rw.Read(buf)
	fmt.Println("Read:", string(buf[:n]))
	rw.Close()
	fmt.Println("Closed, pos reset:", rw.pos)
}
