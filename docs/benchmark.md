```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor
BenchmarkMarshal_all-16                                 10343846               120.3 ns/op           352 B/op          2 allocs/op
BenchmarkAll_concrete_types/bool_true-16                29649908                39.92 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/bool_false-16               29594700                40.05 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/int8-16                     25780116                46.88 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/int16-16                    25268157                47.33 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/int32-16                    25548385                48.27 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/int64-16                    25237653                48.81 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/int-16                      24089805                49.58 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/uint8-16                    24690290                48.01 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/uint16-16                   24696998                45.59 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/uint32-16                   25757484                45.78 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/uint64-16                   27488071                43.86 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/uint-16                     25828390                47.01 ns/op            4 B/op          1 allocs/op
BenchmarkAll_concrete_types/float32-16                  12098773                97.98 ns/op            8 B/op          1 allocs/op
BenchmarkAll_concrete_types/float64-16                  10542684               114.2 ns/op             8 B/op          1 allocs/op
BenchmarkAll_concrete_types/string-16                    5314854               225.9 ns/op            32 B/op          1 allocs/op
BenchmarkAll_concrete_types/simple_slice-16              9699446               124.7 ns/op            48 B/op          1 allocs/op
BenchmarkAll_concrete_types/struct_slice-16              5826050               205.2 ns/op           128 B/op          1 allocs/op
BenchmarkAll_concrete_types/struct_with_map_ptr-16               5474661               218.4 ns/op            64 B/op          1 allocs/op
BenchmarkAll_concrete_types/struct_with_map_embed-16             5782441               203.4 ns/op            48 B/op          1 allocs/op
BenchmarkAll_concrete_types/nil_map-16                          15910476                72.67 ns/op           24 B/op          1 allocs/op
BenchmarkAll_concrete_types/nested_struct_not_anonymous-16       9980430               119.8 ns/op            64 B/op          1 allocs/op
BenchmarkAll_concrete_types/nested_struct_anonymous-16          11537196               104.9 ns/op            48 B/op          1 allocs/op
BenchmarkAll_concrete_types/complex_object-16                    2506570               480.9 ns/op           256 B/op          1 allocs/op
BenchmarkAll_interface/bool_true-16                              5767284               208.5 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/bool_false-16                             5434452               223.8 ns/op            96 B/op          2 allocs/op
BenchmarkAll_interface/int8-16                                   6783025               177.2 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/int16-16                                  6507496               182.9 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/int32-16                                  6545419               182.3 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/int64-16                                  6614941               181.7 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/int-16                                    7110158               168.1 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/uint8-16                                  6344419               182.4 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/uint16-16                                 6091429               191.1 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/uint32-16                                 6173038               187.8 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/uint64-16                                 6415170               189.9 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/uint-16                                   6828225               174.2 ns/op            80 B/op          2 allocs/op
BenchmarkAll_interface/float32-16                                4546837               262.3 ns/op            96 B/op          2 allocs/op
BenchmarkAll_interface/float64-16                                4351116               275.4 ns/op            96 B/op          2 allocs/op
BenchmarkAll_interface/string-16                                 3344209               360.9 ns/op           112 B/op          2 allocs/op
BenchmarkAll_interface/simple_slice-16                           2695722               447.7 ns/op           168 B/op          7 allocs/op
BenchmarkAll_interface/struct_slice-16                            797856              1325 ns/op             408 B/op         22 allocs/op
BenchmarkAll_interface/nested_struct_not_anonymous-16            1442464               825.0 ns/op           208 B/op          7 allocs/op
BenchmarkAll_interface/nested_struct_anonymous-16                1488273               813.2 ns/op           200 B/op          7 allocs/op
BenchmarkMarshal_map_concrete_types/len-1-16                    45631018                24.81 ns/op           16 B/op          1 allocs/op
BenchmarkMarshal_map_concrete_types/len-10-16                   14372926                82.31 ns/op           96 B/op          1 allocs/op
BenchmarkMarshal_map_concrete_types/len-100-16                   1833901               651.4 ns/op          1026 B/op          1 allocs/op
BenchmarkMarshal_slice_concrete_types/len-1-16                  48738880                22.39 ns/op           40 B/op          2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-10-16                  8251402               142.8 ns/op           120 B/op          2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-100-16                  878926              1379 ns/op            1049 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16                    203700              5764 ns/op            3627 B/op        105 allocs/op
PASS
ok      github.com/trim21/go-phpserialize       69.907s
```
