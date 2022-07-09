# benchmark

you can see full result at [benchmark.txt](./benchmark.txt)

compare with https://github.comm/elliotchance/phpserialize (will have to use map to marshal php array)

```text
Benchmark_marshal_compare-16                             2464408               484.1 ns/op           208 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16            222019              5316 ns/op            3450 B/op        101 allocs/op
```

you can find source at [bench_test.go](../bench/bench_test.go)
