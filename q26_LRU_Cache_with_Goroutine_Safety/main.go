package main

import (
	"container/list"
	"fmt"
	"sync"
)

type entry struct{ key, value string }

type LRUCache struct {
	mu       sync.Mutex
	capacity int
	list     *list.List
	items    map[string]*list.Element
}

func NewLRU(cap int) *LRUCache {
	return &LRUCache{capacity: cap, list: list.New(), items: make(map[string]*list.Element)}
}

// TODO: Get returns value for key, moves it to front (most recently used)
// returns ("", false) if key not found
func (c *LRUCache) Get(key string) (string, bool) {
	c.mu.Lock(); defer c.mu.Unlock()
	// TODO: implement
	return "", false
}

// TODO: Put inserts or updates key. If at capacity, evict LRU (list.Back())
func (c *LRUCache) Put(key, value string) {
	c.mu.Lock(); defer c.mu.Unlock()
	// TODO: implement
}

func main() {
	cache := NewLRU(3)
	cache.Put("a", "1")
	cache.Put("b", "2")
	cache.Put("c", "3")

	v, ok := cache.Get("a")
	fmt.Println("get a:", v, ok)   // Expected: 1 true

	cache.Put("d", "4")            // evicts b (LRU)

	_, ok = cache.Get("b")
	fmt.Println("get b:", ok)      // Expected: false (evicted)

	v, _ = cache.Get("d")
	fmt.Println("get d:", v)       // Expected: 4
}
