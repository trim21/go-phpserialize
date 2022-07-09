# go-phpserialize

PHP `serialize()` and `unserialize()` for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...etc.

## Use case:

You serialize all data into php array only, and decoding php serialized array only

decoding php object (or class) is not supported yet.

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

### Disadvantage:

heavy usage of `unsafe`.

#### Limitation:

1. No `omitempty` support (yet).
2. `Marshaling` Anonymous Struct field (embedding struct) working like named field, `Unmarshal` works fine.

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
	Users   []User `php:"users,omitempty"`
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

php serialized array has a length prefix `a:1:{i:0;s:3:"one";}`, when decoding php serialized array into go `slice` or
go `map`,
`go-phpserialize` may call golang's `make()` to create a map or slice with given length.

So a malicious input like `a:100000000:{}` may become `make([]T, 100000000)` and consume high memory.

If you have to decode some un-trusted bytes, make sure only decode them into fixed-length golang array or struct,
never decode them to `interface`, `slice` or `map`.
