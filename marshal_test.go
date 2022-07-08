package phpserialize_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

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

type ContainerAnonymous struct {
	OK   bool
	Item `json:"item" php:"item"`
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

type WithIgnore struct {
	V       int
	Ignored string `php:"-" json:"-"`
	D       any
}

var testCase = []struct {
	Name string
	Data any
}{
	{Name: "bool true", Data: true},
	{Name: "bool false", Data: false},
	{Name: "int8", Data: int8(7)},
	{Name: "int16", Data: int16(7)},
	{Name: "int32", Data: int32(7)},
	{Name: "int64", Data: int64(7)},
	{Name: "int", Data: int(8)},
	{Name: "uint8", Data: uint8(7)},
	{Name: "uint16", Data: uint16(7)},
	{Name: "uint32", Data: uint32(7)},
	{Name: "uint64", Data: uint64(7)},
	{Name: "uint", Data: uint(9)},
	{Name: "float32", Data: float32(3.14)},
	{Name: "float64", Data: float64(3.14)},
	{Name: "string", Data: strings.Repeat("qasd", 5)},
	{Name: "simple slice", Data: []int{1, 4, 6, 2, 3}},
	{Name: "struct slice", Data: []Item{{V: 6}, {V: 5}, {4}, {3}, {2}}},
	{Name: "struct with map ptr", Data: MapPtr{
		Users: []Item{},
		Map:   map[string]int64{"one": 1},
	}},
	{Name: "struct with map embed", Data: MapOnly{
		Map: map[string]int64{"one": 1},
	}},

	{Name: "nil map", Data: (map[string]string)(nil)},

	{Name: "nested struct not anonymous", Data: ContainerNonAnonymous{
		OK:   true,
		Item: Item{V: 5},
		V:    9999,
	}},

	{Name: "nested struct anonymous", Data: ContainerAnonymous{
		Item: Item{V: 5},
	}},

	{Name: "struct with all", Data: TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		B:   false,
		Obj: Inner{V: 2, S: "vvv"},

		Map: map[int]struct{ V int }{7: {V: 4}},
	}},

	{Name: "nested map", Data: NestedMap{
		1: map[uint]string{4: "ok"},
	}},

	{Name: "map[type]any(map)", Data: map[int]any{
		1: map[uint]string{4: "ok"},
	}},

	{Name: "map[type]any(slice)", Data: map[int]any{
		1: []int{3, 1, 4},
	}},

	{Name: "map[type]any(struct)", Data: map[int]any{
		1: User{},
	}},

	{Name: "generic[int]", Data: Generic[int]{1}},
	{Name: "generic[struct]", Data: Generic[User]{User{}}},
	{Name: "generic[map]", Data: Generic[map[string]int]{map[string]int{"one": 1}}},
	{Name: "generic[slice]", Data: Generic[[]string]{[]string{"hello", "world"}}},
}

func TestMarshal_concrete_types(t *testing.T) {
	t.Parallel()
	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			t.Parallel()
			b, err := phpserialize.Marshal(data.Data)
			require.NoError(t, err)

			j := decodeWithRealPhp(t, b)

			require.JSONEq(t, jsonEncode(t, data.Data), j)
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

			j := decodeWithRealPhp(t, b)

			require.JSONEq(t, jsonEncode(t, data), j, "lib: "+string(b)+"\nphp to json: "+j+"\njson.Marshal(data): "+jsonEncode(t, data))
		})
	}
}

// some special case like `empty map`, can't be compared by json unmarshal
func TestMarshal_special(t *testing.T) {
	t.Parallel()
	t.Run("empty map", func(t *testing.T) {
		t.Parallel()
		b, err := phpserialize.Marshal(map[int]string{})
		require.NoError(t, err)
		require.Equal(t, []byte("a:0:{}"), b, string(b))
	})
}

func TestMarshal_int_as_string(t *testing.T) {
	type Container struct {
		I int `php:"i,string"`
	}

	t.Run("negative", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: -104})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"-104";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 0})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:1:"0";}`
		require.Equal(t, expected, string(v))
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
		require.Equal(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{I: 1040})
		require.NoError(t, err)
		expected := `a:1:{s:1:"i";s:4:"1040";}`
		require.Equal(t, expected, string(v))
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
		require.Equal(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		require.Equal(t, expected, string(v))
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
		require.Equal(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: 1.00})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(Container{F: -3.14})
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		require.Equal(t, expected, string(v))
	})
}

func TestMarshal_float64_as_string_reflect(t *testing.T) {
	type Container struct {
		F float64 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		v, err := phpserialize.Marshal(WithIgnore{D: Container{F: 3.14}}.D)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:4:"3.14";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("positive", func(t *testing.T) {
		v, err := phpserialize.Marshal(WithIgnore{D: Container{F: 1.00}}.D)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:1:"1";}`
		require.Equal(t, expected, string(v))
	})

	t.Run("zero", func(t *testing.T) {
		v, err := phpserialize.Marshal(WithIgnore{D: Container{F: -3.14}}.D)
		require.NoError(t, err)
		expected := `a:1:{s:1:"f";s:5:"-3.14";}`
		require.Equal(t, expected, string(v))
	})
}

func decodeWithRealPhp(t *testing.T, s []byte) string {
	os.MkdirAll("./tmp/", 0700)
	fs := t.Name()
	for _, c := range "/{}()<>" {
		fs = strings.ReplaceAll(fs, string(c), "-")
	}

	fs = filepath.Join("./tmp/", fs) + ".php"

	file, err := os.Create(fs)
	require.NoError(t, err)

	fmt.Fprintf(file, `<?php

$val = unserialize('%s');

print json_encode($val);
`, s)

	require.NoError(t, file.Close())

	var buf = bytes.NewBuffer(nil)

	cmd := exec.Command("php", fs)
	cmd.Stdout = buf

	err = cmd.Run()
	require.NoError(t, err)

	return buf.String()
}

func jsonEncode(t *testing.T, value any) string {
	v, err := json.Marshal(value)
	require.NoError(t, err)

	return string(v)
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
		require.Equal(t, []byte(expected), actual)
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
		require.Equal(t, []byte(expected), actual)
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
		require.Equal(t, []byte(expected), actual)
	})

	t.Run("direct", func(t *testing.T) {
		// map in struct is a direct ptr
		type MapOnly struct {
			Map map[string]int64 `php:"map" json:"map"`
		}
		actual, err := phpserialize.Marshal(MapOnly{Map: map[string]int64{"abcdef": 1}})
		require.NoError(t, err)
		expected := `a:1:{s:3:"map";a:1:{s:6:"abcdef";i:1;}}`
		require.Equal(t, []byte(expected), actual)
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
		require.Equal(t, []byte(expected), actual)
	})
}
