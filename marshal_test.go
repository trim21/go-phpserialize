package phpserialize_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v9"

	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/internal/runtime"
	"github.com/trim21/go-phpserialize/internal/test"
)

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
	{
		Name:     "*bool true",
		Data:     null.BoolFrom(true).Ptr(),
		Expected: "b:1;",
	},

	{
		Name: "bool true as string",
		Data: struct {
			B bool `php:"b,string"`
		}{B: true},
		Expected: `a:1:{s:1:"b";s:4:"true";}`,
	},

	{Name: "bool false", Data: false, Expected: "b:0;"},

	{
		Name: "bool true as string",
		Data: struct {
			B bool `php:"b,string"`
		}{},
		Expected: `a:1:{s:1:"b";s:5:"false";}`,
	},

	{Name: "*bool false", Data: new(bool), Expected: "b:0;"},

	{
		Name: "*bool-as-string-direct",
		Data: struct {
			Value *bool `php:"value,string"`
		}{Value: null.BoolFrom(false).Ptr()},
		Expected: `a:1:{s:5:"value";s:5:"false";}`,
	},

	{
		Name: "*bool-as-string-indirect",
		Data: struct {
			Value *bool `php:"value,string"`
			B     *bool
		}{Value: null.BoolFrom(false).Ptr()},
		Expected: `a:2:{s:5:"value";s:5:"false";s:1:"B";N;}`,
	},

	{Name: "int8", Data: int8(7), Expected: "i:7;"},
	{
		Name:     "*int8",
		Data:     null.Int8From(-7).Ptr(),
		Expected: "i:-7;",
	},

	{
		Name: "int8 as string",
		Data: struct {
			Value int8 `php:"value,string"`
		}{Value: 10},
		Expected: `a:1:{s:5:"value";s:2:"10";}`,
	},

	{Name: "int16", Data: int16(7), Expected: "i:7;"},
	{Name: "*int16", Data: null.Int16From(7).Ptr(), Expected: "i:7;"},
	{
		Name: "int16 as string",
		Data: struct {
			Value int16 `php:"value,string"`
		}{Value: 100},
		Expected: `a:1:{s:5:"value";s:3:"100";}`,
	},

	{Name: "int32", Data: int32(7), Expected: "i:7;"},
	{Name: "*int32", Data: null.Int32From(9).Ptr(), Expected: "i:9;"},
	{
		Name: "int32 as string",
		Data: struct {
			Value int32 `php:"value,string"`
		}{Value: 100},
		Expected: `a:1:{s:5:"value";s:3:"100";}`,
	},
	{
		Name: "*int32 as string",
		Data: struct {
			Value *int32 `php:"value,string"`
		}{Value: null.Int32From(100).Ptr()},
		Expected: `a:1:{s:5:"value";s:3:"100";}`,
	},

	{Name: "int64", Data: int64(7), Expected: "i:7;"},
	{Name: "*int64", Data: null.Int64From(10).Ptr(), Expected: "i:10;"},
	{Name: "int", Data: int(8), Expected: "i:8;"},
	{Name: "*int", Data: null.IntFrom(11).Ptr(), Expected: "i:11;"},
	{Name: "uint8", Data: uint8(7), Expected: "i:7;"},
	{Name: "*uint8", Data: null.Uint8From(7).Ptr(), Expected: "i:7;"},
	{Name: "uint16", Data: uint16(7), Expected: "i:7;"},
	{Name: "*uint16", Data: null.Uint16From(7).Ptr(), Expected: "i:7;"},
	{Name: "uint32", Data: uint32(7), Expected: "i:7;"},
	{Name: "*uint32", Data: null.Uint32From(7).Ptr(), Expected: "i:7;"},
	{Name: "uint64", Data: uint64(7777), Expected: "i:7777;"},
	{Name: "*uint64", Data: null.Uint64From(7).Ptr(), Expected: "i:7;"},
	{Name: "uint", Data: uint(9), Expected: "i:9;"},
	{Name: "*uint", Data: null.UintFrom(787).Ptr(), Expected: "i:787;"},
	{Name: "float32", Data: float32(3.14), Expected: "d:3.14;"},
	{Name: "*float32", Data: null.Float32From(3.14).Ptr(), Expected: "d:3.14;"},
	{Name: "float64", Data: float64(3.14), Expected: "d:3.14;"},
	{Name: "*float64", Data: null.Float64From(3.54).Ptr(), Expected: "d:3.54;"},
	{Name: "string", Data: `qwer"qwer`, Expected: `s:9:"qwer"qwer";`},
	{Name: "*string", Data: null.StringFrom(`qwer"qwer`).Ptr(), Expected: `s:9:"qwer"qwer";`},
	{Name: "simple slice", Data: []int{1, 4, 6, 2, 3}, Expected: `a:5:{i:0;i:1;i:1;i:4;i:2;i:6;i:3;i:2;i:4;i:3;}`},
	{
		Name:     "struct-slice",
		Data:     []Item{{V: 6}, {V: 5}, {4}, {3}, {2}},
		Expected: `a:5:{i:0;a:1:{s:1:"v";i:6;}i:1;a:1:{s:1:"v";i:5;}i:2;a:1:{s:1:"v";i:4;}i:3;a:1:{s:1:"v";i:3;}i:4;a:1:{s:1:"v";i:2;}}`,
	},
	{
		Name:     "struct-with-map-indirect",
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
		Name: "uint as string",
		Data: struct {
			Value uint `php:"value,string"`
		}{
			Value: uint(957123587),
		},
		Expected: `a:1:{s:5:"value";s:9:"957123587";}`,
	},
	{
		Name: "uint as string 0",
		Data: struct {
			Value uint `php:"value,string"`
		}{
			Value: uint(0),
		},
		Expected: `a:1:{s:5:"value";s:1:"0";}`,
	},

	{
		Name: "int as string positive",
		Data: struct {
			Value int `php:"value,string"`
		}{
			Value: int(957123587),
		},
		Expected: `a:1:{s:5:"value";s:9:"957123587";}`,
	},
	{
		Name: "int as string 0",
		Data: struct {
			Value int `php:"value,string"`
		}{
			Value: int(0),
		},
		Expected: `a:1:{s:5:"value";s:1:"0";}`,
	},

	{
		Name: "int as string positive",
		Data: struct {
			Value int `php:"value,string"`
		}{
			Value: int(957123587),
		},
		Expected: `a:1:{s:5:"value";s:9:"957123587";}`,
	},
	{
		Name: "int as string negative",
		Data: struct {
			Value int `php:"value,string"`
		}{
			Value: int(-957123587),
		},
		Expected: `a:1:{s:5:"value";s:10:"-957123587";}`,
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
		Name: "omitempty-ptr",
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

			test.StringEqual(t, data.Expected, string(actual))
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
			test.StringEqual(t, expected, string(actual))
		})
	}
}

