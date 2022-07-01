package main

import (
	"fmt"
	"runtime"

	"github.com/goccy/go-reflect"
	"github.com/gookit/goutil/dump"
	"github.com/trim21/go-phpserialize"
)

type Item struct {
	V int
}

type MapPtr struct {
	Users []Item `php:"users"`
	// Obj Inner        `php:"obj"`
	// B   bool         `php:"ok"`
	Map map[string]int64 `php:"map"`
}

func main() {
	var data = MapPtr{
		Map: map[string]int64{"key1": 1},
	}

	dump.P(reflect.ValueOf(data.Map))
	dump.P(reflect.ValueOf(data.Map["key1"]))

	b, err := phpserialize.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("len", len(b))
	fmt.Println(string(b))

	runtime.KeepAlive(data)
}
