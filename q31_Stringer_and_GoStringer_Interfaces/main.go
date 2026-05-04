//go:build ignore
// Remove the above line when implementing

package main

import "fmt"

type Color struct{ R, G, B, A uint8 }

// TODO: implement fmt.Stringer — used by %v and %s
// format: "rgba(R,G,B,A)"
func (c Color) String() string {
	// TODO: implement
	return ""
}

// TODO: implement fmt.GoStringer — used by %#v
// format: "Color{R:R, G:G, B:B, A:A}"
func (c Color) GoString() string {
	// TODO: implement
	return ""
}

func main() {
	red := Color{255, 0, 0, 255}
	fmt.Printf("%%v  → %v
", red)   // Expected: rgba(255,0,0,255)
	fmt.Printf("%%s  → %s
", red)   // Expected: rgba(255,0,0,255)
	fmt.Printf("%%+v → %+v
", red)  // Expected: {R:255 G:0 B:0 A:255}  (struct fields)
	fmt.Printf("%%#v → %#v
", red)  // Expected: Color{R:255, G:0, B:0, A:255}
}
