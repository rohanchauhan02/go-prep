package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

// BUG: returns typed nil — interface is NOT nil!
func getBuggyError(fail bool) error {
	var err *MyError // typed nil pointer
	if fail {
		err = &MyError{"something went wrong"}
	}
	return err // wraps (*MyError)(nil) — non-nil interface!
}

// FIX: return untyped nil directly
func getFixedError(fail bool) error {
	if fail {
		return &MyError{"something went wrong"}
	}
	return nil // untyped nil — interface IS nil
}

func main() {
	// Buggy
	err := getBuggyError(false)
	fmt.Println("Buggy — err == nil?", err == nil) // false! (trap)
	fmt.Printf("Buggy — type: %T value: %v\n", err, err)

	// Fixed
	err2 := getFixedError(false)
	fmt.Println("Fixed — err == nil?", err2 == nil) // true

	err3 := getFixedError(true)
	fmt.Println("Fixed — error:", err3)
}
