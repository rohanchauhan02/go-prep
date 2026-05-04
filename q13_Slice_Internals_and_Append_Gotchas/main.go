//go:build ignore
// Remove the above line when implementing

package main

import "fmt"

// TODO: show shared backing array bug
// orig := []int{1,2,3,4,5}
// sub := orig[1:3]
// mutating sub[0] changes orig[1] — demonstrate this
func sharedArrayBug() {
	orig := []int{1, 2, 3, 4, 5}
	sub := orig[1:3]
	fmt.Println("before:", orig) // [1 2 3 4 5]
	// TODO: mutate sub[0] = 99
	fmt.Println("after :", orig) // Expected: [1 99 3 4 5]
}

// TODO: fix with copy — sub should be independent
func fixWithCopy() {
	orig := []int{1, 2, 3, 4, 5}
	sub := make([]int, 2)
	// TODO: copy orig[1:3] into sub
	// TODO: mutate sub[0] = 99
	fmt.Println("orig unchanged:", orig) // Expected: [1 2 3 4 5]
}

// TODO: show how append grows capacity — print cap before and after growth
func appendGrowth() {
	s := make([]int, 0, 3)
	for i := 0; i < 6; i++ {
		before := cap(s)
		s = append(s, i)
		if cap(s) != before {
			fmt.Printf("cap grew: %d → %d at len=%d
", before, cap(s), len(s))
		}
	}
}

func main() {
	sharedArrayBug()
	fixWithCopy()
	appendGrowth()
	// Expected growth: 3 → 6 at len=4
}
