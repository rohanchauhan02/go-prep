package main

import (
	"context"
	"fmt"
)

// TODO: generate sends nums into a send-only channel, closes it when done or ctx cancelled
func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		// TODO: implement
	}()
	return out
}

// TODO: square reads from in, sends n*n to output channel
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// TODO: implement
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := generate(ctx, 1, 2, 3, 4, 5)
	sq := square(ctx, c)

	for v := range sq {
		fmt.Println(v)
	}
	// Expected output:
	// 1
	// 4
	// 9
	// 16
	// 25
}
