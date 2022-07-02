```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor             
BenchmarkMarshal_concrete_types/bool_true-16      	28589058	        40.13 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/bool_false-16     	29209375	        40.62 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/int8-16           	27247089	        44.32 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/int16-16          	26642510	        45.56 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/int32-16          	24467924	        45.66 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/int64-16          	26642274	        46.60 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/int-16            	23735960	        49.36 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/uint8-16          	22863892	        53.79 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/uint16-16         	23055906	        53.07 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/uint32-16         	26442318	        48.00 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/uint64-16         	24769130	        45.59 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/uint-16           	21492241	        52.00 ns/op	       4 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/float32-16        	11580998	       103.0 ns/op	       8 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/float64-16        	 9827574	       122.0 ns/op	       8 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/string-16         	 5440424	       221.7 ns/op	      32 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/simple_slice-16   	 9489571	       129.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/struct_slice-16   	 5833468	       200.1 ns/op	     128 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/struct_with_map_ptr-16         	 6274210	       191.6 ns/op	      64 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/struct_with_map_embed-16       	 6656156	       181.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/nil_map-16                     	28117135	        42.47 ns/op	       2 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/nested_struct_not_anonymous-16 	 9668304	       124.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/nested_struct_anonymous-16     	12065248	       100.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/complex_object-16              	 2712049	       438.3 ns/op	     256 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/nested_map-16                  	 4915878	       244.5 ns/op	      32 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/map[type]any(map)-16           	 3724242	       319.4 ns/op	      32 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/map[type]any(slice)-16         	 4898326	       238.8 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/map[type]any(struct)-16        	 2583674	       463.0 ns/op	      64 B/op	       3 allocs/op
BenchmarkMarshal_concrete_types/generic[int]-16                	16879725	        72.56 ns/op	      24 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/generic[struct]-16             	11460912	       103.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/generic[map]-16                	 6578026	       180.1 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal_concrete_types/generic[slice]-16              	 6413346	       187.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkMarshal_interface/bool_true-16                        	 5721961	       207.5 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/bool_false-16                       	 5503890	       217.4 ns/op	      96 B/op	       2 allocs/op
BenchmarkMarshal_interface/int8-16                             	 6927450	       173.2 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/int16-16                            	 6617618	       180.9 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/int32-16                            	 6551437	       181.6 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/int64-16                            	 6586975	       181.3 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/int-16                              	 7186837	       166.1 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/uint8-16                            	 6641234	       179.5 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/uint16-16                           	 6234151	       186.7 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/uint32-16                           	 6429265	       187.6 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/uint64-16                           	 6194385	       186.8 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/uint-16                             	 6984451	       173.1 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/float32-16                          	 4498578	       269.2 ns/op	      96 B/op	       2 allocs/op
BenchmarkMarshal_interface/float64-16                          	 4213357	       283.9 ns/op	      96 B/op	       2 allocs/op
BenchmarkMarshal_interface/string-16                           	 3360663	       355.7 ns/op	     112 B/op	       2 allocs/op
BenchmarkMarshal_interface/simple_slice-16                     	 3423482	       342.3 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_interface/struct_slice-16                     	 2773564	       427.1 ns/op	     208 B/op	       2 allocs/op
BenchmarkMarshal_interface/struct_with_map_ptr-16              	 1458866	       822.1 ns/op	     160 B/op	       4 allocs/op
BenchmarkMarshal_interface/struct_with_map_embed-16            	 1846052	       654.0 ns/op	     136 B/op	       3 allocs/op
BenchmarkMarshal_interface/nil_map-16                          	 5694903	       209.8 ns/op	      80 B/op	       2 allocs/op
BenchmarkMarshal_interface/nested_struct_not_anonymous-16      	 1535758	       782.4 ns/op	     192 B/op	       6 allocs/op
BenchmarkMarshal_interface/nested_struct_anonymous-16          	 1610706	       741.7 ns/op	     168 B/op	       5 allocs/op
BenchmarkMarshal_interface/complex_object-16                   	  428176	      2821 ns/op	     801 B/op	      29 allocs/op
BenchmarkMarshal_interface/nested_map-16                       	 2530537	       470.7 ns/op	     112 B/op	       2 allocs/op
BenchmarkMarshal_interface/map[type]any(map)-16                	 2126220	       564.3 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_interface/map[type]any(slice)-16              	 2415004	       495.3 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_interface/map[type]any(struct)-16             	 1621114	       742.1 ns/op	     160 B/op	       4 allocs/op
BenchmarkMarshal_interface/generic[int]-16                     	 3174907	       375.9 ns/op	     120 B/op	       3 allocs/op
BenchmarkMarshal_interface/generic[struct]-16                  	 1771732	       681.6 ns/op	     168 B/op	       5 allocs/op
BenchmarkMarshal_interface/generic[map]-16                     	 2198209	       542.6 ns/op	     136 B/op	       3 allocs/op
BenchmarkMarshal_interface/generic[slice]-16                   	 1767421	       675.3 ns/op	     232 B/op	       7 allocs/op
BenchmarkMarshal_map_concrete_types/len-1-16                   	53885358	        23.14 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_concrete_types/len-10-16                  	17216864	        79.40 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_concrete_types/len-100-16                 	 2001140	       647.0 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_map_as_interface/len-1-16                     	29241903	        38.39 ns/op	      48 B/op	       2 allocs/op
BenchmarkMarshal_map_as_interface/len-10-16                    	14597563	        91.78 ns/op	     128 B/op	       2 allocs/op
BenchmarkMarshal_map_as_interface/len-100-16                   	 1967404	       648.6 ns/op	    1042 B/op	       2 allocs/op
BenchmarkMarshal_map_with_interface_value/len-1-16             	40097972	        25.43 ns/op	      16 B/op	       1 allocs/op
BenchmarkMarshal_map_with_interface_value/len-10-16            	14987298	        88.10 ns/op	      96 B/op	       1 allocs/op
BenchmarkMarshal_map_with_interface_value/len-100-16           	 1673568	       733.5 ns/op	    1025 B/op	       1 allocs/op
BenchmarkMarshal_slice_concrete_types/len-1-16                 	55472026	        20.66 ns/op	      40 B/op	       2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-10-16                	 8798563	       134.0 ns/op	     120 B/op	       2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-100-16               	  884426	      1297 ns/op	    1049 B/op	       2 allocs/op
BenchmarkMarshal_slice_interface/len-1-16                      	30100158	        36.29 ns/op	      72 B/op	       3 allocs/op
BenchmarkMarshal_slice_interface/len-10-16                     	 8270449	       144.7 ns/op	     152 B/op	       3 allocs/op
BenchmarkMarshal_slice_interface/len-100-16                    	  922246	      1317 ns/op	    1065 B/op	       3 allocs/op
Benchmark_marshal_compare-16                                   	 3554985	       335.0 ns/op	     208 B/op	       2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16                 	  217986	      5526 ns/op	    3451 B/op	     101 allocs/op
PASS
ok  	github.com/trim21/go-phpserialize	116.049s
?   	github.com/trim21/go-phpserialize/internal/encoder	[no test files]
```


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
