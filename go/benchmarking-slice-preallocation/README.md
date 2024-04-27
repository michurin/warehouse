# Benchmarking slice preallocation

```
go test -bench=. -count=1 .
```

```
goos: linux
goarch: amd64
pkg: escape-to-the-heap-examples/pkg
cpu: 11th Gen Intel(R) Core(TM) i7-11700 @ 2.50GHz
BenchmarkSum/f1-16            34          35709854 ns/op
BenchmarkSum/f2-16           229           5146334 ns/op
BenchmarkSum/f3-16            31          39921890 ns/op
BenchmarkSum/f4-16           229           5313052 ns/op
PASS
```
