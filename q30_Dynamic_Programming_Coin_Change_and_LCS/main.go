package main

import "fmt"

// TODO: coinChange (bottom-up DP)
// Return minimum number of coins to make amount, or -1 if impossible
// coins=[1,3,4] amount=6 → 2 (3+3)
func coinChange(coins []int, amount int) int {
	// TODO: dp[i] = min coins to make i
	// dp[0]=0, dp[i] = min(dp[i-c]+1) for each coin c <= i
	return -1
}

// TODO: lcs (2D bottom-up DP)
// Return length of longest common subsequence of strings a and b
// lcs("ABCBDAB","BDCAB") → 4
func lcs(a, b string) int {
	// TODO: dp[i][j] = LCS of a[:i] and b[:j]
	return 0
}

func main() {
	fmt.Println(coinChange([]int{1, 3, 4}, 6))  // Expected: 2
	fmt.Println(coinChange([]int{2}, 3))         // Expected: -1
	fmt.Println(coinChange([]int{1}, 0))         // Expected: 0

	fmt.Println(lcs("ABCBDAB", "BDCAB"))         // Expected: 4
	fmt.Println(lcs("AGGTAB", "GXTXAYB"))        // Expected: 4
	fmt.Println(lcs("", "ABC"))                  // Expected: 0
}
