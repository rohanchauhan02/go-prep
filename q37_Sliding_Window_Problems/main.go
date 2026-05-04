package main

import "fmt"

// TODO: max sum of subarray of exactly size k — O(n)
func maxSumK(arr []int, k int) int {
	// TODO: sliding window — add right, remove left
	return 0
}

// TODO: longest substring without repeating characters — O(n)
func longestUnique(s string) int {
	// TODO: freq map + two pointers
	return 0
}

// TODO: minimum window in s containing all chars of t — O(n)
// return "" if impossible
func minWindow(s, t string) string {
	// TODO: need map, have counter, two pointers
	return ""
}

func main() {
	fmt.Println(maxSumK([]int{2, 1, 5, 1, 3, 2}, 3))  // Expected: 9
	fmt.Println(maxSumK([]int{1, 2, 3, 4, 5}, 2))      // Expected: 9

	fmt.Println(longestUnique("abcabcbb")) // Expected: 3
	fmt.Println(longestUnique("bbbbb"))    // Expected: 1
	fmt.Println(longestUnique("pwwkew"))   // Expected: 3

	fmt.Println(minWindow("ADOBECODEBANC", "ABC")) // Expected: BANC
	fmt.Println(minWindow("a", "a"))               // Expected: a
	fmt.Println(minWindow("a", "aa"))              // Expected: ""
}
