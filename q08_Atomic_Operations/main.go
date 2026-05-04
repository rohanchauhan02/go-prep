package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func mutexCounter(n int) int64 {
	var mu sync.Mutex
	var count int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock(); count++; mu.Unlock()
		}()
	}
	wg.Wait()
	return count
}

func atomicCounter(n int) int64 {
	var count atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count.Add(1)
		}()
	}
	wg.Wait()
	return count.Load()
}

func bench(name string, fn func(int) int64) {
	start := time.Now()
	result := fn(100000)
	fmt.Printf("%s: count=%d time=%v\n", name, result, time.Since(start))
}

func main() {
	bench("Mutex  counter", mutexCounter)
	bench("Atomic counter", atomicCounter)
	fmt.Println("Atomic is faster — no lock/unlock overhead")
}
