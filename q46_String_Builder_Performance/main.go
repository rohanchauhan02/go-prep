package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const n = 10000

func bench(name string, fn func()) time.Duration {
	start := time.Now(); fn(); return time.Since(start)
}

func concatPlus() string {
	s := ""
	for i := 0; i < n; i++ { s += "a" }
	return s
}

func concatSprintf() string {
	s := ""
	for i := 0; i < n; i++ { s = fmt.Sprintf("%sa", s) }
	return s
}

func concatBuffer() string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ { buf.WriteByte('a') }
	return buf.String()
}

func concatBuilder() string {
	var sb strings.Builder
	sb.Grow(n) // pre-allocate
	for i := 0; i < n; i++ { sb.WriteByte('a') }
	return sb.String()
}

// CSV serializer using Builder
func toCSV(rows [][]string) string {
	var sb strings.Builder
	for i, row := range rows {
		for j, cell := range row {
			if j > 0 { sb.WriteByte(',') }
			sb.WriteString(cell)
		}
		if i < len(rows)-1 { sb.WriteByte('
') }
	}
	return sb.String()
}

func main() {
	fmt.Printf("+ operator   : %v\n", bench("plus", func() { concatPlus() }))
	fmt.Printf("fmt.Sprintf  : %v\n", bench("sprintf", func() { concatSprintf() }))
	fmt.Printf("bytes.Buffer : %v\n", bench("buffer", func() { concatBuffer() }))
	fmt.Printf("strings.Builder: %v\n", bench("builder", func() { concatBuilder() }))
	fmt.Println("Winner: strings.Builder (pre-allocated, no extra copy)")

	rows := [][]string{
		{"name", "age", "city"},
		{"Alice", "30", "NYC"},
		{"Bob", "25", "LA"},
	}
	fmt.Println("\nCSV output:")
	fmt.Println(toCSV(rows))
}
