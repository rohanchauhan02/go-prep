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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	inFlight.Add(1)
	defer inFlight.Done()
	time.Sleep(2 * time.Second) // simulate long work
	fmt.Fprintln(w, "processed")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/process", processHandler)

	srv := &http.Server{Addr: ":8080", Handler: mux}

	// Start server
	go func() {
		fmt.Println("Server listening on :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, "ListenAndServe error:", err)
		}
	}()

	// Wait for SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown signal received...")

	// Give in-flight requests 5s to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Forced shutdown:", err)
	}
	inFlight.Wait()
	fmt.Println("Server exited cleanly")
}
