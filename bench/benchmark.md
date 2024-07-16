# benchmark

you can see full result at [benchmark.txt](./benchmark.txt)

compare with https://github.com/elliotchance/phpserialize (will have to use map to marshal php array)

with go1.22

```text
goos: windows
goarch: amd64
pkg: bench
cpu: AMD Ryzen 7 5800X 8-Core Processor             
Benchmark_marshal_compare-16                             3612345               326.3 ns/op           208 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16            239784              5260 ns/op            3519 B/op        101 allocs/op
PASS
ok      bench   2.850s
```

you can find source at [bench_test.go](../bench/bench_test.go)
