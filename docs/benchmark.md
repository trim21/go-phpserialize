# benchmark

```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor
BenchmarkMarshal_type/simple_slice-16                       	 9698301	       124.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_type/struct_slice-16                       	 6022387	       198.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkMarshal_type/nested_struct_anonymous-16             	11980591	        99.63 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_type/complex_object-16                      	 2744300	       441.8 ns/op	     256 B/op	       1 allocs/op
BenchmarkMarshal_ifce/simple_slice-16                        	 3572583	       335.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_ifce/struct_slice-16                        	 2793996	       432.9 ns/op	     208 B/op	       2 allocs/op
BenchmarkMarshal_ifce/complex_object-16                      	  444032	      2708 ns/op	     801 B/op	      29 allocs/op
BenchmarkMarshal_map_type/len-1-16                           	46198089	        21.74 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-10-16                          	17399044	        71.77 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-100-16                         	 1546312	       691.2 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-1000-16                        	  175656	      9016 ns/op	   12316 B/op	       1 allocs/op
BenchmarkMarshal_map_as_ifce/len-1-16                        	28166236	        40.92 ns/op	      48 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-10-16                       	13260034	        93.77 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-100-16                      	 1905309	       628.7 ns/op	    1042 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-1000-16                     	  126595	      8394 ns/op	   12443 B/op	       2 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-1-16                	57196239	        22.99 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-10-16               	14403390	        84.54 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-100-16              	 1729622	       834.3 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-1000-16             	  108784	     10036 ns/op	   12317 B/op	       1 allocs/op
BenchmarkMarshal_slice_as_type/len-1-16                      	54770260	        20.90 ns/op	      40 B/op	       2 allocs/op
BenchmarkMarshal_slice_as_type/len-10-16                     	 8606517	       139.4 ns/op	     120 B/op	       2 allocs/op
BenchmarkMarshal_slice_as_type/len-100-16                    	  857528	      1366 ns/op	    1049 B/op	       2 allocs/op
BenchmarkMarshal_slice_as_type/len-1000-16                   	   85200	     14123 ns/op	   12445 B/op	       2 allocs/op
BenchmarkMarshal_slice_as_ifce/len-1-16                      	27253401	        37.00 ns/op	      72 B/op	       3 allocs/op
BenchmarkMarshal_slice_as_ifce/len-10-16                     	 8070274	       145.3 ns/op	     152 B/op	       3 allocs/op
BenchmarkMarshal_slice_as_ifce/len-100-16                    	  903620	      1320 ns/op	    1065 B/op	       3 allocs/op
BenchmarkMarshal_slice_as_ifce/len-1000-16                   	   87000	     13844 ns/op	   12462 B/op	       3 allocs/op
BenchmarkMarshal_slice_of_type/len-1-16                      	 5289433	       228.1 ns/op	     520 B/op	      24 allocs/op
BenchmarkMarshal_slice_of_type/len-10-16                     	 3070195	       407.0 ns/op	     921 B/op	      24 allocs/op
BenchmarkMarshal_slice_of_type/len-100-16                    	  526257	      2218 ns/op	    5339 B/op	      24 allocs/op
BenchmarkMarshal_slice_of_type/len-1000-16                   	   55438	     21777 ns/op	   58792 B/op	      24 allocs/op
BenchmarkMarshal_slice_of_ifce/len-1-16                      	13607493	        87.47 ns/op	     136 B/op	       5 allocs/op
BenchmarkMarshal_slice_of_ifce/len-10-16                     	 2109963	       558.5 ns/op	     680 B/op	      23 allocs/op
BenchmarkMarshal_slice_of_ifce/len-100-16                    	  227103	      5396 ns/op	    6535 B/op	     203 allocs/op
BenchmarkMarshal_slice_of_ifce/len-1000-16                   	   20322	     58138 ns/op	   75204 B/op	    2003 allocs/op
Benchmark_marshal_compare-16                                 	 3547075	       333.5 ns/op	     208 B/op	       2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16               	  216318	      5457 ns/op	    3451 B/op	     101 allocs/op
PASS
```

you can see full result at [benchmark.txt](./benchmark.txt)

compare with https://github.comm/elliotchance/phpserialize (will have to use map to marshal php array)

```text
Benchmark_marshal_compare-16                                     3649764               327.4 ns/op           208 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16                    220365              5382 ns/op            3451 B/op        101 allocs/op
```

```golang
func Benchmark_marshal_compare(b *testing.B) {
type Obj struct {
V int `php:"v"`
S int `php:"s"`
}

type TestData struct {
Users []User `php:"users"`
Obj   Obj    `php:"obj"`
}

var data = TestData{
Users: []User{
{ID: 1, Name: "sai"},
{ID: 2, Name: "trim21"},
},
Obj: Obj{V: 2, S: 3},
}

for i := 0; i < b.N; i++ {
phpserialize.Marshal(data)
}
}

func Benchmark_elliotchance_phpserialize_marshal(b *testing.B) {
var data = map[any]any{
"users": []map[any]any{
{"id": 1, "name": "sai"},
{"id": 2, "name": "trim21"},
},
"obj": map[string]int{"v": 2, "s": 3},
}
for i := 0; i < b.N; i++ {
elliotchance_phpserialize.Marshal(data, nil)
}
}
```
