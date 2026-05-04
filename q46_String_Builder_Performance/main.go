//go:build ignore
// Remove the above line when implementing

package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const n = 10000

// TODO: concatenate "a" n times using + operator (slowest)
func concatPlus() string {
	s := ""
	for i := 0; i < n; i++ {
		// TODO: s += "a"
	}
	return s
}

// TODO: use bytes.Buffer
func concatBuffer() string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		// TODO: buf.WriteByte('a')
	}
	return buf.String()
}

// TODO: use strings.Builder with Grow pre-allocation (fastest)
func concatBuilder() string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		// TODO: sb.WriteByte('a')
	}
	return sb.String()
}

// TODO: CSV serializer — join rows with commas, rows with newlines
func toCSV(rows [][]string) string {
	var sb strings.Builder
	// TODO: implement
	return sb.String()
}

func bench(name string, fn func() string) {
	start := time.Now()
	s := fn()
	fmt.Printf("%-20s len=%d time=%v
", name, len(s), time.Since(start))
}

func main() {
	bench("+ operator", concatPlus)
	bench("bytes.Buffer", concatBuffer)
	bench("strings.Builder", concatBuilder)

	rows := [][]string{{"name","age"},{"Alice","30"},{"Bob","25"}}
	fmt.Println(toCSV(rows))
	// Expected:
	// name,age
	// Alice,30
	// Bob,25
}
