package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ctxKey string

const (
	UserIDKey    ctxKey = "userID"
	RequestIDKey ctxKey = "requestID"
)

type Handler func(ctx context.Context) error

func AuthMiddleware(next Handler) Handler {
	return func(ctx context.Context) error {
		userID := ctx.Value(UserIDKey)
		if userID == nil {
			return errors.New("unauthorized: no userID in context")
		}
		fmt.Printf("[Auth] user=%v\n", userID)
		return next(ctx)
	}
}

func LogMiddleware(next Handler) Handler {
	return func(ctx context.Context) error {
		reqID := ctx.Value(RequestIDKey)
		fmt.Printf("[Log] requestID=%v start\n", reqID)
		err := next(ctx)
		fmt.Printf("[Log] requestID=%v done err=%v\n", reqID, err)
		return err
	}
}

func BusinessHandler(ctx context.Context) error {
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("[Handler] processed request")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user-42")
	ctx = context.WithValue(ctx, RequestIDKey, "req-abc")
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	chain := LogMiddleware(AuthMiddleware(BusinessHandler))
	if err := chain(ctx); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Missing auth scenario ---")
	noAuth := context.Background()
	if err := chain(noAuth); err != nil {
		fmt.Println("Error:", err)
	}
}
