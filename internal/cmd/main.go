package main

import (
	"fmt"

	"github.com/trim21/go-phpserialize"
)

func main() {

	type Container struct {
		Value int `php:"value"`
	}

	var c Container
	raw := `a:1:{s:5:"value";i:10;}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	if err != nil {
		panic(err)
	}
	if c.Value != 10 {
		panic("not equal")
	}

	fmt.Println("correct")
}
