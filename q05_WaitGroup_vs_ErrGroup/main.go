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

func fetchWaitGroup() {
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			fmt.Println("WaitGroup fetched:", url)
		}(u)
	}
	wg.Wait()
}

func fetchErrGroup() error {
	g, ctx := errgroup.WithContext(context.Background())
	for _, u := range urls {
		u := u
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(50 * time.Millisecond):
			}
			if u == "url3" {
				return errors.New("url3 failed")
			}
			fmt.Println("ErrGroup fetched:", u)
			return nil
		})
	}
	return g.Wait()
}

func main() {
	fmt.Println("=== sync.WaitGroup ===")
	fetchWaitGroup()

	fmt.Println("\n=== errgroup (stops on first error) ===")
	if err := fetchErrGroup(); err != nil {
		fmt.Println("ErrGroup stopped with error:", err)
	}
}
