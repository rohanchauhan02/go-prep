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
	logOnce    sync.Once
)

func GetLogger() *Logger {
	logOnce.Do(func() {
		f, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		var out io.Writer
		if err != nil { out = os.Stdout } else { out = io.MultiWriter(os.Stdout, f) }
		logInstance = &Logger{
			logger: log.New(out, "", log.Ldate|log.Ltime),
			level: DEBUG,
		}
	})
	return logInstance
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level { return }
	_, file, line, _ := runtime.Caller(2)
	l.mu.Lock(); defer l.mu.Unlock()
	l.logger.Printf("[%s] %s:%d %s", level, file, line, msg)
}

func (l *Logger) Debug(msg string) { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, msg) }

func main() {
	log := GetLogger()
	log.Debug("starting up")
	log.Info("server ready")
	log.Warn("high memory usage")
	log.Error("connection failed")

	// Same instance from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			GetLogger().Info(fmt.Sprintf("goroutine %d logging", i))
		}(i)
	}
	wg.Wait()
}
