package phpserialize_test

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"github.com/trim21/go-phpserialize"
)

func BenchmarkMarshal_type(b *testing.B) {
	for _, data := range testCase {
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := phpserialize.Marshal(data.Data)
				if err != nil {
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkMarshal_field_as_string(b *testing.B) {
	data := struct {
		F int `php:",string"`
	}{}
	for i := 0; i < b.N; i++ {
		_, err := phpserialize.Marshal(data)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkMarshal_ifce(b *testing.B) {
	for _, data := range testCase {
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := phpserialize.Marshal(data)
				if err != nil {
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkMarshal_map_type(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(m)
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_map_as_ifce(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				var v = struct{ Value any }{m}
				for pb.Next() {
					_, err := phpserialize.Marshal(v)
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_map_with_ifce_value(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]any, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(m)
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_slice_of_value(b *testing.B) {
	type D struct {
		Value []uint
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{m})
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_ifce_slice_as_value(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{m})
					if err != nil {
						b.FailNow()
					}
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_slice_of_type(b *testing.B) {
	type D struct {
		Value []User
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]User, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{m})
					if err != nil {
						b.FailNow()
					}
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_ifce_slice_of_type(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]User, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{Value: m})
					if err != nil {
						b.FailNow()
					}
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_ifce_slice_of_ifce(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]any, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{Value: m})
					if err != nil {
						b.FailNow()
					}
				}
			})
			runtime.KeepAlive(m)
		})
	}
}
