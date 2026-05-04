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
	age   int    // unexported — skipped
}

func serialize(v any) map[string]any {
	result := make(map[string]any)
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem(); rt = rt.Elem()
	}
	if rv.Kind() != reflect.Struct { return result }

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)
		if !field.IsExported() { continue }

		tag := field.Tag.Get("json")
		name := field.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" { name = parts[0] }
			if len(parts) > 1 && parts[1] == "omitempty" && value.IsZero() { continue }
		}
		result[name] = value.Interface()
	}
	return result
}

func main() {
	u := User{ID: 1, Name: "Alice", Email: ""}
	m := serialize(u)
	fmt.Println("Serialized (omitempty hides empty Email):")
	for k, v := range m { fmt.Printf("  %s: %v\n", k, v) }

	u2 := User{ID: 2, Name: "Bob", Email: "bob@example.com"}
	m2 := serialize(&u2) // also works with pointer
	fmt.Println("\nSerialized (Email present):")
	for k, v := range m2 { fmt.Printf("  %s: %v\n", k, v) }

	// Reflect on types
	fmt.Println("\nField types via reflection:")
	rt := reflect.TypeOf(u)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fmt.Printf("  %s %v tag=%q\n", f.Name, f.Type, f.Tag.Get("json"))
	}
}
