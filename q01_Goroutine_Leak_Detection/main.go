package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// LEAKED version - goroutine blocks forever on channel with no sender
func leakedWorker(ch <-chan int) {
	val := <-ch // blocks forever → leak
	fmt.Println("got", val)
}

// FIXED version - respects context cancellation
func fixedWorker(ctx context.Context, ch <-chan int) {
	select {
	case val := <-ch:
		fmt.Println("got", val)
	case <-ctx.Done():
		fmt.Println("worker cancelled:", ctx.Err())
	}
}

func main() {
	fmt.Println("=== Goroutine Leak Demo ===")
	before := runtime.NumGoroutine()
	fmt.Println("goroutines before:", before)

	// Create leak
	ch := make(chan int)
	go leakedWorker(ch)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("goroutines after leak:", runtime.NumGoroutine()) // +1

	// Fix with context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go fixedWorker(ctx, make(chan int))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("goroutines after fix:", runtime.NumGoroutine()) // back to normal
}
