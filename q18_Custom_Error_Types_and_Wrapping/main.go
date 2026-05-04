package main

import (
	"errors"
	"fmt"
)

// TODO: define ValidationError with Code int, Field string, Message string
// implement the error interface
type ValidationError struct {
	Code    int
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	// TODO: return formatted string
	return ""
}

var ErrNotFound = errors.New("not found")

// TODO: validate returns *ValidationError if age < 0
func validate(age int) error {
	// TODO: implement
	return nil
}

// TODO: process wraps validate's error using fmt.Errorf("%w")
func process(age int) error {
	// TODO: implement
	return nil
}

// TODO: lookup wraps ErrNotFound if id == 0
func lookup(id int) error {
	// TODO: implement
	return nil
}

func main() {
	err := process(-5)
	fmt.Println(err) // Expected: process failed: validation error [400] field=age: must be non-negative

	var ve *ValidationError
	fmt.Println(errors.As(err, &ve), ve.Code) // Expected: true 400

	err2 := lookup(0)
	fmt.Println(errors.Is(err2, ErrNotFound)) // Expected: true
}
