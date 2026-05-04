//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// TODO: increment count 100000 times concurrently using sync.Mutex
func mutexCounter(n int) int64 {
	var mu sync.Mutex
	var count int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: lock, increment count, unlock
		}()
	}
	wg.Wait()
	return count
}

// TODO: increment 100000 times using atomic.Int64 — no mutex needed
func atomicCounter(n int) int64 {
	var count atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: count.Add(1)
		}()
	}
	wg.Wait()
	return count.Load()
}

func main() {
	start := time.Now()
	fmt.Println("mutex :", mutexCounter(100000), "time:", time.Since(start))

	start = time.Now()
	fmt.Println("atomic:", atomicCounter(100000), "time:", time.Since(start))
	// Expected: both = 100000, atomic is faster
}
