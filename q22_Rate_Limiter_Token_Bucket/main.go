//go:build ignore
// Remove the above line when implementing

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
	return &TokenBucket{tokens: capacity, capacity: capacity, rate: rate, last: time.Now()}
}

// TODO: refill adds tokens based on time elapsed since last call
func (tb *TokenBucket) refill() {
	// TODO: elapsed := time.Since(tb.last).Seconds()
	// TODO: tb.tokens = min(capacity, tokens + elapsed*rate)
	// TODO: update tb.last
}

// TODO: Allow returns true and consumes 1 token if available
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock(); defer tb.mu.Unlock()
	// TODO: refill, check tokens >= 1, decrement
	return false
}

// TODO: Wait blocks until a token is available or ctx is cancelled
func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		if tb.Allow() { return nil }
		select {
		case <-ctx.Done(): return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func main() {
	rl := NewTokenBucket(3, 1) // 3 capacity, 1 token/sec

	for i := 0; i < 5; i++ {
		if rl.Allow() {
			fmt.Printf("request %d: ALLOWED
", i+1)
		} else {
			fmt.Printf("request %d: DENIED
", i+1)
		}
	}
	// Expected: first 3 ALLOWED, then DENIED

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rl.Wait(ctx); err != nil {
		fmt.Println("Wait error:", err)
	} else {
		fmt.Println("Got token after wait")
	}
}
