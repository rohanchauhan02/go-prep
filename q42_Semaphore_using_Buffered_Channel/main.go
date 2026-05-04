//go:build ignore
// Remove the above line when implementing

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO: Semaphore limits concurrent access to n
// Acquire blocks until a slot is free (or ctx cancelled)
// Release frees a slot
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	// TODO: make buffered channel of size n
	return &Semaphore{}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Semaphore) Release() {
	// TODO: read from channel to free a slot
}

func main() {
	sem := NewSemaphore(3) // max 3 concurrent
	ctx := context.WithValue(context.Background(), "k", "v")

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if err := sem.Acquire(tctx); err != nil {
				fmt.Printf("goroutine-%d: %v
", i, err)
				return
			}
			defer sem.Release()
			fmt.Printf("goroutine-%d working
", i)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("goroutine-%d done
", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("all done — max 3 ran concurrently")
}
