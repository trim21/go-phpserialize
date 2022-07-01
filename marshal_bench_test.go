package phpserialize_test

import (
	"fmt"
	"testing"

	elliotchance_phpserialize "github.com/elliotchance/phpserialize"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/internal/encoder"
)

func BenchmarkMarshal_all(b *testing.B) {
	var data = TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		B:   false,
		Obj: Inner{V: 2, S: "vvv"},

		Map: map[int]struct{ V int }{7: {V: 4}},
	}

	for i := 0; i < b.N; i++ {
		phpserialize.Marshal(data)
	}
}

func BenchmarkMarshal_map_concrete_types(b *testing.B) {
	for i := 1; i < 10000; i = i * 10 {
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			for i := 0; i < b.N; i++ {
				phpserialize.Marshal(m)
			}
		})
	}
}

func BenchmarkMarshal_slice_concrete_types(b *testing.B) {
	for i := 1; i < 10000; i = i * 10 {
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			for i := 0; i < b.N; i++ {
				encoder.Marshal(m)
			}
		})
	}
}

func Benchmark_elliotchance_phpserialize_marshal(b *testing.B) {
	var data = map[any]any{
		"users": []map[any]any{
			{"id": 1, "name": "sai"},
			{"id": 2, "name": "trim21"},
		},
		"obj": map[any]any{"v": 2, "s": "vvv"},
	}
	for i := 0; i < b.N; i++ {
		elliotchance_phpserialize.Marshal(data, nil)
	}
}
