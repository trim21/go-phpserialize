package phpserialize_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

type any = interface{}

func init() {
	color.NoColor = false // force color
}

type Container struct {
	Value any `php:"value"`
}

type Inner struct {
	V int    `php:"v" json:"v"`
	S string `php:"a long string name replace field name" json:"a long string name replace field name"`
}

type TestData struct {
	Users []User                  `php:"users" json:"users"`
	Obj   Inner                   `php:"obj" json:"obj"`
	B     bool                    `php:"ok" json:"ok"`
	Map   map[int]struct{ V int } `php:"map" json:"map"`
}

type User struct {
	ID   uint64 `php:"id" json:"id"`
	Name string `php:"name" json:"name"`
}

type Item struct {
	V int `json:"v" php:"v"`
}

type ContainerNonAnonymous struct {
	OK   bool
	Item Item
	V    int
}

// map in struct is an indirect ptr
type MapPtr struct {
	Users []Item           `php:"users" json:"users"`
	Map   map[string]int64 `php:"map" json:"map"`
}

// map in struct is a direct ptr
type MapOnly struct {
	Map map[string]int64 `php:"map" json:"map"`
}

type NestedMap = map[int]map[uint]string

var testCase = []struct {
	Name     string
	Data     any
	Expected string `php:"-" json:"-"`
}{
	{Name: "bool true", Data: true, Expected: "b:1;"},
	{Name: "bool false", Data: false, Expected: "b:0;"},
	{Name: "int8", Data: int8(7), Expected: "i:7;"},
	{Name: "int16", Data: int16(7), Expected: "i:7;"},
	{Name: "int32", Data: int32(7), Expected: "i:7;"},
	{Name: "int64", Data: int64(7), Expected: "i:7;"},
	{Name: "int", Data: int(8), Expected: "i:8;"},
	{Name: "uint8", Data: uint8(7), Expected: "i:7;"},
	{Name: "uint16", Data: uint16(7), Expected: "i:7;"},
	{Name: "uint32", Data: uint32(7), Expected: "i:7;"},
	{Name: "uint64", Data: uint64(7777), Expected: "i:7777;"},
	{Name: "uint", Data: uint(9), Expected: "i:9;"},
	{Name: "float32", Data: float32(3.14), Expected: "d:3.14;"},
	{Name: "float64", Data: float64(3.14), Expected: "d:3.14;"},
	{Name: "string", Data: `qwer"qwer`, Expected: `s:9:"qwer"qwer";`},
	{Name: "simple slice", Data: []int{1, 4, 6, 2, 3}, Expected: `a:5:{i:0;i:1;i:1;i:4;i:2;i:6;i:3;i:2;i:4;i:3;}`},
	{
		Name:     "struct slice",
		Data:     []Item{{V: 6}, {V: 5}, {4}, {3}, {2}},
		Expected: `a:5:{i:0;a:1:{s:1:"v";i:6;}i:1;a:1:{s:1:"v";i:5;}i:2;a:1:{s:1:"v";i:4;}i:3;a:1:{s:1:"v";i:3;}i:4;a:1:{s:1:"v";i:2;}}`,
	},
	{
		Name:     "struct with map ptr",
		Data:     MapPtr{Users: []Item{}, Map: map[string]int64{"one": 1}},
		Expected: `a:2:{s:5:"users";a:0:{}s:3:"map";a:1:{s:3:"one";i:1;}}`,
	},
	{
		Name:     "struct with map embed",
		Data:     MapOnly{Map: map[string]int64{"one": 1}},
		Expected: `a:1:{s:3:"map";a:1:{s:3:"one";i:1;}}`,
	},
	{
		Name:     "empty map",
		Data:     map[int]string{},
		Expected: "a:0:{}",
	},
	{
		Name:     "nil map",
		Data:     (map[string]string)(nil),
		Expected: `N;`,
	},

	{
		Name: "nested struct not anonymous",
		Data: ContainerNonAnonymous{
			OK:   true,
			Item: Item{V: 5},
			V:    9999,
		},
		Expected: `a:3:{s:2:"OK";b:1;s:4:"Item";a:1:{s:1:"v";i:5;}s:1:"V";i:9999;}`,
	},

	// {
	// 	Name:     "nested struct anonymous",
	// 	Data:     ContainerAnonymous{Item: Item{V: 5}},
	// 	Expected: `a:2:{s:2:"OK";b:0;s:4:"item";a:1:{s:1:"v";i:5;}}`,
	// },

	{
		Name: "struct with all",
		Data: TestData{
			Users: []User{
				{ID: 1, Name: "sai"},
				{ID: 2, Name: "trim21"},
			},
			B:   false,
			Obj: Inner{V: 2, S: "vvv"},

			Map: map[int]struct{ V int }{7: {V: 4}},
		},
		Expected: `a:4:{s:5:"users";a:2:{i:0;a:2:{s:2:"id";i:1;s:4:"name";s:3:"sai";}i:1;a:2:{s:2:"id";i:2;s:4:"name";s:6:"trim21";}}s:3:"obj";a:2:{s:1:"v";i:2;s:37:"a long string name replace field name";s:3:"vvv";}s:2:"ok";b:0;s:3:"map";a:1:{i:7;a:1:{s:1:"V";i:4;}}}`,
	},

	{
		Name:     "nested map",
		Data:     NestedMap{1: map[uint]string{4: "ok"}},
		Expected: `a:1:{i:1;a:1:{i:4;s:2:"ok";}}`,
	},

	{
		Name:     "map[type]any(map)",
		Data:     map[int]any{1: map[uint]string{4: "ok"}},
		Expected: `a:1:{i:1;a:1:{i:4;s:2:"ok";}}`,
	},

	{
		Name:     "map[type]any(slice)",
		Data:     map[int]any{1: []int{3, 1, 4}},
		Expected: `a:1:{i:1;a:3:{i:0;i:3;i:1;i:1;i:2;i:4;}}`,
	},

	{
		Name:     "map[type]any(struct)",
		Data:     map[int]any{1: User{}},
		Expected: `a:1:{i:1;a:2:{s:2:"id";i:0;s:4:"name";s:0:"";}}`,
	},

	{
		Name: "ignore struct field",
		Data: struct {
			V       int
			Ignored string `php:"-"`
		}{
			V:       3,
			Ignored: "vvv",
		},
		Expected: `a:1:{s:1:"V";i:3;}`,
	},
	{
		Name: "private field",
		Data: struct {
			b bool
			D int
		}{D: 10},
		Expected: `a:1:{s:1:"D";i:10;}`,
	},
	{
		Name: "omitempty",
		Data: struct {
			V string `php:",omitempty"`
			D string `php:",omitempty"`
		}{D: "d"},
		Expected: `a:1:{s:1:"D";s:1:"d";}`,
	},
	{
		Name: "omitempty ptr",
		Data: struct {
			V *string `php:",omitempty"`
			D *string `php:",omitempty"`
		}{
			D: new(string),
		},
		Expected: `a:1:{s:1:"D";s:0:"";}`,
	},
}

