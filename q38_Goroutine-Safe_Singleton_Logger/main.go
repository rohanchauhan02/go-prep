//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type Level int

const (DEBUG Level = iota; INFO; WARN; ERROR)

func (l Level) String() string { return [...]string{"DEBUG","INFO","WARN","ERROR"}[l] }

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
	level  Level
}

var (
	logInstance *Logger
	logOnce     sync.Once
)

// TODO: GetLogger returns singleton Logger, initialized once
// Write to both stdout and "app.log" file using io.MultiWriter
func GetLogger() *Logger {
	logOnce.Do(func() {
		// TODO: open app.log, create io.MultiWriter, init logInstance
	})
	return logInstance
}

// TODO: log checks level filter, gets caller info via runtime.Caller(2)
// format: [LEVEL] file:line message
func (l *Logger) log(level Level, msg string) {
	if level < l.level { return }
	_, file, line, _ := runtime.Caller(2)
	l.mu.Lock(); defer l.mu.Unlock()
	// TODO: l.logger.Printf(...)
}

func (l *Logger) Debug(msg string) { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, msg) }

func main() {
	_ = io.MultiWriter // hint
	_ = os.OpenFile    // hint
	lg := GetLogger()
	lg.Info("server started")
	lg.Warn("high memory")
	lg.Error("connection failed")
	fmt.Println("logging done")
}
