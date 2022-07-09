# benchmark

you can see full result at [benchmark.txt](./benchmark.txt)

compare with https://github.com/elliotchance/phpserialize (will have to use map to marshal php array)

```text
goos: windows
goarch: amd64
pkg: bench
cpu: AMD Ryzen 7 5800X 8-Core Processor
Benchmark_marshal_compare-16                             4007931               298.1 ns/op           208 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16            226204              5227 ns/op            3450 B/op        101 allocs/op
PASS
ok      bench   2.886s
```

you can find source at [bench_test.go](../bench/bench_test.go)
