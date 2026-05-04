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

func newNode(v int) *TreeNode { return &TreeNode{Val: v} }

func serialize(root *TreeNode) string {
	if root == nil { return "#" }
	return fmt.Sprintf("%d,%s,%s", root.Val, serialize(root.Left), serialize(root.Right))
}

func deserialize(data string) *TreeNode {
	parts := strings.Split(data, ",")
	idx := 0
	var build func() *TreeNode
	build = func() *TreeNode {
		if idx >= len(parts) || parts[idx] == "#" { idx++; return nil }
		v, _ := strconv.Atoi(parts[idx]); idx++
		return &TreeNode{Val: v, Left: build(), Right: build()}
	}
	return build()
}

func levelOrder(root *TreeNode) [][]int {
	if root == nil { return nil }
	var result [][]int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		var level []int
		for i := 0; i < size; i++ {
			node := queue[0]; queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil { queue = append(queue, node.Left) }
			if node.Right != nil { queue = append(queue, node.Right) }
		}
		result = append(result, level)
	}
	return result
}

func main() {
	//       1
	//      / 	//     2   3
	//    / \   	//   4   5   6
	root := newNode(1)
	root.Left = newNode(2); root.Right = newNode(3)
	root.Left.Left = newNode(4); root.Left.Right = newNode(5)
	root.Right.Right = newNode(6)

	serialized := serialize(root)
	fmt.Println("Serialized:", serialized)

	deserialized := deserialize(serialized)
	fmt.Println("Re-serialized:", serialize(deserialized))

	fmt.Println("Level order:", levelOrder(root))

	// Edge cases
	fmt.Println("Empty tree:", serialize(nil))
	fmt.Println("Single node:", serialize(newNode(42)))
}
