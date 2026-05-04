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

func (eb *EventBus) Subscribe(topic string) <-chan string {
	ch := make(chan string, 10)
	eb.mu.Lock()
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	eb.mu.Unlock()
	return ch
}

func (eb *EventBus) Publish(topic, msg string) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, ch := range eb.subscribers[topic] {
		select {
		case ch <- msg: // non-blocking — slow subscribers are skipped
		default:
			fmt.Printf("  [bus] subscriber slow, dropped: %s\n", msg)
		}
	}
}

func (eb *EventBus) Unsubscribe(topic string, target <-chan string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	subs := eb.subscribers[topic]
	for i, ch := range subs {
		if ch == target {
			close(ch)
			eb.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			return
		}
	}
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
			for msg := range ch {
				fmt.Printf("  [%s] received: %s\n", name, msg)
			}
		}()
	}
	listen("sub1-orders", ch1)
	listen("sub2-orders", ch2)
	listen("sub3-payments", ch3)

	bus.Publish("orders", "order-001 placed")
	bus.Publish("orders", "order-002 placed")
	bus.Publish("payments", "payment-001 received")

	time.Sleep(100 * time.Millisecond)
	bus.Unsubscribe("orders", ch1)
	bus.Unsubscribe("orders", ch2)
	bus.Unsubscribe("payments", ch3)
	wg.Wait()
	fmt.Println("Done")
}
