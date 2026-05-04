package main

import (
	"fmt"
	"sync"
)

type DB struct{ name string }

var (
	instance *DB
	once     sync.Once
)

// TODO: GetDB must return the same *DB instance always
// The DB must be initialized exactly once using sync.Once
func GetDB() *DB {
	// TODO: use once.Do to initialize instance
	return instance
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := GetDB()
			fmt.Println(db.name)
		}()
	}
	wg.Wait()
	// Expected: "prod-db" printed 5 times, initialized only once
}