func TestMarshal_concrete_types(t *testing.T) {
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			actual, err := phpserialize.Marshal(data.Data)
			require.NoError(t, err)

			stringEqual(t, data.Expected, string(actual))
		})
	}
}

func TestMarshal_interface(t *testing.T) {
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)

			expected := fmt.Sprintf(`a:2:{s:4:"Name";s:%d:"%s";s:4:"Data";`, len(data.Name), data.Name) + data.Expected + "}"
			stringEqual(t, expected, string(actual))
		})
	}
}

func TestMarshal_interface_ptr(t *testing.T) {
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			actual, err := phpserialize.Marshal(&data.Data)
			require.NoError(t, err)

			stringEqual(t, data.Expected, string(actual))
		})
	}
}

func TestMarshal_int_as_string(t *testing.T) {
	type Container struct {
		I int `php:"i,string"`
	}

	t.Run("negative", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: -104})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"-104";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_uint_as_string(t *testing.T) {
	type Container struct {
		I uint `php:"i,string"`
	}

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_bool_as_string(t *testing.T) {
	type Container struct {
		B bool `php:",string"`
	}

	t.Run("true", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{true})
		require.NoError(t, err)
		require.Equal(t, `a:1:{s:1:"B";s:4:"true";}`, string(actual))
	})

	t.Run("false", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{false})
		require.NoError(t, err)
		require.Equal(t, `a:1:{s:1:"B";s:5:"false";}`, string(actual))
	})
}

