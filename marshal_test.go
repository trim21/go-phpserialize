package phpserialize_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-json"
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

func TestMarshal(t *testing.T) {
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
		{Name: "slice", Data: []Item{{V: 6}, {V: 5}, {4}, {3}, {2}}},
		{Name: "struct with map ptr", Data: MapPtr{
			Users: []Item{},
			Map:   map[string]int64{"one": 1},
		}},
		{Name: "struct with map embed", Data: MapOnly{
			Map: map[string]int64{"one": 1},
		}},

		{Name: "nil map", Data: MapOnly{}},

		{Name: "nested struct not anonymous", Data: ContainerNonAnonymous{
			OK:   true,
			Item: Item{V: 5},
			V:    9999,
		}},

		{Name: "nested struct anonymous", Data: ContainerAnonymous{
			Item: Item{V: 5},
		}},

		{Name: "complex object", Data: TestData{
			Users: []User{
				{ID: 1, Name: "sai"},
				{ID: 2, Name: "trim21"},
			},
			B:   false,
			Obj: Inner{V: 2, S: "vvv"},

			Map: map[int]struct{ V int }{7: {V: 4}},
		}},
	}

	for _, data := range testCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			b, err := phpserialize.Marshal(data.Data)
			require.NoError(t, err)

			j := decodeWithRealPhp(t, b)

			require.JSONEq(t, jsonEncode(t, data.Data), j)
		})
	}
}

// some special case like `empty map`, can't be compared by json unmarshal
func TestMarshal_special(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		b, err := phpserialize.Marshal(map[int]string{})
		require.NoError(t, err)
		require.Equal(t, []byte("a:0:{}"), b)
	})
}

func decodeWithRealPhp(t *testing.T, s []byte) string {
	os.MkdirAll("./tmp/", 0700)
	fs := filepath.Join("./tmp/", strings.ReplaceAll(t.Name(), "/", "-")) + ".php"

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
