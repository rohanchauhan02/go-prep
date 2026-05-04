//go:build ignore
// Remove the above line when implementing

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO: producer sends 1..10 to output channel, respects ctx
func producer(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// TODO: send 1..10, select on ctx.Done()
	}()
	return out
}

// TODO: worker reads from in, sends "worker-<id> processed <v>" strings to output
func worker(ctx context.Context, id int, in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			time.Sleep(20 * time.Millisecond)
			select {
			case out <- fmt.Sprintf("worker-%d processed %d", id, v):
			case <-ctx.Done(): return
			}
		}
	}()
	return out
}

// TODO: merge fans-in multiple channels into one
func merge(ctx context.Context, channels ...<-chan string) <-chan string {
	merged := make(chan string, len(channels))
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		ch := ch
		go func() {
			defer wg.Done()
			// TODO: forward all values from ch to merged
		}()
	}
	go func() { wg.Wait(); close(merged) }()
	return merged
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	src := producer(ctx)

	// Fan-out to 3 workers
	var workerChans []<-chan string
	for i := 1; i <= 3; i++ {
		workerChans = append(workerChans, worker(ctx, i, src))
	}

	// Fan-in
	for result := range merge(ctx, workerChans...) {
		fmt.Println(result)
	}
	// Expected: 10 lines like "worker-X processed Y"
}
