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

// CreateTable generates CREATE TABLE SQL from struct tags
func CreateTable(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr { t = t.Elem() }

	goToSQL := map[reflect.Kind]string{
		reflect.Int: "INTEGER", reflect.String: "TEXT", reflect.Float64: "REAL", reflect.Bool: "BOOLEAN",
	}
	var cols []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" { continue }
		sqlType := goToSQL[f.Type.Kind()]
		extra := ""
		if tag == "id" { extra = " PRIMARY KEY AUTOINCREMENT" }
		cols = append(cols, fmt.Sprintf("  %s %s%s", tag, sqlType, extra))
	}
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);", strings.ToLower(t.Name()), strings.Join(cols, ",\n"))
}

// Insert generates INSERT SQL and values list
func Insert(v any) (string, []any) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	if rv.Kind() == reflect.Ptr { rv = rv.Elem(); rt = rt.Elem() }

	var cols, placeholders []string
	var values []any
	ph := 1
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" || tag == "id" { continue }
		cols = append(cols, tag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", ph)); ph++
		values = append(values, rv.Field(i).Interface())
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		strings.ToLower(rt.Name()), strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	return sql, values
}

// ScanInto simulates scanning rows into structs using reflection
func ScanInto(rows []map[string]any, target any) {
	slicePtr := reflect.ValueOf(target).Elem()
	elemType := slicePtr.Type().Elem()

	for _, row := range rows {
		elem := reflect.New(elemType).Elem()
		for i := 0; i < elemType.NumField(); i++ {
			f := elemType.Field(i)
			tag := f.Tag.Get("db")
			if val, ok := row[tag]; ok {
				fv := elem.Field(i)
				rv := reflect.ValueOf(val)
				if rv.Type().ConvertibleTo(fv.Type()) {
					fv.Set(rv.Convert(fv.Type()))
				}
			}
		}
		slicePtr.Set(reflect.Append(slicePtr, elem))
	}
}

func main() {
	u := User{}

	fmt.Println("=== CREATE TABLE ===")
	fmt.Println(CreateTable(u))

	fmt.Println("\n=== INSERT ===")
	u2 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	sql, vals := Insert(u2)
	fmt.Println("SQL:", sql)
	fmt.Println("Values:", vals)

	fmt.Println("\n=== SCAN INTO ===")
	rows := []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@example.com", "age": 30},
		{"id": 2, "name": "Bob",   "email": "bob@example.com",   "age": 25},
	}
	var users []User
	ScanInto(rows, &users)
	for _, u := range users {
		fmt.Printf("  %+v\n", u)
	}
}
