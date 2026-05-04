package main

import "fmt"

func deferOrder() {
	fmt.Println("start")
	defer fmt.Println("defer 1 - runs last")
	defer fmt.Println("defer 2 - runs second")
	defer fmt.Println("defer 3 - runs first")
	fmt.Println("end")
}

func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}()
	result = a / b // panics if b == 0
	return
}

// Named return modified by defer
func namedReturn() (total int) {
	defer func() { total *= 2 }() // modifies named return
	total = 10
	return // returns 20, not 10
}

func main() {
	fmt.Println("=== Defer LIFO Order ===")
	deferOrder()

	fmt.Println("\n=== Panic Recovery ===")
	res, err := safeDiv(10, 2)
	fmt.Printf("10/2 = %d, err = %v\n", res, err)
	res, err = safeDiv(10, 0)
	fmt.Printf("10/0 = %d, err = %v\n", res, err)

	fmt.Println("\n=== Named Return Modified by Defer ===")
	fmt.Println("namedReturn() =", namedReturn()) // 20
}
