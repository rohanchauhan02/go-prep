package main

import (
	"fmt"
	"unicode/utf8"
)

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	s := "Hello, 世界"
	fmt.Println("String     :", s)
	fmt.Println("Bytes      :", len(s))
	fmt.Println("Runes      :", utf8.RuneCountInString(s))

	fmt.Println("\n--- Byte indexing (raw bytes) ---")
	for i := 0; i < len(s); i++ {
		fmt.Printf("s[%d] = 0x%x\n", i, s[i])
	}

	fmt.Println("\n--- Range over string (runes) ---")
	for idx, r := range s {
		fmt.Printf("index=%d rune=%c (%U)\n", idx, r, r)
	}

	fmt.Println("\n--- Reverse UTF-8 string ---")
	fmt.Println("Original :", s)
	fmt.Println("Reversed :", reverseString(s))

	// Byte vs rune count difference
	emoji := "A😀B"
	fmt.Printf("\n'%s': bytes=%d runes=%d\n", emoji, len(emoji), utf8.RuneCountInString(emoji))
}
