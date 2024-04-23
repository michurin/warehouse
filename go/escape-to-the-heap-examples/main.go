package main

func returnValue() int {
	n := 42 // stack
	return n
}

func returnPointer() *int {
	n := 42 // escapes to heap due to pointer
	return &n
}

func returnMap() map[int]int {
	m := map[int]int{42: 137} // escapes
	return m
}

func returnNestedMaps() map[int]any {
	t := map[int]int{1: 1} // does not escape
	_ = t
	m := map[int]int{42: 137} // escapes
	n := map[int]any{1: m}    // escapes
	return n
}

func returnSlice() []int {
	t := []int{1, 2}            // stack
	slice := []int{42}          // escapes to heap
	slice = append(slice, t...) // (just consume t)
	return slice
}

func returnArray() [1]int {
	return [1]int{42} // does not escape, returns by value
}

func largeArray() int {
	var largeArray [100000000]int // escapes! ridiculously large
	largeArray[42] = 42
	return largeArray[42]
}

func returnFunc() func() int {
	f := func() int { return 42 } // escaeps
	return f
}

func localChan() chan (int) {
	c := make(chan (int), 100000000) // stack; `c` itself, but not underlying buffers
	c <- 1
	return c
}

func main() { // go run -gcflags="-m -l" main.go
	_ = returnValue() // consume funcs
	_ = *returnPointer()
	_ = returnSlice()
	_ = returnMap()
	_ = returnNestedMaps()
	_ = returnArray()
	_ = largeArray()
	_ = returnFunc()()
	_ = localChan()
}
