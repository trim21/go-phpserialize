package main

import (
	"fmt"
	"strings"

	"github.com/trim21/go-phpserialize"
)

func main() {

	for i := 0; i < 100; i++ {
		fmt.Printf("Field%d int\n", i)
	}
	var data = struct {
		Name string
		Data any
	}{
		Name: "name field value",
		Data: strings.Repeat("qasd", 5),
	}
	// var m = User{ID: uint64(5 + 2), Name: "ue"}

	b, err := phpserialize.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

type User struct {
	ID   uint64 `php:"id" json:"id"`
	Name string `php:"name" json:"name"`
}
