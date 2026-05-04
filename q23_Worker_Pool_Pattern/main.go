package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	OK    bool
}

func workerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobs:
					if !ok { return }
					time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
					ok2 := rand.Intn(5) != 0
					results <- Result{JobID: job.ID, OK: ok2}
					fmt.Printf("  worker-%d processed job-%d ok=%v\n", id, job.ID, ok2)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	go func() { wg.Wait(); close(results) }()
	return results
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 20)
	for i := 1; i <= 10; i++ { jobs <- Job{ID: i} }
	close(jobs)

	results := workerPool(ctx, 3, jobs)

	var done, failed atomic.Int64
	for r := range results {
		if r.OK { done.Add(1) } else { failed.Add(1) }
	}
	fmt.Printf("\nDone: %d  Failed: %d\n", done.Load(), failed.Load())
}
