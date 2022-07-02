# go-phpserialize

PHP `serialize()` and `unserialize()` (in future) for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...etc.

## Use case:

You serialize all data into php array only, php object (or stdClass) is not supported.

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

#### Performance Hint

Due to the technology encoder is using, You should avoid using interface.

Using type is 2x faster than interface.

In the worst condition, it will slow 8x slower.

```text
BenchmarkMarshal_type/complex_object-16            	 2744300	       441.8 ns/op	     256 B/op	       1 allocs/op
BenchmarkMarshal_ifce/complex_object-16            	  444032	      2708 ns/op	     801 B/op	      29 allocs/op
```

### Disadvantage:

heavy usage of `unsafe`.

#### Limitation:

1. No `omitempty` support (yet).
2. Anonymous Struct field (embedding struct) working like named field.

If any of these limitations affect you, please create a issue to let me know.

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
