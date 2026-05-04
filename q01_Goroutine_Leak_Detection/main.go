package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// TODO: leakedWorker blocks forever on ch — goroutine leak
func leakedWorker(ch <-chan int) {
	// TODO: receive from ch (blocks forever → leak)
}

// TODO: fixedWorker must exit when ctx is cancelled
func fixedWorker(ctx context.Context, ch <-chan int) {
	// TODO: use select with ctx.Done() and ch
}

func main() {
	fmt.Println("goroutines before:", runtime.NumGoroutine()) // expected: 1

	ch := make(chan int)
	go leakedWorker(ch)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("goroutines after leak:", runtime.NumGoroutine()) // expected: 2

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go fixedWorker(ctx, make(chan int))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("goroutines after fix:", runtime.NumGoroutine()) // expected: 2 (leak remains) fixed=1
}
