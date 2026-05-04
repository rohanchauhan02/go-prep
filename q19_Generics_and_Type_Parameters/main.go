package main

import "fmt"

// Generic Stack
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T)       { s.items = append(s.items, v) }
func (s *Stack[T]) IsEmpty() bool  { return len(s.items) == 0 }
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if s.IsEmpty() { return zero, false }
	return s.items[len(s.items)-1], true
}
func (s *Stack[T]) Pop() (T, bool) {
	v, ok := s.Peek()
	if ok { s.items = s.items[:len(s.items)-1] }
	return v, ok
}

// Generic Map (functional transform)
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice { result[i] = fn(v) }
	return result
}

// Generic Filter
func Filter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice { if fn(v) { result = append(result, v) } }
	return result
}

// Constraint with ~
type Number interface{ ~int | ~float64 }

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums { total += n }
	return total
}

func main() {
	// Stack[int]
	s := &Stack[int]{}
	s.Push(1); s.Push(2); s.Push(3)
	v, _ := s.Pop()
	fmt.Println("Stack Pop:", v)
	top, _ := s.Peek()
	fmt.Println("Stack Peek:", top)

	// Stack[string]
	ss := &Stack[string]{}
	ss.Push("go"); ss.Push("generics")
	fmt.Println("String stack pop:", func() string { v, _ := ss.Pop(); return v }())

	// Map
	nums := []int{1, 2, 3, 4, 5}
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("Map doubled:", doubled)

	strs := Map(nums, func(n int) string { return fmt.Sprintf("item%d", n) })
	fmt.Println("Map to string:", strs)

	// Filter
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("Filter evens:", evens)

	// Sum with constraint
	fmt.Println("Sum ints:", Sum(nums))
	fmt.Println("Sum floats:", Sum([]float64{1.1, 2.2, 3.3}))
}
