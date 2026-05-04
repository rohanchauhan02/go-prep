package main

import (
	"errors"
	"fmt"
	"time"
)

func mockService(name string, delay time.Duration) <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(delay)
		ch <- name + " response"
	}()
	return ch
}

func queryWithTimeout() (string, error) {
	svc1 := mockService("service-A", 300*time.Millisecond)
	svc2 := mockService("service-B", 150*time.Millisecond)
	timeout := time.After(2 * time.Second)

	select {
	case res := <-svc1:
		return res, nil
	case res := <-svc2:
		return res, nil
	case <-timeout:
		return "", errors.New("both services timed out")
	}
}

func main() {
	res, err := queryWithTimeout()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("First response:", res)
	}

	// Simulate timeout scenario
	fmt.Println("\n--- Simulating slow services ---")
	slow1 := mockService("slow-A", 3*time.Second)
	slow2 := mockService("slow-B", 3*time.Second)
	select {
	case r := <-slow1:
		fmt.Println(r)
	case r := <-slow2:
		fmt.Println(r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Timed out as expected")
	}
}
