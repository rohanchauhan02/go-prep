package main

import (
	"fmt"
	"unsafe"
)

func printSliceInfo(label string, s []int) {
	hdr := (*[3]uintptr)(unsafe.Pointer(&s))
	fmt.Printf("%s: len=%d cap=%d ptr=0x%x data=%v\n", label, len(s), cap(s), hdr[0], s)
}

func sharedBackingArray() {
	orig := []int{1, 2, 3, 4, 5}
	sub := orig[1:3] // shares backing array
	fmt.Println("Before mutation:")
	printSliceInfo("orig", orig)
	printSliceInfo("sub ", sub)

	sub[0] = 99 // MUTATES orig!
	fmt.Println("After sub[0]=99:")
	printSliceInfo("orig", orig) // orig[1] is now 99
	printSliceInfo("sub ", sub)
}

func safeWithCopy() {
	orig := []int{1, 2, 3, 4, 5}
	sub := make([]int, 2)
	copy(sub, orig[1:3]) // independent copy
	sub[0] = 99
	fmt.Println("orig after copy+mutate:", orig) // orig unchanged
}

func appendGrowth() {
	s := make([]int, 0, 3)
	for i := 0; i < 6; i++ {
		before := cap(s)
		s = append(s, i)
		if cap(s) != before {
			fmt.Printf("cap grew: %d → %d at len=%d\n", before, cap(s), len(s))
		}
	}
}

func main() {
	fmt.Println("=== Shared Backing Array Bug ===")
	sharedBackingArray()
	fmt.Println("\n=== Fix with copy() ===")
	safeWithCopy()
	fmt.Println("\n=== Append Capacity Growth ===")
	appendGrowth()
}
