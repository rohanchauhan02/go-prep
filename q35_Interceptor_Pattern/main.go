package main

import (
	"context"
	"fmt"
	"time"
)

type Request struct{ Method string; Payload any }
type Response struct{ Data any }

type Handler func(ctx context.Context, req Request) (Response, error)
type Interceptor func(ctx context.Context, req Request, next Handler) (Response, error)

func ChainInterceptors(h Handler, interceptors ...Interceptor) Handler {
	for i := len(interceptors) - 1; i >= 0; i-- {
		ic := interceptors[i]
		next := h
		h = func(ctx context.Context, req Request) (Response, error) {
			return ic(ctx, req, next)
		}
	}
	return h
}

func LoggingInterceptor(ctx context.Context, req Request, next Handler) (Response, error) {
	fmt.Printf("[Log] method=%s payload=%v\n", req.Method, req.Payload)
	start := time.Now()
	resp, err := next(ctx, req)
	fmt.Printf("[Log] done in %v err=%v\n", time.Since(start), err)
	return resp, err
}

func AuthInterceptor(ctx context.Context, req Request, next Handler) (Response, error) {
	token, _ := ctx.Value("token").(string)
	if token != "valid" {
		return Response{}, fmt.Errorf("auth: invalid token")
	}
	fmt.Println("[Auth] OK")
	return next(ctx, req)
}

func MetricsInterceptor(ctx context.Context, req Request, next Handler) (Response, error) {
	resp, err := next(ctx, req)
	fmt.Printf("[Metrics] method=%s ok=%v\n", req.Method, err == nil)
	return resp, err
}

func businessHandler(ctx context.Context, req Request) (Response, error) {
	return Response{Data: fmt.Sprintf("result for %s", req.Method)}, nil
}

func main() {
	handler := ChainInterceptors(businessHandler, LoggingInterceptor, AuthInterceptor, MetricsInterceptor)

	ctx := context.WithValue(context.Background(), "token", "valid")
	resp, err := handler(ctx, Request{Method: "GetUser", Payload: 42})
	fmt.Println("Response:", resp.Data, "Error:", err)

	fmt.Println("\n--- Without valid token ---")
	resp2, err2 := handler(context.Background(), Request{Method: "GetUser"})
	fmt.Println("Response:", resp2, "Error:", err2)
}