func TestMarshal_interface_ptr(t *testing.T) {
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			actual, err := phpserialize.Marshal(&data.Data)
			require.NoError(t, err)

			test.StringEqual(t, data.Expected, string(actual))
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
		test.StringEqual(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("ptr-direct", func(t *testing.T) {
		data := struct {
			I *int `php:"i,string"`
		}{
			I: null.IntFrom(0).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("int indirect", func(t *testing.T) {
		data := struct {
			II *int `php:"ii,string,omitempty"`
			I  *int `php:"i,string"`
		}{
			I: null.IntFrom(0).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("int indirect", func(t *testing.T) {
		data := struct {
			II *int `php:"ii,string"`
			I  *int `php:"i,string"`
		}{
			I: null.IntFrom(0).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:2:{s:2:"ii";N;s:1:"i";s:1:"0";}`
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("positive", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{Value: S{F: 1.00}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("zero", func(t *testing.T) {
		actual, err := phpserialize.Marshal(Container{Value: S{F: -3.14}}.Value)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		test.StringEqual(t, expected, string(actual))
	})
}

func TestMarshal_ptr(t *testing.T) {
	t.Run("int-indirect-no-omit", func(t *testing.T) {
		type Indirect struct {
			A *int `php:"a"`
			B *int `php:"b"`
		}

		var i int = 50

		indirect := runtime.IfaceIndir(runtime.Type2RType(reflect.TypeOf(Indirect{})))
		require.True(t, indirect, "struct should be indirect")

		actual, err := phpserialize.Marshal(Indirect{B: &i})
		require.NoError(t, err)
		expected := `a:2:{s:1:"a";N;s:1:"b";i:50;}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("int-indirect-omitempty", func(t *testing.T) {
		type Indirect struct {
			A *int `php:"a"`
			B *int `php:"b,omitempty"`
		}

		var i int = 50

		indirect := runtime.IfaceIndir(runtime.Type2RType(reflect.TypeOf(Indirect{})))
		require.True(t, indirect, "struct should be indirect")

		actual, err := phpserialize.Marshal(Indirect{A: &i})
		require.NoError(t, err)
		expected := `a:1:{s:1:"a";i:50;}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("int-direct", func(t *testing.T) {
		type Direct struct {
			Value *int `php:"value"`
		}

		indirect := runtime.IfaceIndir(runtime.Type2RType(reflect.TypeOf(Direct{})))
		require.False(t, indirect, "struct should be indirect")

		var i int = 50

		t.Run("encode", func(t *testing.T) {
			actual, err := phpserialize.Marshal(Direct{Value: &i})
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";i:50;}`
			test.StringEqual(t, expected, string(actual))
		})
	})

	t.Run("nil", func(t *testing.T) {
		type Data struct {
			Value *int `php:"value"`
		}
		var data = Data{}

		actual, err := phpserialize.Marshal(data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";N;}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("*string", func(t *testing.T) {
		type Data struct {
			Value *string `php:"value"`
		}
		var s = "abcdefg"
		var data = Data{&s}

		actual, err := phpserialize.Marshal(data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";s:7:"abcdefg";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("multiple ptr", func(t *testing.T) {
		type Data struct {
			Value *string `php:"value"`
			D     *int    `php:"d,omitempty"`
		}
		var s = "abcdefg"
		var data = Data{Value: &s}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";s:7:"abcdefg";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("struct", func(t *testing.T) {
		t.Run("*struct", func(t *testing.T) {
			type Data struct {
				Value int    `php:"value"`
				ID    uint32 `php:"id"`
			}
			var data = Data{}

			actual, err := phpserialize.Marshal(&data)
			require.NoError(t, err)
			expected := `a:2:{s:5:"value";i:0;s:2:"id";i:0;}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("*struct-nil", func(t *testing.T) {
			type Data struct {
				Value int    `php:"value"`
				ID    uint32 `php:"id"`
			}
			var data *Data

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `N;`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("indirect", func(t *testing.T) {
			type Data struct {
				B     *int  `php:"b"`
				Value *User `php:"value"`
			}

			var b = 20
			var data = Data{B: &b}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:2:{s:1:"b";i:20;s:5:"value";N;}`
			test.StringEqual(t, expected, string(actual))
		})

		u := User{
			ID:   4,
			Name: "one",
		}

		t.Run("direct", func(t *testing.T) {
			type Data struct {
				Value *User `php:"value"`
			}
			var data = Data{}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";N;}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("encode direct", func(t *testing.T) {
			type Data struct {
				Value *User `php:"value"`
			}
			var data = Data{Value: &u}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:2:{s:2:"id";i:4;s:4:"name";s:3:"one";}}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("encode indirect", func(t *testing.T) {
			type Data struct {
				B     *int  `php:"b"`
				Value *User `php:"value"`
			}
			var data = Data{Value: &u}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:2:{s:1:"b";N;s:5:"value";a:2:{s:2:"id";i:4;s:4:"name";s:3:"one";}}`
			test.StringEqual(t, expected, string(actual))
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("nil-direct", func(t *testing.T) {
			type Data struct {
				Value *[5]int `php:"value"`
			}
			var data = Data{}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";N;}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("nil-indirect", func(t *testing.T) {
			type Data struct {
				Value *[5]int `php:"value"`
				B     *bool   `php:"b"`
			}

			var b = true
			var data = Data{B: &b}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:2:{s:5:"value";N;s:1:"b";b:1;}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("omitempty", func(t *testing.T) {
			type Data struct {
				Value *[5]int `php:"value,omitempty"`
			}
			var s = [5]int{1, 6, 4, 7, 9}
			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:5:{i:0;i:1;i:1;i:6;i:2;i:4;i:3;i:7;i:4;i:9;}}`
			test.StringEqual(t, expected, string(actual))
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
			test.StringEqual(t, expected, string(actual))
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
			test.StringEqual(t, expected, string(actual))
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
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("nil", func(t *testing.T) {
			type Data struct {
				Value *[]string `php:"value"`
			}

			var data = Data{}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";N;}`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("encode", func(t *testing.T) {
			type Data struct {
				Value *[]string `php:"value"`
			}

			var s = []string{"1", "2"}

			var data = Data{&s}

			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)
			expected := `a:1:{s:5:"value";a:2:{i:0;s:1:"1";i:1;s:1:"2";}}`
			test.StringEqual(t, expected, string(actual))
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
			test.StringEqual(t, expected, string(actual))
		})

	})

	t.Run("struct-map", func(t *testing.T) {
		t.Run("direct", func(t *testing.T) {
			type Data struct {
				Value *map[int]int `php:"value"`
			}

			t.Run("nil direct", func(t *testing.T) {
				var data = Data{}
				actual, err := phpserialize.Marshal(data)
				require.NoError(t, err)
				expected := `a:1:{s:5:"value";N;}`
				test.StringEqual(t, expected, string(actual))
			})

			t.Run("encode", func(t *testing.T) {
				var s = map[int]int{1: 2}

				actual, err := phpserialize.Marshal(&s)
				require.NoError(t, err)
				expected := `a:1:{i:1;i:2;}`
				test.StringEqual(t, expected, string(actual))
			})

			t.Run("omitempty encode", func(t *testing.T) {
				type Data struct {
					Value *map[int]int `php:"value,omitempty"`
				}

				var s = map[int]int{1: 2}
				var data = Data{&s}

				actual, err := phpserialize.Marshal(data)
				require.NoError(t, err)
				expected := `a:1:{s:5:"value";a:1:{i:1;i:2;}}`
				test.StringEqual(t, expected, string(actual))
			})

			t.Run("omitempty nil", func(t *testing.T) {
				type Data struct {
					Value *map[int]int `php:"value,omitempty"`
				}
				var data = Data{}

				actual, err := phpserialize.Marshal(data)
				require.NoError(t, err)
				expected := `a:0:{}`
				test.StringEqual(t, expected, string(actual))
			})
		})

		t.Run("indirect", func(t *testing.T) {
			type Data struct {
				Value *map[int]int `php:"value"`
				Bool  *bool        `php:"b"`
			}

			t.Run("nil direct", func(t *testing.T) {
				var data = Data{}
				actual, err := phpserialize.Marshal(data)
				require.NoError(t, err)
				expected := `a:2:{s:5:"value";N;s:1:"b";N;}`
				test.StringEqual(t, expected, string(actual))
			})

			t.Run("encode", func(t *testing.T) {
				var s = map[int]int{1: 2}

				actual, err := phpserialize.Marshal(&s)
				require.NoError(t, err)
				expected := `a:1:{i:1;i:2;}`
				test.StringEqual(t, expected, string(actual))
			})

			t.Run("omitempty", func(t *testing.T) {
				type Data struct {
					Value *map[int]int `php:"value,omitempty"`
					Bool  *bool        `php:"b"`
				}

				var s = map[int]int{1: 2}
				var data = Data{Value: &s}

				actual, err := phpserialize.Marshal(data)
				require.NoError(t, err)
				expected := `a:1:{s:5:"value";a:1:{i:1;i:2;}}`
				test.StringEqual(t, expected, string(actual))
			})
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			type Data struct {
				Value *map[int]int `php:"value"`
			}

			var data = Data{}

			actual, err := phpserialize.Marshal(data.Value)
			require.NoError(t, err)
			expected := `N;`
			test.StringEqual(t, expected, string(actual))
		})

		t.Run("encode", func(t *testing.T) {
			var s = map[int]int{1: 2}

			actual, err := phpserialize.Marshal(&s)
			require.NoError(t, err)
			expected := `a:1:{i:1;i:2;}`
			test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("nested", func(t *testing.T) {
		type Container struct {
			Value ***uint `php:"value"`
		}

		var v uint = 8
		var p = &v
		var a = &p

		_, err := phpserialize.Marshal(Container{Value: &a})
		require.Error(t, err)
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
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("direct", func(t *testing.T) {
		// map in struct is a direct ptr
		type MapOnly struct {
			Map map[string]int64 `php:"map" json:"map"`
		}
		actual, err := phpserialize.Marshal(MapOnly{Map: map[string]int64{"abcdef": 1}})
		require.NoError(t, err)
		expected := `a:1:{s:3:"map";a:1:{s:6:"abcdef";i:1;}}`
		test.StringEqual(t, expected, string(actual))
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
		test.StringEqual(t, expected, string(actual))
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
	test.StringEqual(t, expected, string(actual))
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

	test.StringEqual(t, `a:3:{s:1:"A";i:3;s:1:"B";i:2;s:1:"C";i:1;}`, string(actual))
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

	test.StringEqual(t, `a:3:{s:1:"A";i:3;s:1:"B";i:2;s:1:"C";i:1;}`, string(actual))
}

func TestRecursivePanic(t *testing.T) {
	t.Parallel()

	type O struct {
		Name string
		E    []O
	}

	actual, err := phpserialize.Marshal(O{
		Name: "hello",
		E: []O{
			{
				Name: "BB",
				E: []O{
					{Name: "C C D D E E F F"},
				},
			},
		},
	})
	require.NoError(t, err)

	test.StringEqual(t, `a:2:{s:4:"Name";s:5:"hello";s:1:"E";a:1:{i:0;a:2:{s:4:"Name";s:2:"BB";s:1:"E";a:1:{i:0;a:2:{s:4:"Name";s:15:"C C D D E E F F";s:1:"E";N;}}}}}`, string(actual))
}