func TestMarshal_float32_as_string(t *testing.T) {
	type Container struct {
		F float32 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_float64_as_string(t *testing.T) {
	type Container struct {
		F float64 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_float64_as_string_reflect(t *testing.T) {
	type Container struct {
		Value any `php:"value"`
	}
	type S struct {
		F float64 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{Value: S{F: 3.14}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{Value: S{F: 1.00}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{Value: S{F: -3.14}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_ptr(t *testing.T) {
	t.Run("*string", func(t *testing.T) {
		type Data struct {
			Value *string `php:"value"`
		}
		var s = "abcdefg"
		var data = Data{&s}

		actual, err := phpserialize.Marshal(data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";s:7:"abcdefg";}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("array", func(t *testing.T) {

		t.Run("omitempty", func(t *testing.T) {
			type Data struct {
				Value *[5]int `php:"value,omitempty"`
			}
			var s = [5]int{1, 6, 4, 7, 9}
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:5:{i:0;i:1;i:1;i:6;i:2;i:4;i:3;i:7;i:4;i:9;}}`
			stringEqual(t, expected, string(actual))
		})

		t.Run("no omitempty", func(t *testing.T) {
			type Data struct {
				Value *[5]int `php:"value"`
			}
			var s = [5]int{1, 6, 4, 7, 9}
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:5:{i:0;i:1;i:1;i:6;i:2;i:4;i:3;i:7;i:4;i:9;}}`
			stringEqual(t, expected, string(actual))
		})
	})

	t.Run("slice", func(t *testing.T) {
		t.Run("omitempty", func(t *testing.T) {
			type Data struct {
				Value *[]string `php:"value,omitempty"`
			}
			var s = strings.Split("abcdefg", "")
			require.Len(t, s, 7)
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:7:{i:0;s:1:"a";i:1;s:1:"b";i:2;s:1:"c";i:3;s:1:"d";i:4;s:1:"e";i:5;s:1:"f";i:6;s:1:"g";}}`
			stringEqual(t, expected, string(actual))
		})

		t.Run("no omitempty", func(t *testing.T) {
			type Data struct {
				Value *[]string `php:"value"`
			}
			var s = strings.Split("abcdefg", "")
			require.Len(t, s, 7)
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:7:{i:0;s:1:"a";i:1;s:1:"b";i:2;s:1:"c";i:3;s:1:"d";i:4;s:1:"e";i:5;s:1:"f";i:6;s:1:"g";}}`
			stringEqual(t, expected, string(actual))
		})
	})

	t.Run("*string omitempty", func(t *testing.T) {
		type Data struct {
			Value *string `php:"value,omitempty"`
		}

		t.Run("not_nil", func(t *testing.T) {
			var s = "abcdefg"
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";s:7:"abcdefg";}`
			stringEqual(t, expected, string(actual))
		})

	})

	t.Run("map", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			type Data struct {
				Value *map[int]int `php:"value"`
			}

			var data = Data{}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";N;}`
			stringEqual(t, expected, string(actual))
		})

		t.Run("omitempty", func(t *testing.T) {
			type Data struct {
				Value *map[int]int `php:"value,omitempty"`
			}

			var s = map[int]int{1: 2}
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:1:{i:1;i:2;}}`
			stringEqual(t, expected, string(actual))
		})
	})

	t.Run("int", func(t *testing.T) {
		type Data struct {
			Value *int `php:"value"`
		}
		var s = 644
		var data = Data{&s}

		actual, err := phpserialize.Marshal(data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";i:644;}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("recursive", func(t *testing.T) {
		type Container struct {
			Value any `php:"value"`
		}

		var v uint = 8
		var p = &v
		var a any = &p

		expected := `a:1:{s:5:"value";i:8;}`
		actual, err := phpserialize.Marshal(Container{Value: &a})
		require.NoError(t, err)
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_map(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		// map in struct is a direct ptr
		type MapOnly struct {
			Map map[string]int64 `php:"map" json:"map"`
		}
		actual, err := phpserialize.Marshal(MapOnly{Map: nil})
		require.NoError(t, err)
		expected := `a:1:{s:3:"map";N;}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("direct", func(t *testing.T) {
		// map in struct is a direct ptr
		type MapOnly struct {
			Map map[string]int64 `php:"map" json:"map"`
		}
		actual, err := phpserialize.Marshal(MapOnly{Map: map[string]int64{"abcdef": 1}})
		require.NoError(t, err)
		expected := `a:1:{s:3:"map";a:1:{s:6:"abcdef";i:1;}}`
		stringEqual(t, expected, string(actual))
	})

	t.Run("indirect", func(t *testing.T) {
		// map in struct is an indirect ptr
		type MapPtr struct {
			Users []Item           `php:"users" json:"users"`
			Map   map[string]int64 `php:"map" json:"map"`
		}

		actual, err := phpserialize.Marshal(MapPtr{Map: map[string]int64{"abcdef": 1}})
		require.NoError(t, err)
		expected := `a:2:{s:5:"users";N;s:3:"map";a:1:{s:6:"abcdef";i:1;}}`
		stringEqual(t, expected, string(actual))
	})
}

type M interface {
	Bool() bool
}

type mImpl struct {
}

func (m mImpl) Bool() bool {
	return true
}

func TestMarshal_interface_with_method(t *testing.T) {
	var data M = mImpl{}
	actual, err := phpserialize.Marshal(Container{Value: data})
	require.NoError(t, err)
	expected := `a:1:{s:5:"value";a:0:{}}`
	stringEqual(t, expected, string(actual))
}

func stringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("Result not as expected:\n%v", CharacterDiff(expected, actual))
		t.FailNow()
	}
}

func TestMarshal_anonymous_field(t *testing.T) {
	type N struct {
		A int
		B int
	}

	type M struct {
		N
		C int
	}

	actual, err := phpserialize.Marshal(M{N: N{
		A: 3,
		B: 2,
	}, C: 1})
	require.NoError(t, err)

	stringEqual(t, `a:3:{s:1:"A";i:3;s:1:"B";i:2;s:1:"C";i:1;}`, string(actual))
}

func TestMarshal_anonymous_field_omitempty(t *testing.T) {
	type L struct {
		E int `php:"E,omitempty"`
	}

	type N struct {
		L
		A int
		B int
	}

	type M struct {
		N
		C int
	}

	actual, err := phpserialize.Marshal(M{N: N{
		A: 3,
		B: 2,
	}, C: 1})
	require.NoError(t, err)

	stringEqual(t, `a:3:{s:1:"A";i:3;s:1:"B";i:2;s:1:"C";i:1;}`, string(actual))
}
