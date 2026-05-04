package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// TODO: UnixTime wraps time.Time but marshals as Unix timestamp (int64)
type UnixTime struct{ time.Time }

// TODO: MarshalJSON — if zero return "null", else marshal Unix seconds
func (u UnixTime) MarshalJSON() ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// TODO: UnmarshalJSON — if "null" set zero, else parse int64 as Unix seconds
func (u *UnixTime) UnmarshalJSON(data []byte) error {
	// TODO: implement
	return nil
}

type Event struct {
	Name      string   `json:"name"`
	CreatedAt UnixTime `json:"created_at"`
}

func main() {
	e := Event{Name: "launch", CreatedAt: UnixTime{time.Now().UTC()}}
	b, _ := json.Marshal(e)
	fmt.Println("marshaled:", string(b))
	// Expected: {"name":"launch","created_at":1234567890}

	var e2 Event
	json.Unmarshal(b, &e2)
	fmt.Println("name:", e2.Name)
	fmt.Println("time:", e2.CreatedAt.Time.UTC().Format(time.RFC3339))
	// Expected: original time

	// null handling
	json.Unmarshal([]byte(`{"name":"x","created_at":null}`), &e2)
	fmt.Println("zero?", e2.CreatedAt.IsZero()) // Expected: true
}
