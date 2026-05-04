//go:build ignore
// Remove the above line when implementing

package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct{ ID int }

// TODO: workerPool creates numWorkers goroutines reading from jobs channel
// Each worker sleeps 20ms, marks job done or failed (randomly or always ok)
// Returns a results channel (Job IDs processed)
// Must stop cleanly when ctx is cancelled or jobs channel is closed
func workerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan int {
	out := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobs:
					if !ok { return }
					time.Sleep(20 * time.Millisecond)
					// TODO: send job.ID to out
					fmt.Printf("  worker-%d done job-%d
", id, job.ID)
					out <- job.ID
				case <-ctx.Done():
					return
				}
			}
		}(i + 1)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 10)
	for i := 1; i <= 10; i++ { jobs <- Job{ID: i} }
	close(jobs)

	results := workerPool(ctx, 3, jobs)

	var count atomic.Int64
	for range results { count.Add(1) }
	fmt.Println("Total processed:", count.Load()) // Expected: 10
}
