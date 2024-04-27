package pkg_test

import (
	"fmt"
	"testing"

	pkg "x"
)

func BenchmarkSum(b *testing.B) {
	for i, f := range []func(int) int{
		pkg.Sum1,
		pkg.Sum2,
		pkg.Sum3,
		pkg.Sum4,
	} {
		b.Run(fmt.Sprintf("f%d", i+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f(1000000)
			}
		})
	}
}
