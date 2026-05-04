package main

import "fmt"

// TODO: describe uses type switch to handle multiple types
// int → "int: N (doubled=2N)"
// float64 → "float64: N"
// string → "string: "s" (len=N)"
// []int → "slice: len=N sum=S"
// nil → "nil"
// default → "unknown: <type>"
func describe(i interface{}) string {
	switch v := i.(type) {
	// TODO: add all cases
	default:
		return fmt.Sprintf("unknown: %T", v)
	}
}

func typeAssertDemo() {
	var i interface{} = "hello"

	// Safe (comma-ok)
	s, ok := i.(string)
	fmt.Println("string:", s, ok) // Expected: hello true

	n, ok := i.(int)
	fmt.Println("int:", n, ok)    // Expected: 0 false

	// Unsafe — catches panic with recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic caught:", r)
		}
	}()
	_ = i.(int) // Expected: panic caught: interface conversion: ...
}

func main() {
	vals := []interface{}{42, 3.14, "go", []int{1,2,3}, nil, true}
	for _, v := range vals {
		fmt.Println(describe(v))
	}
	fmt.Println()
	typeAssertDemo()
}
