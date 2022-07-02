# benchmark

```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor             
BenchmarkMarshal_type/simple_slice-16             	             9865743	       123.1 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_ifce/simple_slice-16                        	 3536656	       338.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_type/struct_slice-16                        	 5961864	       198.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkMarshal_ifce/struct_slice-16                        	 2821891	       426.0 ns/op	     208 B/op	       2 allocs/op
BenchmarkMarshal_type/nested_struct_not_anonymous-16         	 9823286	       121.8 ns/op	      64 B/op	       1 allocs/op
BenchmarkMarshal_type/nested_struct_anonymous-16             	12110666	        98.68 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_type/complex_object-16                      	 2607030	       455.6 ns/op	     256 B/op	       1 allocs/op
BenchmarkMarshal_ifce/complex_object-16                      	  428196	      2739 ns/op	     801 B/op	      29 allocs/op
BenchmarkMarshal_map_type/len-1-16                           	52128130	        21.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-10-16                          	16340535	        77.92 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-100-16                         	 1996768	       637.9 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_map_type/len-1000-16                        	  171327	      7428 ns/op	   12320 B/op	       1 allocs/op
BenchmarkMarshal_map_as_ifce/len-1-16                        	31290742	        37.91 ns/op	      48 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-10-16                       	14223368	        89.86 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-100-16                      	 2000554	       621.3 ns/op	    1042 B/op	       2 allocs/op
BenchmarkMarshal_map_as_ifce/len-1000-16                     	  128023	      7885 ns/op	   12447 B/op	       2 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-1-16                	48151060	        25.95 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-10-16               	14288197	        88.11 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-100-16              	 1712043	       837.7 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_map_with_ifce_value/len-1000-16             	  130514	     10660 ns/op	   12314 B/op	       1 allocs/op
BenchmarkMarshal_slice_of_value/len-1-16                     	45599980	        26.44 ns/op	      56 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_value/len-10-16                    	 8269929	       144.4 ns/op	     136 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_value/len-100-16                   	  938629	      1343 ns/op	    1049 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_value/len-1000-16                  	   86728	     13799 ns/op	   12431 B/op	       2 allocs/op
BenchmarkMarshal_ifce_slice_as_value/len-1-16                	20397339	        60.85 ns/op	     136 B/op	       7 allocs/op
BenchmarkMarshal_ifce_slice_as_value/len-10-16               	13426633	        93.19 ns/op	     216 B/op	       7 allocs/op
BenchmarkMarshal_ifce_slice_as_value/len-100-16              	 2299296	       451.6 ns/op	    1129 B/op	       7 allocs/op
BenchmarkMarshal_ifce_slice_as_value/len-1000-16             	  218365	      5507 ns/op	   12541 B/op	       7 allocs/op
BenchmarkMarshal_slice_of_type/len-1-16                      	25546209	        43.06 ns/op	     104 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_type/len-10-16                     	 4940810	       232.8 ns/op	     504 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_type/len-100-16                    	  552106	      2180 ns/op	    4913 B/op	       2 allocs/op
BenchmarkMarshal_slice_of_type/len-1000-16                   	   47683	     25289 ns/op	   58884 B/op	       2 allocs/op
BenchmarkMarshal_ifce_slice_of_type/len-1-16                 	 5300565	       227.9 ns/op	     520 B/op	      24 allocs/op
BenchmarkMarshal_ifce_slice_of_type/len-10-16                	 3100856	       382.5 ns/op	     921 B/op	      24 allocs/op
BenchmarkMarshal_ifce_slice_of_type/len-100-16               	  578076	      2145 ns/op	    5339 B/op	      24 allocs/op
BenchmarkMarshal_ifce_slice_of_type/len-1000-16              	   54580	     22132 ns/op	   59057 B/op	      24 allocs/op
BenchmarkMarshal_ifce_slice_of_ifce/len-1-16                 	13511793	        89.73 ns/op	     136 B/op	       5 allocs/op
BenchmarkMarshal_ifce_slice_of_ifce/len-10-16                	 2116438	       563.4 ns/op	     680 B/op	      23 allocs/op
BenchmarkMarshal_ifce_slice_of_ifce/len-100-16               	  218760	      5733 ns/op	    6534 B/op	     203 allocs/op
BenchmarkMarshal_ifce_slice_of_ifce/len-1000-16              	   20196	     58888 ns/op	   75302 B/op	    2003 allocs/op
Benchmark_marshal_compare-16                                 	 3548409	       335.0 ns/op	     208 B/op	       2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16               	  217281	      5529 ns/op	    3451 B/op	     101 allocs/op

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
