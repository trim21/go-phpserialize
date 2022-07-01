PHP `serialize()` and `unserialize()` (in future) for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...eta

Advantage:

Low memory allocation and fast:

```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor
BenchmarkMarshal_all-16                                 14072532                92.43 ns/op          184 B/op          2 allocs/op
BenchmarkMarshal_map_concrete_types/len-1-16            43740704                26.50 ns/op           16 B/op          1 allocs/op
BenchmarkMarshal_map_concrete_types/len-10-16           14448308                92.70 ns/op           96 B/op          1 allocs/op
BenchmarkMarshal_map_concrete_types/len-100-16           1503818               795.3 ns/op          1025 B/op          1 allocs/op
BenchmarkMarshal_slice_concrete_types/len-1-16          50314042                25.24 ns/op           40 B/op          2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-10-16          8027953               143.3 ns/op           120 B/op          2 allocs/op
BenchmarkMarshal_slice_concrete_types/len-100-16          924982              1387 ns/op            1049 B/op          2 allocs/op
Benchmark_elliotchance_phpserialize_marshal-16            189254              6064 ns/op            3627 B/op        105 allocs/op
PASS
ok      github.com/trim21/go-phpserialize       11.304s
```

Disadvantage:

heavy usage of `unsafe`.

Limitation:

1. You can't use `interface` (working on it, not done yet)
2. Marshal go `struct`, `map` into php array and array only, php object is not supported.\
3. No `omitempty` support (yet).

example:

```golang
package main

import (
	"fmt"

	"github.com/trim21/go-phpserialize"
)

type Inner struct {
	V int    `php:"v"`
	S string `php:"a long string name replace field name"`
}

type With struct {
	Users []User `php:"users"`
	Obj   Inner  `php:"obj"`
}

type User struct {
	ID   uint32 `php:"id"`
	Name string `php:"name"`
}

func main() {
	var data = With{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		Obj: Inner{V: 2, S: "vvv"},
	}
	var b, err = phpserialize.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
```

Marshaler is heavily inspired by https://github.com/goccy/go-json

this is different with https://github.com/elliotchance/phpserialize
