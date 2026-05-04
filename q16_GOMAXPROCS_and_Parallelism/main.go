//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// TODO: cpuWork does heavy computation (no IO) — sum of squares up to n
func cpuWork(n int) int {
	// TODO: implement sum of i*i for i in 0..n
	return 0
}

// TODO: run `workers` goroutines each calling cpuWork(5_000_000)
// set GOMAXPROCS to `procs` before running, return elapsed time
func runParallel(procs, workers int) time.Duration {
	runtime.GOMAXPROCS(procs)
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: call cpuWork
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	cpus := runtime.NumCPU()
	fmt.Println("CPUs:", cpus)

	t1 := runParallel(1, 4)
	tN := runParallel(cpus, 4)

	fmt.Printf("GOMAXPROCS=1   : %v
", t1)
	fmt.Printf("GOMAXPROCS=%-2d  : %v
", cpus, tN)
	fmt.Printf("Speedup: %.2fx
", float64(t1)/float64(tN))
	// Expected: ~Nx speedup where N = num CPUs
}
