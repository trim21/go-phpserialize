package main

import (
	"fmt"
	"runtime"

	"github.com/goccy/go-reflect"
	"github.com/gookit/goutil/dump"
	"github.com/trim21/go-phpserialize"
)

func main() {
	data := make(map[int]uint, 10)
	data[1] = 10

	dump.P(reflect.ValueOf(data))

	b, err := phpserialize.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("len", len(b))
	fmt.Println(string(b))

	runtime.KeepAlive(data)
}
