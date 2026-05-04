package main

import (
	"fmt"
	"sync"
	"time"
)

// ─── Barrier using sync.Cond ─────────────────────
type Barrier struct {
	mu      sync.Mutex
	cond    *sync.Cond
	count   int
	total   int
	generation int
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{total: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	gen := b.generation
	b.count++
	if b.count == b.total {
		b.count = 0
		b.generation++
		b.cond.Broadcast() // release all waiters
	} else {
		for gen == b.generation { b.cond.Wait() }
	}
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
			// Phase 1
			time.Sleep(time.Duration(id*30) * time.Millisecond)
			fmt.Printf("  goroutine-%d reached barrier (phase 1)\n", id)
			barrier.Wait()
			// Phase 2 — all start together after barrier
			fmt.Printf("  goroutine-%d past barrier  (phase 2)\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("All goroutines completed both phases")
}
