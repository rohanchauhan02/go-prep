package main

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct{ root *TrieNode }

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

// TODO: Insert adds word char-by-char, marks last node isEnd=true
func (t *Trie) Insert(word string) {
	// TODO: implement
}

// TODO: Search returns true only if word exists AND last node isEnd=true
func (t *Trie) Search(word string) bool {
	// TODO: implement
	return false
}

// TODO: StartsWith returns true if any word begins with prefix
func (t *Trie) StartsWith(prefix string) bool {
	// TODO: implement
	return false
}

// TODO: Autocomplete returns all words starting with prefix (DFS from prefix node)
func (t *Trie) Autocomplete(prefix string) []string {
	// TODO: navigate to prefix end, then DFS collecting isEnd words
	return nil
}

func main() {
	t := NewTrie()
	for _, w := range []string{"app", "apple", "apply", "application", "banana"} {
		t.Insert(w)
	}

	fmt.Println(t.Search("app"))         // true
	fmt.Println(t.Search("ap"))          // false
	fmt.Println(t.StartsWith("app"))     // true
	fmt.Println(t.StartsWith("xyz"))     // false
	fmt.Println(t.Autocomplete("app"))   // [app apple apply application] (any order)
}
