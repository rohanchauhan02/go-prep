package main

import "fmt"

// 1. Max sum subarray of size k
func maxSumK(arr []int, k int) int {
	sum, maxSum := 0, 0
	for i := 0; i < k; i++ { sum += arr[i] }
	maxSum = sum
	for i := k; i < len(arr); i++ {
		sum += arr[i] - arr[i-k]
		if sum > maxSum { maxSum = sum }
	}
	return maxSum
}

// 2. Longest substring without repeating characters
func longestUnique(s string) int {
	freq := make(map[byte]int)
	left, maxLen := 0, 0
	for right := 0; right < len(s); right++ {
		freq[s[right]]++
		for freq[s[right]] > 1 { freq[s[left]]--; left++ }
		if right-left+1 > maxLen { maxLen = right - left + 1 }
	}
	return maxLen
}

// 3. Minimum window substring
func minWindow(s, t string) string {
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ { need[t[i]]++ }
	have, total := 0, len(need)
	window := make(map[byte]int)
	left, minLen, start := 0, len(s)+1, 0
	for right := 0; right < len(s); right++ {
		c := s[right]; window[c]++
		if need[c] > 0 && window[c] == need[c] { have++ }
		for have == total {
			if right-left+1 < minLen { minLen = right - left + 1; start = left }
			lc := s[left]; window[lc]--
			if need[lc] > 0 && window[lc] < need[lc] { have-- }
			left++
		}
	}
	if minLen > len(s) { return "" }
	return s[start : start+minLen]
}

func main() {
	fmt.Println("=== Max Sum Subarray k=3 ===")
	fmt.Println(maxSumK([]int{2, 1, 5, 1, 3, 2}, 3)) // 9

	fmt.Println("\n=== Longest Unique Substring ===")
	for _, s := range []string{"abcabcbb", "bbbbb", "pwwkew"} {
		fmt.Printf("  %-12s → %d\n", s, longestUnique(s))
	}

	fmt.Println("\n=== Minimum Window Substring ===")
	fmt.Println(minWindow("ADOBECODEBANC", "ABC")) // BANC
	fmt.Println(minWindow("a", "a"))               // a
	fmt.Println(minWindow("a", "aa"))              // ""
}
