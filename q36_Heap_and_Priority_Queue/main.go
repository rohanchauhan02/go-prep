package main

import (
	"container/heap"
	"fmt"
)

// TODO: implement heap.Interface for a min-heap of ints
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	// TODO: remove and return last element
	return 0
}

// Item for K-way merge
type Item struct{ val, row, col int }
type ItemHeap []Item

func (h ItemHeap) Len() int           { return len(h) }
func (h ItemHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h ItemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ItemHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *ItemHeap) Pop() any {
	old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

// TODO: mergeKSorted merges K sorted arrays using min-heap
// Push first element of each array, pop min, push next from same row
func mergeKSorted(arrays [][]int) []int {
	// TODO: implement
	return nil
}

func main() {
	h := &IntHeap{5, 2, 9, 1, 7}
	heap.Init(h)
	fmt.Print("sorted: ")
	for h.Len() > 0 { fmt.Print(heap.Pop(h), " ") }
	fmt.Println()
	// Expected: 1 2 5 7 9

	arrays := [][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}
	fmt.Println(mergeKSorted(arrays))
	// Expected: [1 2 3 4 5 6 7 8 9]
}
