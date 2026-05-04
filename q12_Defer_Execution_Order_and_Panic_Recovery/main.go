package main

import "fmt"

// TODO: demonstrate defer LIFO — three defers, show order
func deferOrder() {
	fmt.Println("start")
	defer fmt.Println("defer 1") // runs last
	defer fmt.Println("defer 2") // runs second
	defer fmt.Println("defer 3") // runs first
	fmt.Println("end")
	// Expected output order: start, end, defer3, defer2, defer1
}

// TODO: safeDiv catches divide-by-zero panic via recover()
// returns (result int, err error)
func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: set err from recovered value
		}
	}()
	result = a / b
	return
}

// TODO: namedReturn — defer modifies named return value
// return value should be 20 (total=10, defer multiplies by 2)
func namedReturn() (total int) {
	defer func() {
		// TODO: double total
	}()
	total = 10
	return
}

func main() {
	deferOrder()

	r, err := safeDiv(10, 2)
	fmt.Println(r, err)     // Expected: 10 <nil>

	r, err = safeDiv(10, 0)
	fmt.Println(r, err)     // Expected: 0 recovered panic: runtime error: integer divide by zero

	fmt.Println(namedReturn()) // Expected: 20
}
