package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

// TODO: this function has the nil interface trap bug
// When fail=false, it returns a typed nil — fix it
func getBuggyError(fail bool) error {
	var err *MyError
	if fail {
		err = &MyError{"something went wrong"}
	}
	return err // BUG: typed nil — interface is NOT nil!
}

// TODO: fix the bug — when fail=false return untyped nil
func getFixedError(fail bool) error {
	// TODO: implement correctly
	return nil
}

func main() {
	// Bug demonstration
	err := getBuggyError(false)
	fmt.Println("buggy err == nil?", err == nil) // Expected: false (trap!)

	// Fixed
	err2 := getFixedError(false)
	fmt.Println("fixed err == nil?", err2 == nil) // Expected: true

	err3 := getFixedError(true)
	fmt.Println("fixed with error:", err3) // Expected: something went wrong
}
