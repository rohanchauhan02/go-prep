//go:build ignore
// Remove the above line when implementing

package main

import (
	"context"
	"fmt"
	"time"
)

type Request struct{ Method string }
type Response struct{ Data string }
type Handler func(ctx context.Context, req Request) (Response, error)
type Interceptor func(ctx context.Context, req Request, next Handler) (Response, error)

// TODO: ChainInterceptors wraps handler with interceptors (last interceptor outermost)
func ChainInterceptors(h Handler, interceptors ...Interceptor) Handler {
	// TODO: reverse-iterate interceptors, wrap h each time
	return h
}

// TODO: LoggingInterceptor prints method, calls next, prints duration
func LoggingInterceptor(ctx context.Context, req Request, next Handler) (Response, error) {
	fmt.Printf("[Log] %s
", req.Method)
	start := time.Now()
	resp, err := next(ctx, req)
	fmt.Printf("[Log] done %v
", time.Since(start))
	return resp, err
}

// TODO: AuthInterceptor checks ctx value "token" == "valid", else return error
func AuthInterceptor(ctx context.Context, req Request, next Handler) (Response, error) {
	// TODO: implement
	return next(ctx, req)
}

func businessHandler(ctx context.Context, req Request) (Response, error) {
	return Response{Data: "result:" + req.Method}, nil
}

func main() {
	handler := ChainInterceptors(businessHandler, LoggingInterceptor, AuthInterceptor)

	ctx := context.WithValue(context.Background(), "token", "valid")
	resp, err := handler(ctx, Request{Method: "GetUser"})
	fmt.Println(resp.Data, err) // Expected: result:GetUser <nil>

	resp2, err2 := handler(context.Background(), Request{Method: "GetUser"})
	fmt.Println(resp2, err2) // Expected: {} auth: invalid token
}
