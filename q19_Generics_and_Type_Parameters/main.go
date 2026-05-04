package main

import "fmt"

// TODO: generic Stack[T any] with Push, Pop (T, bool), Peek (T, bool), IsEmpty
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	// TODO: append v
}

func (s *Stack[T]) IsEmpty() bool {
	// TODO: return len == 0
	return true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	// TODO: return last item or (zero, false)
	return zero, false
}

func (s *Stack[T]) Pop() (T, bool) {
	// TODO: peek then shrink slice
	var zero T
	return zero, false
}

// TODO: Map transforms []T → []U using fn
func Map[T, U any](slice []T, fn func(T) U) []U {
	// TODO: implement
	return nil
}

// TODO: Filter returns elements where fn returns true
func Filter[T any](slice []T, fn func(T) bool) []T {
	// TODO: implement
	return nil
}

// TODO: Number constraint — Sum works for int and float64
type Number interface{ ~int | ~float64 }

func Sum[T Number](nums []T) T {
	var total T
	// TODO: sum all
	return total
}

func main() {
	s := &Stack[int]{}
	s.Push(1); s.Push(2); s.Push(3)

	v, _ := s.Pop()
	fmt.Println("Pop:", v)   // Expected: 3

	top, _ := s.Peek()
	fmt.Println("Peek:", top) // Expected: 2

	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Map x2:", Map(nums, func(n int) int { return n * 2 })) // [2 4 6 8 10]
	fmt.Println("Filter even:", Filter(nums, func(n int) bool { return n%2 == 0 })) // [2 4]
	fmt.Println("Sum:", Sum(nums)) // 15
}
