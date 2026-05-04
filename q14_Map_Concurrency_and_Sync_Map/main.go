//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"sync"
)

// TODO: thread-safe map using sync.Mutex
type SafeMap struct {
	mu sync.Mutex
	m  map[string]int
}

func NewSafeMap() *SafeMap { return &SafeMap{m: make(map[string]int)} }

func (sm *SafeMap) Set(k string, v int) {
	// TODO: lock, write, unlock
}

func (sm *SafeMap) Get(k string) (int, bool) {
	// TODO: lock, read, unlock
	return 0, false
}

// TODO: use sync.Map — Store, Load, Range
func syncMapDemo() {
	var m sync.Map
	// TODO: Store 3 key-value pairs
	// TODO: Load one key and print it
	// TODO: Range over all and print
}

func main() {
	sm := NewSafeMap()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(fmt.Sprintf("k%d", i), i)
		}(i)
	}
	wg.Wait()
	v, ok := sm.Get("k5")
	fmt.Println("k5:", v, ok) // Expected: 5 true

	fmt.Println("
sync.Map demo:")
	syncMapDemo()
}
