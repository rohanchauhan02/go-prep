package main

import "fmt"

type Graph struct{ adj map[int][]int }

func NewGraph() *Graph { return &Graph{adj: make(map[int][]int)} }
func (g *Graph) AddEdge(u, v int) { g.adj[u] = append(g.adj[u], v) }

func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true
	var order []int
	for len(queue) > 0 {
		node := queue[0]; queue = queue[1:]
		order = append(order, node)
		for _, nb := range g.adj[node] {
			if !visited[nb] { visited[nb] = true; queue = append(queue, nb) }
		}
	}
	return order
}

func (g *Graph) DFSIterative(start int) []int {
	visited := make(map[int]bool)
	stack := []int{start}
	var order []int
	for len(stack) > 0 {
		node := stack[len(stack)-1]; stack = stack[:len(stack)-1]
		if visited[node] { continue }
		visited[node] = true; order = append(order, node)
		for _, nb := range g.adj[node] { if !visited[nb] { stack = append(stack, nb) } }
	}
	return order
}

func (g *Graph) hasCycleUtil(node int, visited, recStack map[int]bool) bool {
	visited[node] = true; recStack[node] = true
	for _, nb := range g.adj[node] {
		if !visited[nb] && g.hasCycleUtil(nb, visited, recStack) { return true }
		if recStack[nb] { return true }
	}
	recStack[node] = false; return false
}

func (g *Graph) HasCycle() bool {
	visited, recStack := make(map[int]bool), make(map[int]bool)
	for node := range g.adj {
		if !visited[node] && g.hasCycleUtil(node, visited, recStack) { return true }
	}
	return false
}

func (g *Graph) TopologicalSort() []int {
	inDegree := make(map[int]int)
	for u := range g.adj { for _, v := range g.adj[u] { inDegree[v]++ } }
	queue := []int{}
	for node := range g.adj { if inDegree[node] == 0 { queue = append(queue, node) } }
	var order []int
	for len(queue) > 0 {
		node := queue[0]; queue = queue[1:]; order = append(order, node)
		for _, nb := range g.adj[node] {
			inDegree[nb]--
			if inDegree[nb] == 0 { queue = append(queue, nb) }
		}
	}
	return order
}

func main() {
	g := NewGraph()
	for _, e := range [][2]int{{1,2},{1,3},{2,4},{3,4},{4,5}} { g.AddEdge(e[0], e[1]) }
	fmt.Println("BFS from 1:", g.BFS(1))
	fmt.Println("DFS from 1:", g.DFSIterative(1))
	fmt.Println("Has cycle: ", g.HasCycle())
	fmt.Println("Topo sort :", g.TopologicalSort())

	cyclic := NewGraph()
	cyclic.AddEdge(1,2); cyclic.AddEdge(2,3); cyclic.AddEdge(3,1)
	fmt.Println("\nCyclic graph has cycle:", cyclic.HasCycle())
}
