package main

import (
	"fmt"
	"unicode"
)

func main() {
	var r *unicode.RangeTable
	// r = unicode.Latin
	// r = unicode.Digit
	r = unicode.Punct
	// r = unicode.Cypriot
	// r = unicode.Cyrillic
	// r = unicode.GraphicRanges[1]
	for i, x := range r.R16 {
		fmt.Printf("R16%3[1]d %[2]c(%[2]d) - %[3]c(%[3]d): ", i, x.Lo, x.Hi)
		for c := max(x.Lo, 32); c <= min(x.Hi, 127); c++ {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	for i, x := range r.R32 {
		fmt.Printf("R32%3[1]d %[2]c(%[2]d) - %[3]c(%[3]d): ", i, x.Lo, x.Hi)
		for c := max(x.Lo, 32); c <= min(x.Hi, 127); c++ {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	fmt.Println()
}
