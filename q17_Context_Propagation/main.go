package main

import (
	"context"
	"fmt"
	"time"
)

type ctxKey string

const (
	UserIDKey    ctxKey = "userID"
	RequestIDKey ctxKey = "requestID"
)

type Handler func(ctx context.Context) error

// TODO: AuthMiddleware reads UserIDKey from ctx
// if missing → return error "unauthorized"
func AuthMiddleware(next Handler) Handler {
	return func(ctx context.Context) error {
		// TODO: check ctx.Value(UserIDKey) != nil
		return next(ctx)
	}
}

// TODO: LogMiddleware prints [Log] requestID=<id> before and after next
func LogMiddleware(next Handler) Handler {
	return func(ctx context.Context) error {
		// TODO: log start, call next, log done
		return next(ctx)
	}
}

// TODO: BusinessHandler sleeps 100ms respecting ctx cancellation
func BusinessHandler(ctx context.Context) error {
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("[Handler] done")
		return nil
	case <-ctx.Done():
		// TODO: return ctx.Err()
		return nil
	}
}

func main() {
	ctx := context.WithValue(context.Background(), UserIDKey, "user-42")
	ctx = context.WithValue(ctx, RequestIDKey, "req-abc")
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	chain := LogMiddleware(AuthMiddleware(BusinessHandler))
	if err := chain(ctx); err != nil {
		fmt.Println("Error:", err)
	}
	// Expected: [Log] start, [Auth] OK, [Handler] done, [Log] done

	fmt.Println("--- no auth ---")
	if err := chain(context.Background()); err != nil {
		fmt.Println("Error:", err) // Expected: Error: unauthorized
	}
}
