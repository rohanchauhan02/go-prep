package main

import "fmt"

// ─── Coin Change ────────────────────────────────
func coinChangeDP(coins []int, amount int) int {
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ { dp[i] = amount + 1 }
	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] { dp[i] = dp[i-c] + 1 }
		}
	}
	if dp[amount] > amount { return -1 }
	return dp[amount]
}

// Top-down memoization
func coinChangeMemo(coins []int, amount int) int {
	memo := make(map[int]int)
	var dp func(n int) int
	dp = func(n int) int {
		if n == 0 { return 0 }
		if n < 0 { return -1 }
		if v, ok := memo[n]; ok { return v }
		best := -1
		for _, c := range coins {
			res := dp(n - c)
			if res >= 0 && (best < 0 || res+1 < best) { best = res + 1 }
		}
		memo[n] = best; return best
	}
	return dp(amount)
}

// ─── LCS ────────────────────────────────────────
func lcs(a, b string) int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp { dp[i] = make([]int, n+1) }
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] { dp[i][j] = dp[i-1][j-1] + 1 } else {
				if dp[i-1][j] > dp[i][j-1] { dp[i][j] = dp[i-1][j] } else { dp[i][j] = dp[i][j-1] }
			}
		}
	}
	return dp[m][n]
}

func main() {
	fmt.Println("=== Coin Change ===")
	fmt.Println("coins=[1,3,4] amount=6 →", coinChangeDP([]int{1, 3, 4}, 6))
	fmt.Println("coins=[2] amount=3 →", coinChangeDP([]int{2}, 3))
	fmt.Println("memo coins=[1,3,4] amount=6 →", coinChangeMemo([]int{1, 3, 4}, 6))

	fmt.Println("\n=== LCS ===")
	fmt.Println("lcs(ABCBDAB, BDCAB) =", lcs("ABCBDAB", "BDCAB")) // 4
	fmt.Println("lcs(AGGTAB, GXTXAYB) =", lcs("AGGTAB", "GXTXAYB")) // 4
}
