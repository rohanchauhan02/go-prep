package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTime struct{ time.Time }

func (u UnixTime) MarshalJSON() ([]byte, error) {
	if u.IsZero() { return []byte("null"), nil }
	return json.Marshal(u.Unix())
}

func (u *UnixTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" { u.Time = time.Time{}; return nil }
	var ts int64
	if err := json.Unmarshal(data, &ts); err != nil { return err }
	u.Time = time.Unix(ts, 0).UTC()
	return nil
}

type Event struct {
	Name      string   `json:"name"`
	CreatedAt UnixTime `json:"created_at"`
	ExpiresAt UnixTime `json:"expires_at,omitempty"`
}

func main() {
	now := UnixTime{time.Now().UTC()}
	e := Event{Name: "launch", CreatedAt: now}

	b, _ := json.Marshal(e)
	fmt.Println("Marshaled:", string(b))

	// Round-trip
	var e2 Event
	json.Unmarshal(b, &e2)
	fmt.Println("Unmarshaled name:", e2.Name)
	fmt.Println("Unmarshaled time:", e2.CreatedAt.Time)

	// null handling
	nullJSON := `{"name":"empty","created_at":null}`
	var e3 Event
	json.Unmarshal([]byte(nullJSON), &e3)
	fmt.Println("Zero time?", e3.CreatedAt.IsZero())
}
