package main

import "fmt"

// TODO: stackAlloc — return int by value (stays on stack)
func stackAlloc() int {
	x := 42
	return x // x stays on stack
}

// TODO: heapAlloc — return pointer (x escapes to heap)
func heapAlloc() *int {
	x := 42
	return &x // TODO: why does x escape?
}

// TODO: interfaceEscape — boxing into interface causes escape
func interfaceEscape() interface{} {
	x := 42
	return x // TODO: why does x escape here?
}

func main() {
	fmt.Println(stackAlloc())      // 42
	fmt.Println(*heapAlloc())      // 42
	fmt.Println(interfaceEscape()) // 42

	// Run to see escape decisions:
	// go build -gcflags="-m -m" main.go
	// Look for: "x escapes to heap" vs "x does not escape"
}
