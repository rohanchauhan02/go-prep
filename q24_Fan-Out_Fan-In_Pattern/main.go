package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= 10; i++ {
			select {
			case out <- i:
			case <-ctx.Done(): return
			}
		}
	}()
	return out
}

func worker(ctx context.Context, id int, in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			time.Sleep(20 * time.Millisecond) // simulate work
			select {
			case out <- fmt.Sprintf("worker-%d processed %d", id, v):
			case <-ctx.Done(): return
			}
		}
	}()
	return out
}

func merge(ctx context.Context, channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	merged := make(chan string, len(channels))

	fanIn := func(ch <-chan string) {
		defer wg.Done()
		for v := range ch {
			select {
			case merged <- v:
			case <-ctx.Done(): return
			}
		}
	}
	wg.Add(len(channels))
	for _, ch := range channels { go fanIn(ch) }
	go func() { wg.Wait(); close(merged) }()
	return merged
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	src := producer(ctx)
	// Fan-out to 3 workers (all read same channel — work is distributed)
	var workerChans []<-chan string
	for i := 1; i <= 3; i++ {
		workerChans = append(workerChans, worker(ctx, i, src))
	}
	// Fan-in results
	for result := range merge(ctx, workerChans...) {
		fmt.Println(result)
	}
}
