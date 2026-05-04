package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var urls = []string{"url1", "url2", "url3", "url4", "url5"}

// TODO: fetch each url concurrently using sync.WaitGroup
// Print "fetched: <url>" for each, wait for all to finish
func fetchWaitGroup() {
	var wg sync.WaitGroup
	// TODO: implement
	wg.Wait()
}

// TODO: fetch using errgroup — stop all on first error
// url3 should return an error
func fetchErrGroup() error {
	g, ctx := errgroup.WithContext(context.Background())
	_ = ctx
	for _, u := range urls {
		u := u
		g.Go(func() error {
			time.Sleep(50 * time.Millisecond)
			// TODO: return error if u == "url3"
			fmt.Println("fetched:", u)
			return nil
		})
	}
	return g.Wait()
}

func main() {
	fmt.Println("=== WaitGroup ===")
	fetchWaitGroup()

	fmt.Println("=== ErrGroup ===")
	if err := fetchErrGroup(); err != nil {
		fmt.Println("stopped:", err) // Expected: stopped: url3 failed
	}
	_ = errors.New // hint
}
