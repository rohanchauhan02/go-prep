package main

import "fmt"

// TODO: exact binary search — return index or -1
func binarySearch(arr []int, target int) int {
	// TODO: lo, hi, mid — standard binary search
	return -1
}

// TODO: first occurrence of target — return leftmost index or -1
func firstOccurrence(arr []int, target int) int {
	// TODO: when found, save result and continue left (hi = mid-1)
	return -1
}

// TODO: last occurrence of target — return rightmost index or -1
func lastOccurrence(arr []int, target int) int {
	// TODO: when found, save result and continue right (lo = mid+1)
	return -1
}

// TODO: search in rotated sorted array — O(log n)
// Hint: one half is always sorted — check which and narrow accordingly
func searchRotated(arr []int, target int) int {
	// TODO: implement
	return -1
}

func main() {
	arr := []int{1, 2, 2, 2, 3, 4, 5}
	fmt.Println(binarySearch(arr, 2))       // any valid index (1,2,3)
	fmt.Println(firstOccurrence(arr, 2))    // Expected: 1
	fmt.Println(lastOccurrence(arr, 2))     // Expected: 3
	fmt.Println(binarySearch(arr, 6))       // Expected: -1

	rotated := []int{4, 5, 6, 7, 0, 1, 2}
	fmt.Println(searchRotated(rotated, 0))  // Expected: 4
	fmt.Println(searchRotated(rotated, 3))  // Expected: -1
	fmt.Println(searchRotated(rotated, 7))  // Expected: 3
}
