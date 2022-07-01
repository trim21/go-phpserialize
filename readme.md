# go-phpserialize

PHP `serialize()` and `unserialize()` (in future) for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...etc.

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

### Disadvantage:

heavy usage of `unsafe`.

#### Limitation:

1. Marshal go `struct`, `map` into php array and array only, php object is not supported.
2. `interface` contain any `map` (in progress) or `interface` as `map`'s value type
3. No `omitempty` support (yet).
4. Anonymous Struct field (embedding struct) working like named field.

If any of this Limitation affect you (except `1.`), please create a PR to let me know.

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
