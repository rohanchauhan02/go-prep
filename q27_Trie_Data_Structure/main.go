package main

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct{ root *TrieNode }

func NewTrie() *Trie { return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}} }

func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok { return false }
		node = node.children[ch]
	}
	return node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for _, ch := range prefix {
		if _, ok := node.children[ch]; !ok { return false }
		node = node.children[ch]
	}
	return true
}

func (t *Trie) Autocomplete(prefix string) []string {
	node := t.root
	for _, ch := range prefix {
		if _, ok := node.children[ch]; !ok { return nil }
		node = node.children[ch]
	}
	var results []string
	var dfs func(n *TrieNode, cur string)
	dfs = func(n *TrieNode, cur string) {
		if n.isEnd { results = append(results, cur) }
		for ch, child := range n.children { dfs(child, cur+string(ch)) }
	}
	dfs(node, prefix)
	return results
}

func main() {
	t := NewTrie()
	words := []string{"apple", "app", "application", "apply", "banana", "band", "bandana"}
	for _, w := range words { t.Insert(w) }

	fmt.Println("Search 'app'        :", t.Search("app"))
	fmt.Println("Search 'ap'         :", t.Search("ap"))
	fmt.Println("StartsWith 'app'    :", t.StartsWith("app"))
	fmt.Println("StartsWith 'xyz'    :", t.StartsWith("xyz"))
	fmt.Println("Autocomplete 'app'  :", t.Autocomplete("app"))
	fmt.Println("Autocomplete 'ban'  :", t.Autocomplete("ban"))
}
