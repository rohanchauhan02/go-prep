package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var inFlight sync.WaitGroup

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: write 200 {"status":"ok"}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	inFlight.Add(1)
	defer inFlight.Done()
	// TODO: simulate 2s work, then respond "processed"
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "processed")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/process", processHandler)

	srv := &http.Server{Addr: ":8080", Handler: mux}

	// TODO: start server in goroutine
	go func() {
		fmt.Println("listening :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// TODO: listen for SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down...")

	// TODO: give in-flight requests 5s to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("forced:", err)
	}
	inFlight.Wait()
	fmt.Println("clean exit")
}
