package main

import (
	"fmt"
	"sync"
	"time"
)

type MutexCache struct {
	mu    sync.Mutex
	store map[string]string
}

func (c *MutexCache) Get(k string) string {
	c.mu.Lock(); defer c.mu.Unlock()
	return c.store[k]
}
func (c *MutexCache) Set(k, v string) {
	c.mu.Lock(); defer c.mu.Unlock()
	c.store[k] = v
}

type RWCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func (c *RWCache) Get(k string) string {
	c.mu.RLock(); defer c.mu.RUnlock()
	return c.store[k]
}
func (c *RWCache) Set(k, v string) {
	c.mu.Lock(); defer c.mu.Unlock()
	c.store[k] = v
}

func bench(name string, get func(), iters int) {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < iters; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); get() }()
	}
	wg.Wait()
	fmt.Printf("%s: %v for %d reads\n", name, time.Since(start), iters)
}

func main() {
	mc := &MutexCache{store: map[string]string{"key": "value"}}
	rw := &RWCache{store: map[string]string{"key": "value"}}

	bench("sync.Mutex  ", func() { mc.Get("key") }, 10000)
	bench("sync.RWMutex", func() { rw.Get("key") }, 10000)
	fmt.Println("RWMutex wins on read-heavy workloads (multiple readers at once)")
}
