package phpserialize_test

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

func init() {
	color.NoColor = false // force color
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

type Generic[T any] struct {
	Value T
}

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
		Name:     "generic[int]",
		Data:     Generic[int]{1},
		Expected: `a:1:{s:5:"Value";i:1;}`,
	},

	{
		Name:     "generic[struct]",
		Data:     Generic[User]{User{}},
		Expected: `a:1:{s:5:"Value";a:2:{s:2:"id";i:0;s:4:"name";s:0:"";}}`,
	},

	{
		Name:     "generic[map]",
		Data:     Generic[map[string]int]{map[string]int{"one": 1}},
		Expected: `a:1:{s:5:"Value";a:1:{s:3:"one";i:1;}}`,
	},

	{
		Name:     "generic[slice]",
		Data:     Generic[[]string]{[]string{"hello", "world"}},
		Expected: `a:1:{s:5:"Value";a:2:{i:0;s:5:"hello";i:1;s:5:"world";}}`,
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
}

func TestMarshal_concrete_types(t *testing.T) {
	t.Parallel()
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			b, err := phpserialize.Marshal(data.Data)
			require.NoError(t, err)

			stringEqual(t, data.Expected, string(b))
		})
	}
}

func TestMarshal_interface(t *testing.T) {
	t.Parallel()
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			t.Parallel()
			b, err := phpserialize.Marshal(data)
			require.NoError(t, err)

			expected := fmt.Sprintf(`a:2:{s:4:"Name";s:%d:"%s";s:4:"Data";`, len(data.Name), data.Name) + data.Expected + "}"
			actual := string(b)
			stringEqual(t, expected, actual)
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
		v, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		stringEqual(t, expected, string(v))
	})
}

func TestMarshal_uint_as_string(t *testing.T) {
	type Container struct {
		I uint `php:"i,string"`
	}

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		stringEqual(t, expected, string(v))
	})
}

func TestMarshal_bool_as_string(t *testing.T) {
	type Container struct {
		B bool `php:",string"`
	}

	t.Run("true", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{true})
		require.NoError(t, err)
		require.Equal(t, `a:1:{s:1:"B";s:4:"true";}`, string(v))
	})

	t.Run("false", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{false})
		require.NoError(t, err)
		require.Equal(t, `a:1:{s:1:"B";s:5:"false";}`, string(v))
	})
}

func TestMarshal_float32_as_string(t *testing.T) {
	type Container struct {
		F float32 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(v))
	})
}

func TestMarshal_float64_as_string(t *testing.T) {
	type Container struct {
		F float64 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(v))
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
		v, err := phpserialize.Marshal(Container{Value: S{F: 3.14}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{Value: S{F: 1.00}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		stringEqual(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{Value: S{F: -3.14}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		stringEqual(t, expected, string(v))
	})
}

func TestMarshal_ptr(t *testing.T) {
	t.Parallel()

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

	t.Run("*int", func(t *testing.T) {
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

		expected := `a:1:{s:3:"f1q";s:10:"0147852369";}`
		actual, err := phpserialize.Marshal(Container{Value: &a})
		require.NoError(t, err)
		stringEqual(t, expected, string(actual))
	})
}

func TestMarshal_map(t *testing.T) {
	t.Parallel()

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

func stringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("Result not as expected:\n%v", CharacterDiff(expected, actual))
		t.FailNow()
	}
}
