//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

// TODO: CreateTable generates "CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, ...)"
// Use reflect to read db tags and map Go types to SQL types:
// int → INTEGER, string → TEXT, float64 → REAL, bool → BOOLEAN
// Field with tag "id" gets "PRIMARY KEY AUTOINCREMENT"
func CreateTable(v any) string {
	// TODO: implement using reflect.TypeOf
	_ = strings.ToLower // hint
	return ""
}

// TODO: Insert generates "INSERT INTO user (name, email, age) VALUES ($1, $2, $3)"
// Skip the "id" field. Return (sql string, values []any)
func Insert(v any) (string, []any) {
	// TODO: implement using reflect.TypeOf and reflect.ValueOf
	return "", nil
}

// TODO: ScanInto maps []map[string]any rows into a slice of structs
// Match map keys to db tags, set field values using reflection
func ScanInto(rows []map[string]any, target any) {
	// TODO: implement using reflect.ValueOf(target).Elem()
}

func main() {
	fmt.Println("=== CREATE TABLE ===")
	fmt.Println(CreateTable(User{}))
	// Expected:
	// CREATE TABLE IF NOT EXISTS user (
	//   id INTEGER PRIMARY KEY AUTOINCREMENT,
	//   name TEXT,
	//   email TEXT,
	//   age INTEGER
	// );

	fmt.Println("
=== INSERT ===")
	sql, vals := Insert(User{Name: "Alice", Email: "alice@x.com", Age: 30})
	fmt.Println(sql)
	fmt.Println(vals)
	// Expected:
	// INSERT INTO user (name, email, age) VALUES ($1, $2, $3)
	// [Alice alice@x.com 30]

	fmt.Println("
=== SCAN ===")
	rows := []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@x.com", "age": 30},
		{"id": 2, "name": "Bob",   "email": "bob@x.com",   "age": 25},
	}
	var users []User
	ScanInto(rows, &users)
	for _, u := range users { fmt.Printf("  %+v
", u) }
}
