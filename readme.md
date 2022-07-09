# go-phpserialize

PHP `serialize()` and `unserialize()` for Go.

Support All go type including `map`, `slice`, `strcut`, and simple type like `int`, `uint` ...etc.

## supported and tested go version

- 1.17
- 1.18 (go 1.18+ with have extra generic type tests)

## Use case:

You serialize all data into php array only. Decoding from php serialized array or class are both supported.

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

### Disadvantage:

heavy usage of `unsafe`.

#### Limitation:

`Marshaling` Anonymous Struct field (embedding struct) working like named field, `Unmarshal` works fine.

## example

### Marshal

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

you can un-serialize it in php to

```php
array(
    "users" => array(
        0 => array(
            "id" => "1",
            "name" => "sai",
        ),
        1 => array(
            "id" => "2",
            "name" => "trim21",
        ),
    ),
    "obj" => array(
        "v" => 2,
        "a long string name replace field name" => "vvv",
    ),
);
```

### Unmarshal

```golang
package main

import (
	"encoding/json"
	"fmt"

	"github.com/trim21/go-phpserialize"
)

func main() {
	var v struct {
		Value map[string]string `php:"value" json:"value"`
	}
	raw := `a:1:{s:5:"value";a:5:{s:3:"one";s:1:"1";s:3:"two";s:1:"2";s:5:"three";s:1:"3";s:4:"four";s:1:"4";s:4:"five";s:1:"5";}}`

	err := phpserialize.Unmarshal([]byte(raw), &v)
	if err != nil {
		panic(err)
	}

	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
```

you will see

```json
{
  "value": {
    "five": "5",
    "four": "4",
    "one": "1",
    "three": "3",
    "two": "2"
  }
}
```

`any` type will be decoded to `map[any]any`, where key maybe `int64` or `string`, and value can be anything.

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

## License

MIT License
