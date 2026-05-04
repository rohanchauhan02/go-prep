package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// TODO: serialize tree to string using preorder traversal
// null nodes → "#", separator → ","
// e.g. tree [1,2,3] → "1,2,#,#,3,#,#"
func serialize(root *TreeNode) string {
	// TODO: implement
	return ""
}

// TODO: deserialize string back to tree
func deserialize(data string) *TreeNode {
	parts := strings.Split(data, ",")
	idx := 0
	var build func() *TreeNode
	build = func() *TreeNode {
		if idx >= len(parts) || parts[idx] == "#" { idx++; return nil }
		v, _ := strconv.Atoi(parts[idx]); idx++
		// TODO: build left and right children
		return &TreeNode{Val: v}
	}
	return build()
}

// TODO: levelOrder returns nodes level by level as [][]int (BFS)
func levelOrder(root *TreeNode) [][]int {
	// TODO: implement using queue
	return nil
}

func main() {
	root := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 6}},
	}

	s := serialize(root)
	fmt.Println("serialized:", s)

	r2 := deserialize(s)
	fmt.Println("re-serialized:", serialize(r2))
	// Expected: same string as original

	fmt.Println("level order:", levelOrder(root))
	// Expected: [[1] [2 3] [4 5 6]]
}
