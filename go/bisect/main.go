package main

import (
	"cmp"
	"fmt"
)

// Search returning index of place to insert value v
func Search[T cmp.Ordered](a []T, v T) int {
	// Lemma:
	// all a[idx < l]:  a[idx] <= v
	// all a[idx >= r]: a[idx] >  v
	// final:
	// v = 3
	//       l   r
	//        \ /
	// 1 2 3 3 4
	// we have to return r to let new value appears in this posistion
	//         v----new
	// 1 2 3 3 3 4
	l := 0      // a[idx < l] — empty set
	r := len(a) // a[idx >= r] — empty set
	for l < r {
		m := (l + r) / 2 // better: l + (r-l)/2
		if a[m] <= v {
			l = m + 1 // it's important to use strict less for left, to do this trick
		} else {
			r = m
		}
	}
	return r
}

// FindAnyValue finds index of any element with given value
func FindAnyValue[T cmp.Ordered](a []T, v T) int {
	// Lemma:
	// all a[idx < l]: a[idx] < v
	// all a[idx > r]: a[idx] > v
	// final:
	// v = 3
	//   r l
	// 1 2 4 5
	//     l   r
	// 1 2 3 3 3 4
	l := 0
	r := len(a) - 1
	for l <= r {
		m := (l + r) / 2
		if a[m] < v {
			l = m + 1
		} else if a[m] > v {
			r = m - 1
		} else {
			return m
		}
	}
	return -1
}

// FindValueRange returns leftmost index of v and rightmost index of v, or (-1, -1) if no v
func FindValueRange[T cmp.Ordered](a []T, v T) (int, int) {
	// 1. Find left most
	// Lemma:
	// all a[idx < l]: a[idx] < v
	// all a[idx >= r]: a[idx] >= v
	// final:
	// v = 3
	//   l   r
	//    \ /
	// 1 2 3 3 4
	l := 0
	r := len(a)
	for l < r {
		m := (l + r) / 2
		if a[m] < v {
			l = m + 1
		} else {
			r = m
		}
	}
	if r == len(a) { // all elements are less than v
		return -1, -1
	}
	if a[r] != v { // no v
		return -1, -1
	}
	left := r
	// 2. Find right most
	// Lemma:
	// all a[idx < l]: a[idx] <= v
	// all a[idx >= r]: a[idx] > v
	// final:
	// v = 3
	//     l   r
	//      \ /
	// 1 2 3 3 4
	l = 0
	r = len(a)
	for l < r {
		m := (l + r) / 2
		if a[m] <= v { // diff!
			l = m + 1
		} else {
			r = m
		}
	}
	return left, r - 1
}

func main() {
	cases := [][]int{
		nil,                   // nil
		{},                    // empty
		{1},                   // one less
		{3},                   // one equal
		{5},                   // one great
		{1, 1, 2},             // all less
		{1, 2, 4},             // not present
		{1, 3, 4},             // present in center
		{3, 4, 5},             // present left
		{1, 2, 3},             // present right
		{1, 3, 3, 4},          // present many times in center
		{3, 3, 4, 5},          // present many times left
		{1, 2, 3, 3},          // present many times right
		{4, 5, 6},             // all high
		{1, 3, 3, 3, 3, 3, 4}, // long match
	}
	title("Classic binary search")
	for _, a := range cases {
		i := Search(a, 3)
		prn(a, i)
	}
	title("Find any by value")
	for _, a := range cases {
		i := FindAnyValue(a, 3)
		prn(a, i)
	}
	title("Find range by value")
	for _, a := range cases {
		i, j := FindValueRange(a, 3)
		prn(a, i, j)
	}
}

// helpers

func prn[T any](a []T, idx ...int) {
	clrs := make([][2]string, len(a)+1)
	for _, i := range idx {
		if i >= 0 && i <= len(a) {
			clrs[i] = [2]string{"\033[41;93;1m", "\033[0m"}
		}
	}
	for i, v := range a {
		c := clrs[i]
		fmt.Printf("%s%v%s ", c[0], v, c[1])
	}
	c := clrs[len(a)]
	fmt.Printf("%s%v%s (i=%d)\n", c[0], "-", c[1], idx)
}

func title(s string) {
	fmt.Printf("\033[92;1m%s\033[0m\n", s)
}
