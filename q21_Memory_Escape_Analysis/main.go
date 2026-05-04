package main

import "fmt"

// Stays on stack — small value, doesn't escape
func stackAlloc() int {
	x := 42 // stays on stack
	return x
}

// Escapes to heap — pointer returned, value outlives function
func heapAlloc() *int {
	x := 42
	return &x // x escapes to heap
}

// Interface boxing causes escape
func interfaceEscape() interface{} {
	x := 42
	return x // x escapes: stored in interface
}

// Large array may escape
func largeArray() [1024]int {
	var arr [1024]int
	for i := range arr { arr[i] = i }
	return arr // returned by value — may or may not escape
}

// Goroutine closure causes escape
func goroutineEscape() {
	x := 42
	go func() { fmt.Println(x) }() // x escapes — captured by goroutine
}

func main() {
	fmt.Println("stack alloc result:", stackAlloc())
	fmt.Println("heap alloc result:", *heapAlloc())
	fmt.Println("interface escape:", interfaceEscape())

	arr := largeArray()
	fmt.Println("large array[0]:", arr[0])

	goroutineEscape()

	// Run: go build -gcflags='-m' main.go
	// to see escape analysis output from the compiler
	fmt.Println("\nRun: go build -gcflags='-m -m' main.go to see escape analysis")
}
