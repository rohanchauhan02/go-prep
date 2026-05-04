package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Semaphore struct{ ch chan struct{} }

func NewSemaphore(n int) *Semaphore { return &Semaphore{ch: make(chan struct{}, n)} }

func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}: return nil
	case <-ctx.Done(): return ctx.Err()
	}
}

func (s *Semaphore) Release() { <-s.ch }

// Connection pool backed by semaphore
type Pool struct {
	sem *Semaphore
	id  int
	mu  sync.Mutex
}

func NewPool(size int) *Pool { return &Pool{sem: NewSemaphore(size)} }

func (p *Pool) Do(ctx context.Context, task func(connID int)) error {
	if err := p.sem.Acquire(ctx); err != nil { return err }
	defer p.sem.Release()
	p.mu.Lock(); id := p.id; p.id++; p.mu.Unlock()
	task(id)
	return nil
}

func main() {
	pool := NewPool(3) // max 3 concurrent
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := pool.Do(ctx, func(connID int) {
				fmt.Printf("goroutine-%d acquired conn-%d\n", i, connID)
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("goroutine-%d released conn-%d\n", i, connID)
			})
			if err != nil {
				fmt.Printf("goroutine-%d: %v\n", i, err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("All done")
}
