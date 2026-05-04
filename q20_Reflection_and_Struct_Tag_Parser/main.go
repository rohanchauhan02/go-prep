package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	age   int    // unexported — must be skipped
}

// TODO: serialize reads struct tags (json:"name") and builds map[string]any
// Rules:
//   - skip unexported fields
//   - use tag name if present, else field name
//   - if tag has ",omitempty" and field is zero value → skip
func serialize(v any) map[string]any {
	result := make(map[string]any)
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem(); rt = rt.Elem()
	}
	// TODO: iterate fields, apply rules above
	_ = strings.Split // hint
	return result
}

func main() {
	u := User{ID: 1, Name: "Alice", Email: ""}
	m := serialize(u)
	fmt.Println("serialized:", m)
	// Expected: map[id:1 name:Alice]  (email omitted — empty + omitempty)

	u2 := User{ID: 2, Name: "Bob", Email: "bob@example.com"}
	m2 := serialize(u2)
	fmt.Println("serialized:", m2)
	// Expected: map[email:bob@example.com id:2 name:Bob]
}
