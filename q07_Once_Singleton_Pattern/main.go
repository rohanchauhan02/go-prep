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

func GetDB() *DB {
	once.Do(func() {
		fmt.Println("  [initializing DB — runs only once]")
		instance = &DB{name: "prod-db"}
	})
	return instance
}

// BAD version (race condition without sync.Once):
var badInstance *DB
var initialized bool

func GetDBBad() *DB {
	if !initialized { // RACE: multiple goroutines can pass this check
		fmt.Println("  [BAD: initializing]")
		badInstance = &DB{name: "prod-db"}
		initialized = true
	}
	return badInstance
}

func main() {
	fmt.Println("=== sync.Once Singleton ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := GetDB()
			fmt.Println("  got:", db.name)
		}()
	}
	wg.Wait()
	fmt.Println("Same instance every time (init ran once)")
}
