//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type Middleware func(http.Handler) http.Handler

// TODO: Chain applies middlewares right-to-left so they execute left-to-right
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	// TODO: iterate mws in reverse, wrap h
	return h
}

// TODO: Logger middleware — print method, path, and duration after next
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[Log] %s %s %v
", r.Method, r.URL.Path, time.Since(start))
	})
}

// TODO: Auth middleware — check "Authorization: Bearer secret"
// if missing/wrong → 401 Unauthorized, don't call next
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement
		next.ServeHTTP(w, r)
	})
}

func businessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func main() {
	handler := Chain(http.HandlerFunc(businessHandler), Logger, Auth)

	// Valid token
	req := httptest.NewRequest("GET", "/api", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	fmt.Println("status:", rr.Code) // Expected: 200

	// No token
	req2 := httptest.NewRequest("GET", "/api", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	fmt.Println("status:", rr2.Code) // Expected: 401
}
