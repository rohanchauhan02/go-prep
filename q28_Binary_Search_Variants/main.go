package main

import "fmt"

// 1. Exact match
func binarySearch(arr []int, target int) int {
	lo, hi := 0, len(arr)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target { return mid } else if arr[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
	}
	return -1
}

// 2. First occurrence
func firstOccurrence(arr []int, target int) int {
	lo, hi, res := 0, len(arr)-1, -1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target { res = mid; hi = mid - 1 } else if arr[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
	}
	return res
}

// 3. Last occurrence
func lastOccurrence(arr []int, target int) int {
	lo, hi, res := 0, len(arr)-1, -1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target { res = mid; lo = mid + 1 } else if arr[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
	}
	return res
}

// 4. Rotated sorted array
func searchRotated(arr []int, target int) int {
	lo, hi := 0, len(arr)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target { return mid }
		if arr[lo] <= arr[mid] { // left half sorted
			if arr[lo] <= target && target < arr[mid] { hi = mid - 1 } else { lo = mid + 1 }
		} else { // right half sorted
			if arr[mid] < target && target <= arr[hi] { lo = mid + 1 } else { hi = mid - 1 }
		}
	}
	return -1
}

func main() {
	sorted := []int{1, 2, 2, 2, 3, 4, 5}
	fmt.Println("=== Binary Search Variants ===")
	fmt.Println("Exact 2         :", binarySearch(sorted, 2))
	fmt.Println("First 2         :", firstOccurrence(sorted, 2))
	fmt.Println("Last 2          :", lastOccurrence(sorted, 2))
	fmt.Println("Exact 6 (miss)  :", binarySearch(sorted, 6))

	rotated := []int{4, 5, 6, 7, 0, 1, 2}
	fmt.Println("\nRotated array:", rotated)
	for _, t := range []int{0, 3, 7} {
		fmt.Printf("  search %d → index %d\n", t, searchRotated(rotated, t))
	}
}
