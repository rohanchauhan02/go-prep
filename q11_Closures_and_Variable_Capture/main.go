package main

import (
	"fmt"
	"sync"
)

func buggyClosures() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		funcs[i] = func() { fmt.Print(i, " ") } // captures &i — all print 5
	}
	return funcs
}

func fixedByArg() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		i := i // new variable per iteration
		funcs[i] = func() { fmt.Print(i, " ") }
	}
	return funcs
}

func fixedByParam() []func() {
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		funcs[i] = func(n int) func() {
			return func() { fmt.Print(n, " ") }
		}(i)
	}
	return funcs
}

func goroutineBug() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i // fix: shadow i
		go func() { defer wg.Done(); fmt.Print(i, " ") }()
	}
	wg.Wait()
}

func main() {
	fmt.Println("=== Buggy (all print 5) ===")
	for _, f := range buggyClosures() { f() }

	fmt.Println("\n=== Fixed by shadowing ===")
	for _, f := range fixedByArg() { f() }

	fmt.Println("\n=== Fixed by parameter ===")
	for _, f := range fixedByParam() { f() }

	fmt.Println("\n=== Goroutine fix ===")
	goroutineBug()
	fmt.Println()
}
