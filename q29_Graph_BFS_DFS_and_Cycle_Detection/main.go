//go:build ignore
// Remove the above line when implementing

package main

import "fmt"

type Graph struct{ adj map[int][]int }

func NewGraph() *Graph { return &Graph{adj: make(map[int][]int)} }

func (g *Graph) AddEdge(u, v int) { g.adj[u] = append(g.adj[u], v) }

// TODO: BFS from start — return nodes in visited order
func (g *Graph) BFS(start int) []int {
	// TODO: use queue (slice), visited map
	return nil
}

// TODO: DFS iterative from start — return nodes in visited order
func (g *Graph) DFS(start int) []int {
	// TODO: use stack (slice), visited map
	return nil
}

// TODO: HasCycle — detect cycle in directed graph using DFS + recursion stack
func (g *Graph) HasCycle() bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool)
	var dfs func(node int) bool
	dfs = func(node int) bool {
		// TODO: implement
		return false
	}
	for node := range g.adj {
		if !visited[node] && dfs(node) { return true }
	}
	return false
}

// TODO: TopologicalSort using Kahn's algorithm (BFS with in-degree)
func (g *Graph) TopologicalSort() []int {
	// TODO: compute in-degrees, process nodes with 0 in-degree
	return nil
}

func main() {
	g := NewGraph()
	for _, e := range [][2]int{{1,2},{1,3},{2,4},{3,4},{4,5}} {
		g.AddEdge(e[0], e[1])
	}
	fmt.Println("BFS:", g.BFS(1))           // Expected: [1 2 3 4 5]
	fmt.Println("DFS:", g.DFS(1))           // Expected: [1 3 4 5 2] (stack order)
	fmt.Println("Cycle:", g.HasCycle())     // Expected: false
	fmt.Println("Topo:", g.TopologicalSort()) // Expected: [1 2 3 4 5] or valid topo order

	cyclic := NewGraph()
	cyclic.AddEdge(1, 2); cyclic.AddEdge(2, 3); cyclic.AddEdge(3, 1)
	fmt.Println("Cyclic:", cyclic.HasCycle()) // Expected: true
}
