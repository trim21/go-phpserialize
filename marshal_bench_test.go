package phpserialize_test

import (
	"fmt"
	"runtime"
	"strings"
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
			{ID: 3, Name: "g"},
		},
		// B:   false,
		// Obj: Inner{V: 2, S: "vvv"},

		// Map: map[int]struct{ V int }{7: {V: 4}},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encoder.Marshal(data)
		}
	})
}

func BenchmarkAll_concrete_types(b *testing.B) {
	for _, data := range testCase {
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				phpserialize.Marshal(data.Data)
			}
		})
	}
}

func BenchmarkAll_interface(b *testing.B) {
	for _, data := range testCase {
		if strings.Contains(data.Name, "map") {
			continue
		}
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				phpserialize.Marshal(data)
			}
		})
	}
}

func BenchmarkMarshal_map_concrete_types(b *testing.B) {
	for i := 1; i < 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(m)
				}
			})
		})
	}
}

func BenchmarkMarshal_slice_concrete_types(b *testing.B) {
	for i := 1; i < 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(m)
				}
			})
			runtime.KeepAlive(m)
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
