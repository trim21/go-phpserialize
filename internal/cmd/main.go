package main

import (
	"fmt"

	"github.com/trim21/go-phpserialize"
)

func main() {
	var data = struct {
		Name string
		Data any
	}{Name: "int32", Data: int32(7)}

	b, err := phpserialize.Marshal(data)
	if err != nil {
		panic(err)

	}

	fmt.Println(string(b))
}

type MapPtr struct {
	Users []Item           `php:"users,omitempty" json:"users"`
	Map   map[string]int64 `php:"map" json:"map"`
}

type Item struct {
	V int `json:"v" php:"v"`
}
