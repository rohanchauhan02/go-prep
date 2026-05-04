//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"unicode/utf8"
)

// TODO: reverse a UTF-8 string correctly (by runes, not bytes)
func reverseString(s string) string {
	// TODO: convert to []rune, reverse, convert back
	return ""
}

// TODO: count bytes vs runes — they differ for multi-byte chars
func bytesVsRunes(s string) (bytes, runes int) {
	// TODO: use len(s) for bytes, utf8.RuneCountInString for runes
	return
}

func main() {
	s := "Hello, 世界"

	b, r := bytesVsRunes(s)
	fmt.Printf("bytes=%d runes=%d
", b, r)
	// Expected: bytes=13 runes=9

	fmt.Println("reversed:", reverseString(s))
	// Expected: 界世 ,olleH

	// Range over string gives runes
	fmt.Println("
Range (rune indices):")
	for idx, ch := range s {
		fmt.Printf("  [%d] %c
", idx, ch)
	}
	// Note: indices jump by byte size (世 = 3 bytes)

	fmt.Println("
utf8 package:")
	fmt.Println("RuneCount:", utf8.RuneCountInString(s))
}
