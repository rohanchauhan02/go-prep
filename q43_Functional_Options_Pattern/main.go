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

type Option func(*Server)

func WithAddr(addr string) Option       { return func(s *Server) { s.addr = addr } }
func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }
func WithMaxConns(n int) Option          { return func(s *Server) { s.maxConns = n } }
func WithTLS(cfg *tls.Config) Option     { return func(s *Server) { s.tls = cfg } }

func NewServer(opts ...Option) *Server {
	// sane defaults
	s := &Server{
		addr:     ":8080",
		timeout:  30 * time.Second,
		maxConns: 100,
	}
	for _, o := range opts { o(s) }
	return s
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{addr=%s timeout=%v maxConns=%d tls=%v}",
		s.addr, s.timeout, s.maxConns, s.tls != nil)
}

func main() {
	// Default server
	s1 := NewServer()
	fmt.Println("Default:", s1)

	// Custom server
	s2 := NewServer(
		WithAddr(":9090"),
		WithTimeout(10*time.Second),
		WithMaxConns(500),
		WithTLS(&tls.Config{}),
	)
	fmt.Println("Custom:", s2)

	// Partial override
	s3 := NewServer(WithMaxConns(50))
	fmt.Println("Partial:", s3)
}
