package main

import (
	"fmt"
	"sync"
)

// TODO: fixedOrder acquires muA then muB in both goroutines (consistent order)
// This prevents deadlock. Both goroutines must lock in the same order.
func fixedOrder(muA, muB *sync.Mutex, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: lock muA first, then muB, do work, unlock both
}

func main() {
	var muA, muB sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go fixedOrder(&muA, &muB, "goroutine-1", &wg)
	go fixedOrder(&muA, &muB, "goroutine-2", &wg)
	wg.Wait()

	fmt.Println("No deadlock!") // expected: No deadlock!
}
