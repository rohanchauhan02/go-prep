package main

import (
	"container/heap"
	"fmt"
)

// MinHeap for integers
type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

// ─── K-way merge ────────────────────────────────
type Item struct{ val, row, col int }
type ItemHeap []Item
func (h ItemHeap) Len() int           { return len(h) }
func (h ItemHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h ItemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ItemHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *ItemHeap) Pop() any {
	old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

func mergeKSorted(arrays [][]int) []int {
	h := &ItemHeap{}
	heap.Init(h)
	for i, arr := range arrays {
		if len(arr) > 0 { heap.Push(h, Item{arr[0], i, 0}) }
	}
	var result []int
	for h.Len() > 0 {
		item := heap.Pop(h).(Item)
		result = append(result, item.val)
		if item.col+1 < len(arrays[item.row]) {
			heap.Push(h, Item{arrays[item.row][item.col+1], item.row, item.col + 1})
		}
	}
	return result
}

func main() {
	h := &IntHeap{5, 2, 9, 1, 7}
	heap.Init(h)
	fmt.Print("Min-heap pops: ")
	for h.Len() > 0 { fmt.Print(heap.Pop(h), " ") }
	fmt.Println()

	arrays := [][]int{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
	}
	fmt.Println("K-way merge:", mergeKSorted(arrays))
}
