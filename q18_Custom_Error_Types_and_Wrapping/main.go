package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Code    int
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error [%d] field=%s: %s", e.Code, e.Field, e.Message)
}

var ErrNotFound = errors.New("not found")

func validate(age int) error {
	if age < 0 {
		return &ValidationError{Code: 400, Field: "age", Message: "must be non-negative"}
	}
	return nil
}

func process(age int) error {
	if err := validate(age); err != nil {
		return fmt.Errorf("process failed: %w", err) // wrap with %w
	}
	return nil
}

func lookup(id int) error {
	if id == 0 {
		return fmt.Errorf("lookup id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func main() {
	// Unwrap with errors.As
	err := process(-5)
	fmt.Println("error:", err)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("Code=%d Field=%s\n", ve.Code, ve.Field)
	}

	// Unwrap with errors.Is
	err2 := lookup(0)
	fmt.Println("\nerror:", err2)
	fmt.Println("errors.Is ErrNotFound:", errors.Is(err2, ErrNotFound))

	// Unwrap chain
	fmt.Println("Unwrapped:", errors.Unwrap(err2))
}
