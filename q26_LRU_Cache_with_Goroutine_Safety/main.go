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

func (c *LRUCache) Get(key string) (string, bool) {
	c.mu.Lock(); defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.list.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return "", false
}

func (c *LRUCache) Put(key, value string) {
	c.mu.Lock(); defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.list.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}
	if c.list.Len() == c.capacity {
		back := c.list.Back()
		c.list.Remove(back)
		delete(c.items, back.Value.(*entry).key)
	}
	el := c.list.PushFront(&entry{key, value})
	c.items[key] = el
}

func main() {
	cache := NewLRU(3)
	cache.Put("a", "1")
	cache.Put("b", "2")
	cache.Put("c", "3")

	v, ok := cache.Get("a") // a is now most recent
	fmt.Printf("Get a: %s %v\n", v, ok)

	cache.Put("d", "4") // evicts b (LRU)

	_, ok = cache.Get("b")
	fmt.Printf("Get b (evicted): %v\n", ok)

	v, _ = cache.Get("c")
	fmt.Printf("Get c: %s\n", v)

	v, _ = cache.Get("d")
	fmt.Printf("Get d: %s\n", v)

	// Concurrent test
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("k%d", i%5)
			cache.Put(k, fmt.Sprintf("v%d", i))
			cache.Get(k)
		}(i)
	}
	wg.Wait()
	fmt.Println("Concurrent access done safely")
}
