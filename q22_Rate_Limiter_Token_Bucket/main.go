package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // tokens per second
	last     time.Time
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{
		tokens: capacity, capacity: capacity,
		rate: rate, last: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)
	tb.last = now
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock(); defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		if tb.Allow() { return nil }
		select {
		case <-ctx.Done(): return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func min(a, b float64) float64 {
	if a < b { return a }
	return b
}

func main() {
	rl := NewTokenBucket(5, 2) // 5 capacity, 2 tokens/sec

	fmt.Println("=== Token Bucket Rate Limiter ===")
	allowed, denied := 0, 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			fmt.Printf("request %2d: ALLOWED\n", i+1)
			allowed++
		} else {
			fmt.Printf("request %2d: DENIED\n", i+1)
			denied++
		}
	}
	fmt.Printf("Allowed: %d, Denied: %d\n", allowed, denied)

	fmt.Println("\n--- Waiting for token ---")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rl.Wait(ctx); err != nil {
		fmt.Println("Wait error:", err)
	} else {
		fmt.Println("Got token after waiting")
	}
}
