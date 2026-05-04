//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"sync"
	"time"
)

// TODO: implement Get/Set using sync.Mutex
type MutexCache struct {
	mu    sync.Mutex
	store map[string]string
}

func (c *MutexCache) Get(k string) string {
	// TODO: lock, read, unlock
	return ""
}
func (c *MutexCache) Set(k, v string) {
	// TODO: lock, write, unlock
}

// TODO: implement Get/Set using sync.RWMutex (RLock for reads)
type RWCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func (c *RWCache) Get(k string) string {
	// TODO: RLock for concurrent reads
	return ""
}
func (c *RWCache) Set(k, v string) {
	// TODO: full Lock for writes
}

func bench(name string, fn func(), n int) {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); fn() }()
	}
	wg.Wait()
	fmt.Printf("%s: %v
", name, time.Since(start))
}

func main() {
	mc := &MutexCache{store: map[string]string{"k": "v"}}
	rw := &RWCache{store: map[string]string{"k": "v"}}

	bench("Mutex  ", func() { mc.Get("k") }, 10000)
	bench("RWMutex", func() { rw.Get("k") }, 10000)
	// Expected: RWMutex is faster on read-heavy workload
}
