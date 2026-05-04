//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"unsafe"
)

// BAD ordering — wastes memory due to padding
type BadLayout struct {
	a bool    // 1 byte + 7 padding
	b float64 // 8 bytes
	c bool    // 1 byte + 7 padding
	d int32   // 4 bytes + 4 padding
}

// TODO: GoodLayout — reorder fields to minimize padding
// Hint: largest alignment first (float64 > int32 > bool)
type GoodLayout struct {
	// TODO: reorder a, b, c, d from BadLayout
	b float64
	d int32
	a bool
	c bool
}

func main() {
	bad := BadLayout{}
	good := GoodLayout{}

	fmt.Printf("BadLayout  size=%d bytes
", unsafe.Sizeof(bad))
	fmt.Printf("GoodLayout size=%d bytes
", unsafe.Sizeof(good))
	// Expected: BadLayout ~32, GoodLayout ~16

	fmt.Printf("
BadLayout  offsets: a=%d b=%d c=%d d=%d
",
		unsafe.Offsetof(bad.a), unsafe.Offsetof(bad.b),
		unsafe.Offsetof(bad.c), unsafe.Offsetof(bad.d))

	fmt.Printf("GoodLayout offsets: b=%d d=%d a=%d c=%d
",
		unsafe.Offsetof(good.b), unsafe.Offsetof(good.d),
		unsafe.Offsetof(good.a), unsafe.Offsetof(good.c))

	saved := unsafe.Sizeof(bad) - unsafe.Sizeof(good)
	fmt.Printf("
Saved %d bytes per struct
", saved)
}
