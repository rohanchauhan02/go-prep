package main

import (
	"fmt"
	"sync"
	"time"
)

type MutexMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewMutexMap() *MutexMap { return &MutexMap{m: make(map[string]int)} }
func (mm *MutexMap) Set(k string, v int) { mm.mu.Lock(); defer mm.mu.Unlock(); mm.m[k] = v }
func (mm *MutexMap) Get(k string) (int, bool) {
	mm.mu.RLock(); defer mm.mu.RUnlock()
	v, ok := mm.m[k]; return v, ok
}

func benchMutexMap(n int) time.Duration {
	mm := NewMutexMap()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("k%d", i%100)
			mm.Set(k, i)
			mm.Get(k)
		}(i)
	}
	wg.Wait()
	return time.Since(start)
}

func benchSyncMap(n int) time.Duration {
	var sm sync.Map
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("k%d", i%100)
			sm.Store(k, i)
			sm.Load(k)
		}(i)
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	n := 10000
	fmt.Printf("Mutex map :  %v\n", benchMutexMap(n))
	fmt.Printf("sync.Map  :  %v\n", benchSyncMap(n))
	fmt.Println("\nsync.Map wins for write-once/read-many; Mutex wins for frequent writes")

	// Iterate sync.Map
	var sm sync.Map
	sm.Store("go", 1); sm.Store("lang", 2)
	sm.Range(func(k, v any) bool {
		fmt.Printf("  %v = %v\n", k, v)
		return true
	})
}
