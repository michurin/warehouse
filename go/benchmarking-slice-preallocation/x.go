package pkg

func nums1(n int) []int {
	r := []int(nil)
	for i := 1; i <= n; i++ {
		r = append(r, i)
	}
	return r
}

func Sum1(n int) int {
	a := nums1(n)
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}

func nums2(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i + 1
	}
	return r
}

func Sum2(n int) int {
	a := nums2(n)
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}

func Sum3(n int) int {
	a := []int(nil)
	for i := 1; i <= n; i++ {
		a = append(a, i)
	}
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}

func Sum4(n int) int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i + 1
	}
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}
