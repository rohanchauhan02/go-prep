package main

import (
	"fmt"
	"unsafe"
)

// Bad ordering — lots of padding
type BadLayout struct {
	a bool    // 1 byte  + 7 padding
	b float64 // 8 bytes
	c bool    // 1 byte  + 7 padding
	d int32   // 4 bytes + 4 padding
} // total: ~32 bytes

// Good ordering — minimal padding
type GoodLayout struct {
	b float64 // 8 bytes
	d int32   // 4 bytes
	a bool    // 1 byte
	c bool    // 1 byte
	// 2 bytes padding to align to 8
} // total: ~16 bytes

func printLayout(name string, size, ab, bb, cb, db uintptr) {
	fmt.Printf("%-12s size=%-3d a@%-3d b@%-3d c@%-3d d@%-3d\n", name, size, ab, bb, cb, db)
}

func main() {
	bad := BadLayout{}
	good := GoodLayout{}

	fmt.Println("=== Struct Memory Layout ===")
	printLayout("BadLayout",
		unsafe.Sizeof(bad),
		unsafe.Offsetof(bad.a),
		unsafe.Offsetof(bad.b),
		unsafe.Offsetof(bad.c),
		unsafe.Offsetof(bad.d),
	)
	printLayout("GoodLayout",
		unsafe.Sizeof(good),
		unsafe.Offsetof(good.a),
		unsafe.Offsetof(good.b),
		unsafe.Offsetof(good.c),
		unsafe.Offsetof(good.d),
	)

	savings := unsafe.Sizeof(bad) - unsafe.Sizeof(good)
	fmt.Printf("\nMemory saved per struct: %d bytes (%.0f%%)\n",
		savings, float64(savings)/float64(unsafe.Sizeof(bad))*100)

	// Alignment of types
	fmt.Println("\nType alignments:")
	fmt.Println("  bool   :", unsafe.Alignof(bool(false)))
	fmt.Println("  int32  :", unsafe.Alignof(int32(0)))
	fmt.Println("  float64:", unsafe.Alignof(float64(0)))
	fmt.Println("  string :", unsafe.Alignof(string("")))
}
