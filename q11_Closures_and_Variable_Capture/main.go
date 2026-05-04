package main

import (
	"fmt"
	"sync"
)

// TODO: buggyClosures returns 5 funcs — all print the same wrong value (5)
// because they all capture the same loop variable i
func buggyClosures() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		funcs[i] = func() { fmt.Print(i, " ") } // BUG: captures &i
	}
	return funcs
}

// TODO: fix by shadowing i inside the loop with i := i
func fixedByShadow() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		// TODO: shadow i so each closure captures its own copy
		funcs[i] = func() { fmt.Print(i, " ") }
	}
	return funcs
}

// TODO: fix by passing i as a function argument
func fixedByParam() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		// TODO: wrap in an immediately-invoked func that takes i as param
		funcs[i] = func() { fmt.Print(i, " ") }
	}
	return funcs
}

// TODO: goroutines — fix the same bug with i := i before go func()
func fixedGoroutines() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		// TODO: fix closure capture before go func()
		go func() { defer wg.Done(); fmt.Print(i, " ") }()
	}
	wg.Wait()
}

func main() {
	fmt.Println("Buggy (all 5):")
	for _, f := range buggyClosures() { f() } // Expected: 5 5 5 5 5
	fmt.Println()

	fmt.Println("Fixed shadow:")
	for _, f := range fixedByShadow() { f() } // Expected: 0 1 2 3 4
	fmt.Println()

	fmt.Println("Fixed param:")
	for _, f := range fixedByParam() { f() }  // Expected: 0 1 2 3 4
	fmt.Println()

	fmt.Println("Fixed goroutines:")
	fixedGoroutines()                          // Expected: 0 1 2 3 4 (any order)
	fmt.Println()
}
