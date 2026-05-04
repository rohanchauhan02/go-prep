package main

import (
	"errors"
	"fmt"
	"time"
)

// TODO: mockService returns a channel that sends name after delay
func mockService(name string, delay time.Duration) <-chan string {
	// TODO: implement
	return nil
}

// TODO: queryWithTimeout returns first response from svc-A(300ms) or svc-B(150ms)
// Returns error if neither responds within 2 seconds
func queryWithTimeout() (string, error) {
	// TODO: use select with time.After(2s)
	return "", errors.New("not implemented")
}

func main() {
	res, err := queryWithTimeout()
	fmt.Println(res, err)
	// Expected: "service-B response" <nil>  (B is faster at 150ms)

	// Timeout scenario
	svc := mockService("slow", 3*time.Second)
	select {
	case r := <-svc:
		fmt.Println(r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Timed out") // Expected: Timed out
	}
}
