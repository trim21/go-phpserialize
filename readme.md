PHP `serialize()` and `unserialize()` (in future) for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...eta

Advantage:

Low memory allocation and fast:

```text
goos: windows
goarch: amd64
pkg: github.com/trim21/go-phpserialize
cpu: AMD Ryzen 7 5800X 8-Core Processor

BenchmarkMarshal_map_concrete_types/len-1-16             6954546               169.4 ns/op           144 B/op          2 allocs/op
BenchmarkMarshal_map_concrete_types/len-10-16            2814568               423.8 ns/op           224 B/op          2 allocs/op
BenchmarkMarshal_map_concrete_types/len-100-16            399626              3058 ns/op            1152 B/op          2 allocs/op
BenchmarkMarshal_map_concrete_types/len-1000-16            31783             37477 ns/op           12425 B/op          2 allocs/op
```

Disadvantage:

heavy usage of `unsafe`.

Limitation:

1. You can't use `interface` (working on it, not done yet)
2. Marshal go `struct`, `map` into php array and array only, php object is not supported.

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

you will
get `a:2:{s:5:"users";a:2:{i:0;a:2:{s:2:"id";i:1;s:4:"name";s:3:"sai";}i:1;a:2:{s:2:"id";i:2;s:4:"name";s:6:"trim21";}}s:3:"obj";a:2:{s:1:"v";i:2;s:37:"a long string name replace field name";s:3:"vvv";}}`

this is different with https://github.com/elliotchance/phpserialize
