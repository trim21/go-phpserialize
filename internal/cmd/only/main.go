package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/trim21/go-phpserialize"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		marshal1()
	}()

	go func() {
		defer wg.Done()
		marshal2()
	}()

	go func() {
		defer wg.Done()
		marshal3()
	}()

	wg.Wait()
}

const N = 100000000

func marshal1() {
	data := make(map[int]uint, 10)
	data[1] = 10

	var b []byte
	var err error

	for i := 0; i < N; i++ {
		b, err = phpserialize.Marshal(data)
		if err != nil {
			panic(err)
		}
		runtime.GC()
	}

	fmt.Println("len", len(b))
	fmt.Println(string(b))

	runtime.KeepAlive(data)
}

func marshal2() {
	i := 100
	var data = make(map[int]uint, i)
	for j := 0; j < i; j++ {
		data[j+1] = uint(j + 2)
	}

	var b []byte
	var err error

	for i := 0; i < N; i++ {
		b, err = phpserialize.Marshal(data)
		if err != nil {
			panic(err)
		}
		runtime.GC()

	}

	fmt.Println("len", len(b))
	fmt.Println(string(b))

	runtime.KeepAlive(data)
}

func marshal3() {
	var data = TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
			{ID: 3, Name: "g"},
		},
		B:   false,
		Obj: Inner{V: 2, S: "vvv"},

		Map: map[int]struct{ V int }{7: {V: 4}},
	}

	var b []byte
	var err error

	for i := 0; i < N; i++ {
		b, err = phpserialize.Marshal(data)
		if err != nil {
			panic(err)
		}
		runtime.GC()

	}

	fmt.Println("len", len(b))
	fmt.Println(string(b))

	runtime.KeepAlive(data)
}

type Inner struct {
	V int    `php:"v"`
	S string `php:"a long string name replace field name"`
}

type TestData struct {
	Users []User                  `php:"users"`
	Obj   Inner                   `php:"obj"`
	B     bool                    `php:"ok"`
	Map   map[int]struct{ V int } `php:"map"`
}

type User struct {
	ID   uint64 `php:"id"`
	Name string `php:"name"`
}
