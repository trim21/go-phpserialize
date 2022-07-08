package bench

import (
	"testing"

	elliotchance_phpserialize "github.com/elliotchance/phpserialize"
	"github.com/trim21/go-phpserialize"
)

func Benchmark_marshal_compare(b *testing.B) {
	type Obj struct {
		V int `php:"v"`
		S int `php:"s"`
	}

	type User struct {
		ID   uint64 `php:"id" json:"id"`
		Name string `php:"name" json:"name"`
	}

	type TestData struct {
		Users []User `php:"users"`
		Obj   Obj    `php:"obj"`
	}

	var data = TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		Obj: Obj{V: 2, S: 3},
	}

	for i := 0; i < b.N; i++ {
		phpserialize.Marshal(data)
	}
}

func Benchmark_elliotchance_phpserialize_marshal(b *testing.B) {
	var data = map[any]any{
		"users": []map[any]any{
			{"id": 1, "name": "sai"},
			{"id": 2, "name": "trim21"},
		},
		"obj": map[string]int{"v": 2, "s": 3},
	}
	for i := 0; i < b.N; i++ {
		elliotchance_phpserialize.Marshal(data, nil)
	}
}
