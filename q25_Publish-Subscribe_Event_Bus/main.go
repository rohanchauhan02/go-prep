//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"sync"
	"time"
)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]chan string)}
}

// TODO: Subscribe creates a buffered channel (cap 10) and registers it for topic
func (eb *EventBus) Subscribe(topic string) <-chan string {
	// TODO: implement
	return nil
}

// TODO: Publish sends msg to all subscribers of topic
// Use non-blocking send (select + default) to skip slow subscribers
func (eb *EventBus) Publish(topic, msg string) {
	eb.mu.RLock(); defer eb.mu.RUnlock()
	// TODO: implement
}

// TODO: Unsubscribe removes a subscriber channel and closes it
func (eb *EventBus) Unsubscribe(topic string, ch <-chan string) {
	// TODO: implement
}

func main() {
	bus := NewEventBus()
	ch1 := bus.Subscribe("orders")
	ch2 := bus.Subscribe("orders")
	ch3 := bus.Subscribe("payments")

	var wg sync.WaitGroup
	listen := func(name string, ch <-chan string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range ch { fmt.Printf("[%s] %s
", name, msg) }
		}()
	}
	listen("sub1", ch1)
	listen("sub2", ch2)
	listen("sub3-pay", ch3)

	bus.Publish("orders", "order-001")
	bus.Publish("payments", "payment-001")
	time.Sleep(50 * time.Millisecond)

	bus.Unsubscribe("orders", ch1)
	bus.Unsubscribe("orders", ch2)
	bus.Unsubscribe("payments", ch3)
	wg.Wait()
	fmt.Println("done")
}
