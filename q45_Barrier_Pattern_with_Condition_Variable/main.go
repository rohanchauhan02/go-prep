//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"sync"
	"time"
)

// TODO: Barrier blocks N goroutines until all have arrived, then releases all
// Use sync.Cond for implementation
type Barrier struct {
	mu         sync.Mutex
	cond       *sync.Cond
	count      int
	total      int
	generation int
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{total: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// TODO: Wait increments count, if count==total broadcast to release all
// otherwise Wait on cond until generation advances
func (b *Barrier) Wait() {
	b.mu.Lock()
	// TODO: implement
	b.mu.Unlock()
}

func main() {
	n := 5
	barrier := NewBarrier(n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*30) * time.Millisecond)
			fmt.Printf("goroutine-%d reached barrier
", id)
			barrier.Wait()
			// All goroutines cross together
			fmt.Printf("goroutine-%d past barrier
", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("all phases complete")
	// Expected: all "reached" printed, then all "past" printed (roughly together)
}
