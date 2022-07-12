package phpserialize_test

import (
	"fmt"

	"github.com/trim21/go-phpserialize"
)

func ExampleMarshal() {
	type User struct {
		ID   uint32 `php:"id,string"`
		Name string `php:"name"`
	}

	type Inner struct {
		V int    `php:"v"`
		S string `php:"a long string name replace field name"`
	}

	type With struct {
		Users   []User `php:"users,omitempty"`
		Obj     Inner  `php:"obj"`
		Ignored bool   `php:"-"`
	}

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
	// Output: a:2:{s:5:"users";a:2:{i:0;a:2:{s:2:"id";s:1:"1";s:4:"name";s:3:"sai";}i:1;a:2:{s:2:"id";s:1:"2";s:4:"name";s:6:"trim21";}}s:3:"obj";a:2:{s:1:"v";i:2;s:37:"a long string name replace field name";s:3:"vvv";}}
}

func ExampleUnmarshal() {
	var v struct {
		Value map[string]string `php:"value" json:"value"`
	}
	raw := `a:1:{s:5:"value";a:5:{s:3:"one";s:1:"1";s:3:"two";s:1:"2";s:5:"three";s:1:"3";s:4:"four";s:1:"4";s:4:"five";s:1:"5";}}`

	err := phpserialize.Unmarshal([]byte(raw), &v)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.Value["five"])
	// Output: 5
}
