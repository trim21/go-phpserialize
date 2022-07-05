package main

import (
	"fmt"

	"github.com/trim21/go-phpserialize"
)

func main() {

	type Container struct {
		F float64 `php:"f1q"`
	}

	var c Container
	raw := `a:1:{s:3:"f1q";d:147852369;}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	if err != nil {
		panic(err)
	}
	if c.F != 147852369 {
		panic("not equal")
	}

	fmt.Println("correct")
}
