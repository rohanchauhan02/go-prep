package main

import "fmt"

type Color struct{ R, G, B, A uint8 }

// fmt.Stringer — used by %v, %s
func (c Color) String() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%d)", c.R, c.G, c.B, c.A)
}

// fmt.GoStringer — used by %#v
func (c Color) GoString() string {
	return fmt.Sprintf("Color{R:%d, G:%d, B:%d, A:%d}", c.R, c.G, c.B, c.A)
}

func main() {
	red := Color{255, 0, 0, 255}

	fmt.Printf("%%v  → %v\n", red)   // calls String()
	fmt.Printf("%%s  → %s\n", red)   // calls String()
	fmt.Printf("%%+v → %+v\n", red)  // struct fields (ignores Stringer)
	fmt.Printf("%%#v → %#v\n", red)  // calls GoString()
	fmt.Printf("%%T  → %T\n", red)   // type only

	// In a collection
	palette := []Color{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}}
	fmt.Println("\nPalette:", palette)
}
