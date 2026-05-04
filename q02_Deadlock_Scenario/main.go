package main

import (
	"fmt"
	"sync"
	"time"
)

var muA, muB sync.Mutex

// DEADLOCK: goroutine1 locks A then B, goroutine2 locks B then A
func deadlockDemo() {
	done := make(chan struct{})
	go func() {
		muA.Lock(); defer muA.Unlock()
		time.Sleep(10 * time.Millisecond)
		fmt.Println("trying muB from goroutine1...")
		muB.Lock(); defer muB.Unlock()
		close(done)
	}()
	muB.Lock(); defer muB.Unlock()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("trying muA from main...")
	muA.Lock(); defer muA.Unlock()
	<-done
}

// FIX: always acquire locks in the same order (A then B)
func fixedOrder() {
	var wg sync.WaitGroup
	var mA, mB sync.Mutex
	task := func(name string) {
		defer wg.Done()
		mA.Lock(); defer mA.Unlock()
		mB.Lock(); defer mB.Unlock()
		fmt.Println(name, "completed safely")
	}
	wg.Add(2)
	go task("goroutine-1")
	go task("goroutine-2")
	wg.Wait()
}

func main() {
	fmt.Println("=== Deadlock Fix Demo (skipping actual deadlock) ===")
	// deadlockDemo() // uncomment to observe deadlock
	fixedOrder()
	fmt.Println("No deadlock! Consistent lock ordering works.")
}
