# go-phpserialize

PHP `serialize()` and `unserialize()` (in future) for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...etc.

## Use case:

You serialize all data into php array only, php object (or stdClass) is not supported.

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

#### Performance Hint

Encoder will try to build an optimized path for a type. When encoding `interface`,
encoder will fall back to reflect, which is much slower. 

If you care about performance, you should avoid using interface.

Using type is 2x faster than interface in average.

In the worst condition, it may be 8x slower (or more).

```text
BenchmarkMarshal_type/struct_with_all-16      2814640        431.0 ns/op      256 B/op     1 allocs/op
BenchmarkMarshal_ifce/struct_with_all-16       374654         3247 ns/op      849 B/op    35 allocs/op
```

### Disadvantage:

heavy usage of `unsafe`.

#### Limitation:

1. No `omitempty` support (yet).
2. Anonymous Struct field (embedding struct) working like named field.

If any of these limitations affect you, please create an issue to let me know.

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
	Users   []User `php:"users"`
	Obj     Inner  `php:"obj"`
  	Ignored bool   `php:"-"`
}

type User struct {
	ID   uint32 `php:"id,string"`
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

Heavily inspired by https://github.com/goccy/go-json


## Security

TL;DR: Don't unmarshal content you can't trust.

Attackers may consume large memory with very few bytes.

php serialized array has a length prefix `a:1:{i:0;s:3:"one";}`, when decoding php serialized array into go `slice` or go `map`, 
`go-phpserialize` may call golang's `make()` to create a map or slice with given length.

So a malicious input like `a:100000000:{}` may become `make([]T, 100000000)` and consume high memory.

If you have to decode some un-trusted bytes, make sure only decode then into fixed-length golang array or struct, 
never decode them to `interface`, `slice` or `map`.
