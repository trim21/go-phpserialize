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

			phpserialize.Marshal(D{Value: m}) // warm up

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{Value: m})
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}
func BenchmarkMarshal_1_type_slice_of_type(b *testing.B) {
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
					_, err := phpserialize.Marshal(D{Value: m})
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_1_ifce_slice_of_ifce(b *testing.B) {
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

			phpserialize.Marshal(D{Value: m}) // warm up

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := phpserialize.Marshal(D{Value: m})
					if err != nil {
						b.FailNow()
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_large_struct_10(b *testing.B) {
	var v struct {
		Field0 bool
		Field1 bool
		Field2 bool
		Field3 bool
		Field4 bool
		Field5 bool
		Field6 bool
		Field7 bool
		Field8 bool
		Field9 bool
	}

	for i := 0; i < b.N; i++ {
		_, err := phpserialize.Marshal(v)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkMarshal_large_struct_50(b *testing.B) {
	var v struct {
		Field0  bool
		Field1  bool
		Field2  bool
		Field3  bool
		Field4  bool
		Field5  bool
		Field6  bool
		Field7  bool
		Field8  bool
		Field9  bool
		Field10 bool
		Field11 bool
		Field12 bool
		Field13 bool
		Field14 bool
		Field15 bool
		Field16 bool
		Field17 bool
		Field18 bool
		Field19 bool
		Field20 bool
		Field21 bool
		Field22 bool
		Field23 bool
		Field24 bool
		Field25 bool
		Field26 bool
		Field27 bool
		Field28 bool
		Field29 bool
		Field30 bool
		Field31 bool
		Field32 bool
		Field33 bool
		Field34 bool
		Field35 bool
		Field36 bool
		Field37 bool
		Field38 bool
		Field39 bool
		Field40 bool
		Field41 bool
		Field42 bool
		Field43 bool
		Field44 bool
		Field45 bool
		Field46 bool
		Field47 bool
		Field48 bool
		Field49 bool
	}
	for i := 0; i < b.N; i++ {
		_, err := phpserialize.Marshal(v)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkMarshal_large_struct_100(b *testing.B) {
	var v struct {
		Field0  bool
		Field1  bool
		Field2  bool
		Field3  bool
		Field4  bool
		Field5  bool
		Field6  bool
		Field7  bool
		Field8  bool
		Field9  bool
		Field10 bool
		Field11 bool
		Field12 bool
		Field13 bool
		Field14 bool
		Field15 bool
		Field16 bool
		Field17 bool
		Field18 bool
		Field19 bool
		Field20 bool
		Field21 bool
		Field22 bool
		Field23 bool
		Field24 bool
		Field25 bool
		Field26 bool
		Field27 bool
		Field28 bool
		Field29 bool
		Field30 bool
		Field31 bool
		Field32 bool
		Field33 bool
		Field34 bool
		Field35 bool
		Field36 bool
		Field37 bool
		Field38 bool
		Field39 bool
		Field40 bool
		Field41 bool
		Field42 bool
		Field43 bool
		Field44 bool
		Field45 bool
		Field46 bool
		Field47 bool
		Field48 bool
		Field49 bool
		Field50 bool
		Field51 bool
		Field52 bool
		Field53 bool
		Field54 bool
		Field55 bool
		Field56 bool
		Field57 bool
		Field58 bool
		Field59 bool
		Field60 bool
		Field61 bool
		Field62 bool
		Field63 bool
		Field64 bool
		Field65 bool
		Field66 bool
		Field67 bool
		Field68 bool
		Field69 bool
		Field70 bool
		Field71 bool
		Field72 bool
		Field73 bool
		Field74 bool
		Field75 bool
		Field76 bool
		Field77 bool
		Field78 bool
		Field79 bool
		Field80 bool
		Field81 bool
		Field82 bool
		Field83 bool
		Field84 bool
		Field85 bool
		Field86 bool
		Field87 bool
		Field88 bool
		Field89 bool
		Field90 bool
		Field91 bool
		Field92 bool
		Field93 bool
		Field94 bool
		Field95 bool
		Field96 bool
		Field97 bool
		Field98 bool
		Field99 bool
	}

	for i := 0; i < b.N; i++ {
		_, err := phpserialize.Marshal(v)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkMarshal_int(b *testing.B) {
	for i := 10; i <= 1000; i *= 10 {
		i := i
		b.Run(fmt.Sprintf("marshal int %d", i), func(b *testing.B) {
			var s = make([]int, i)
			for j := 0; j < i; j++ {
				s[j] = j
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, err := phpserialize.Marshal(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_uint(b *testing.B) {
	for i := 10; i <= 1000; i *= 10 {
		var s = make([]uint, i)
		for j := 0; j < i; j++ {
			s[j] = uint(j)
		}
		b.Run(fmt.Sprintf("marshal uint %d", i), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, err := phpserialize.Marshal(s); err != nil {
						b.Fatal(err.Error())
					}
				}
			})
		})
	}
}

func BenchmarkMarshal_many_map_field(b *testing.B) {
	var s = struct {
		Map1 map[string]int `php:"map_1"`
		Map2 map[int]string `php:"map_2"`
		Map3 map[uint]int   `php:"map_3"`
	}{
		Map1: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		Map2: map[int]string{
			1: "1111",
			2: "2222",
			3: "3333",
			4: "4444",
			5: "5555",
			6: "6666",
		},
		Map3: map[uint]int{
			1: 1111,
			2: 2222,
			3: 3333,
			4: 4444,
			5: 5555,
			6: 6666,
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := phpserialize.Marshal(s)
			if err != nil {
				b.FailNow()
			}
		}
	})
}
