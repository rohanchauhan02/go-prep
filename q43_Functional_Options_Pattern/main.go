package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Server struct {
	addr     string
	timeout  time.Duration
	maxConns int
	tls      *tls.Config
}

// Option is a function that configures a Server
type Option func(*Server)

// TODO: WithAddr sets the server address
func WithAddr(addr string) Option {
	return func(s *Server) { s.addr = addr }
}

// TODO: WithTimeout sets the server timeout
func WithTimeout(d time.Duration) Option {
	// TODO: implement
	return func(s *Server) {}
}

// TODO: WithMaxConns sets the max connections
func WithMaxConns(n int) Option {
	// TODO: implement
	return func(s *Server) {}
}

// TODO: WithTLS sets TLS config
func WithTLS(cfg *tls.Config) Option {
	// TODO: implement
	return func(s *Server) {}
}

// TODO: NewServer creates Server with sane defaults then applies opts
// defaults: addr=":8080", timeout=30s, maxConns=100
func NewServer(opts ...Option) *Server {
	s := &Server{
		// TODO: set defaults
	}
	for _, o := range opts { o(s) }
	return s
}

func main() {
	s1 := NewServer()
	fmt.Println(s1.addr, s1.timeout, s1.maxConns)
	// Expected: :8080 30s 100

	s2 := NewServer(WithAddr(":9090"), WithTimeout(10*time.Second), WithMaxConns(500))
	fmt.Println(s2.addr, s2.timeout, s2.maxConns)
	// Expected: :9090 10s 500

	s3 := NewServer(WithTLS(&tls.Config{}))
	fmt.Println(s3.tls != nil)
	// Expected: true
}
