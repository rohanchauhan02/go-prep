package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }
	return h
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[Logger] %s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer secret" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("[Auth] token valid")
		next.ServeHTTP(w, r)
	})
}

func RateLimit(next http.Handler) http.Handler {
	// simplified: allow all in this demo
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[RateLimit] allowed")
		next.ServeHTTP(w, r)
	})
}

func businessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from business handler")
}

func main() {
	handler := Chain(
		http.HandlerFunc(businessHandler),
		Logger,
		RateLimit,
		Auth,
	)

	// Test with valid token
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	fmt.Println("Status:", rr.Code, "Body:", rr.Body.String())

	// Test without token
	req2 := httptest.NewRequest("GET", "/api/data", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	fmt.Println("Status:", rr2.Code)
}
