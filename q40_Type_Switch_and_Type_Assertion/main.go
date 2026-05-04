package main

import "fmt"

func describe(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d (doubled=%d)", v, v*2)
	case float64:
		return fmt.Sprintf("float64: %.2f", v)
	case string:
		return fmt.Sprintf("string: %q (len=%d)", v, len(v))
	case []int:
		sum := 0; for _, x := range v { sum += x }
		return fmt.Sprintf("[]int: len=%d sum=%d", len(v), sum)
	case map[string]int:
		return fmt.Sprintf("map[string]int: len=%d", len(v))
	case bool:
		return fmt.Sprintf("bool: %v", v)
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

func typeAssertionDemo() {
	var i interface{} = "hello"

	// Safe assertion (comma-ok idiom)
	s, ok := i.(string)
	fmt.Printf(".(string): %q ok=%v\n", s, ok)

	n, ok := i.(int)
	fmt.Printf(".(int):    %d  ok=%v\n", n, ok)

	// Unsafe assertion — panics if wrong type
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic:", r)
		}
	}()
	_ = i.(int) // panics: interface holds string not int
}

func main() {
	values := []interface{}{42, 3.14, "go", []int{1, 2, 3}, map[string]int{"a": 1}, true, nil}
	for _, v := range values {
		fmt.Println(describe(v))
	}
	fmt.Println()
	typeAssertionDemo()
}
