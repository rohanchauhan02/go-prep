package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func cpuWork(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i * i
	}
	return sum
}

func runParallel(procs, workers int) time.Duration {
	runtime.GOMAXPROCS(procs)
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cpuWork(5_000_000)
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	workers := 4
	cpus := runtime.NumCPU()
	fmt.Printf("CPU cores available: %d\n\n", cpus)

	t1 := runParallel(1, workers)
	fmt.Printf("GOMAXPROCS=1    : %v\n", t1)

	tN := runParallel(cpus, workers)
	fmt.Printf("GOMAXPROCS=%-4d : %v\n", cpus, tN)

	fmt.Printf("\nSpeedup: %.2fx\n", float64(t1)/float64(tN))
	fmt.Println("More GOMAXPROCS = true parallelism on CPU-bound work")
}
