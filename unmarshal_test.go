package phpserialize_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trim21/go-phpserialize"
)

func TestUnmarshal_as_string(t *testing.T) {
	t.Parallel()

	t.Run("ptr", func(t *testing.T) {
		type Container struct {
			V *int `php:",omitempty,string"`
		}

		var c Container
		raw := `a:1:{s:1:"V";s:1:"1";}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, 1, *c.V)
	})
}

func TestUnmarshal_struct_string(t *testing.T) {
	t.Parallel()

	t.Run("value", func(t *testing.T) {
		type Container struct {
			F string `php:"f1q"`
			V bool   `php:"1a9"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";s:10:"0147852369";}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, "0147852369", c.F)
	})

	t.Run("empty", func(t *testing.T) {
		type Container struct {
			F string `php:"f"`
		}

		var c Container
		raw := `a:0:{}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, "", c.F)
	})
}

func TestUnmarshal_stdClass(t *testing.T) {
	raw := `O:8:"stdClass":1:{s:1:"a";s:13:"a str value q";}`

	t.Run("struct", func(t *testing.T) {
		var v struct {
			A string `php:"a"`
		}

		require.NoError(t, phpserialize.Unmarshal([]byte(raw), &v))

		require.Equal(t, "a str value q", v.A)
	})

	t.Run("any", func(t *testing.T) {
		var v any
		require.NoError(t, phpserialize.Unmarshal([]byte(raw), &v))

		m, ok := v.(map[string]any)
		require.True(t, ok, "type cast fail", v)

		require.Equal(t, "a str value q", m["a"])
	})

	t.Run("skip", func(t *testing.T) {
		raw := `a:2:{s:1:"a";O:8:"stdClass":1:{s:1:"a";s:13:"a str value q";}s:5:"value";b:1;}`
		var v struct {
			Value bool `php:"value"`
		}
		require.NoError(t, phpserialize.Unmarshal([]byte(raw), &v))

		require.True(t, v.Value)
	})
}

func TestUnmarshal_struct_bytes(t *testing.T) {
	t.Parallel()

	t.Run("value", func(t *testing.T) {
		type Container struct {
			F []byte `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";s:10:"0147852369";}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, []byte("0147852369"), c.F)
	})

	t.Run("empty", func(t *testing.T) {
		type Container struct {
			F []byte `php:"f"`
		}

		var c Container
		raw := `a:0:{}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Nil(t, c.F)
	})
}

func TestUnmarshal_struct_float(t *testing.T) {
	t.Parallel()

	t.Run("negative", func(t *testing.T) {
		type Container struct {
			F float64 `php:"f"`
		}
		var c Container
		raw := `a:1:{s:1:"f";d:3.14;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, float64(3.14), c.F)
	})

	t.Run("positive", func(t *testing.T) {
		type Container struct {
			F float64 `php:"f"`
		}
		var c Container
		raw := `a:1:{s:1:"f";d:1;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, float64(1), c.F)
	})

	t.Run("zero", func(t *testing.T) {
		type Container struct {
			F float64 `php:"f"`
		}
		var c Container
		raw := `a:1:{s:1:"f";d:-3.14;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, float64(-3.14), c.F)
	})

	t.Run("float32", func(t *testing.T) {
		type Container struct {
			F float32 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";d:147852369;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, float32(147852369), c.F)
	})

	t.Run("float64", func(t *testing.T) {
		type Container struct {
			F float64 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";d:147852369;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, float64(147852369), c.F)
	})
}

func TestUnmarshal_struct_uint(t *testing.T) {
	t.Parallel()

	t.Run("uint", func(t *testing.T) {
		type Container struct {
			F uint `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:147852369;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, uint(147852369), c.F)
	})

	t.Run("uint8", func(t *testing.T) {
		type Container struct {
			F uint8 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:255;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, uint8(255), c.F)
	})

	t.Run("uint16", func(t *testing.T) {
		type Container struct {
			F uint16 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:574;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, uint16(574), c.F)
	})

	t.Run("uint32", func(t *testing.T) {
		type Container struct {
			F uint32 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:57400;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, uint32(57400), c.F)
	})

	t.Run("uint64", func(t *testing.T) {
		type Container struct {
			F uint64 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:5740000;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, uint64(5740000), c.F)
	})
}

func TestUnmarshal_struct_int(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		type Container struct {
			F int `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:147852369;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, int(147852369), c.F)
	})

	t.Run("int8", func(t *testing.T) {
		type Container struct {
			F int8 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:65;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, int8(65), c.F)
	})

	t.Run("int16", func(t *testing.T) {
		type Container struct {
			F int16 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:574;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, int16(574), c.F)
	})

	t.Run("int32", func(t *testing.T) {
		type Container struct {
			F int32 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:57400;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, int32(57400), c.F)
	})

	t.Run("int64", func(t *testing.T) {
		type Container struct {
			F int64 `php:"f1q"`
		}

		var c Container
		raw := `a:1:{s:3:"f1q";i:5740000;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, int64(5740000), c.F)
	})
}

func TestUnmarshal_slice(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		type Container struct {
			Value []string `php:"value"`
		}

		var c Container
		raw := `a:1:{s:5:"value";a:0:{}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Len(t, c.Value, 0)
	})

	t.Run("string", func(t *testing.T) {
		type Container struct {
			Value []string `php:"value"`
		}
		var c Container
		raw := `a:3:{s:2:"bb";b:1;s:5:"value";a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, []string{"one", "two", "q"}, c.Value)
	})

	t.Run("string more length", func(t *testing.T) {
		type Container struct {
			Value []string `php:"value"`
		}
		var c Container
		raw := `a:1:{s:5:"value";a:6:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";i:3;s:1:"a";i:4;s:2:"zx";i:5;s:3:"abc";}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, []string{"one", "two", "q", "a", "zx"}, c.Value[:5])
	})
}

func TestUnmarshal_array(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		type Container struct {
			Value [5]string `php:"value"`
		}

		var c Container
		raw := `a:1:{s:5:"value";a:0:{}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, [5]string{}, c.Value)
	})

	t.Run("string less length", func(t *testing.T) {
		type Container struct {
			Value [5]string `php:"value"`
		}
		var c Container
		raw := `a:1:{s:5:"value";a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, [5]string{"one", "two", "q"}, c.Value)
	})

	t.Run("string more length", func(t *testing.T) {
		type Container struct {
			Value [5]string `php:"value"`
		}
		var c Container
		raw := `a:1:{s:5:"value";a:6:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";i:3;s:1:"a";i:4;s:2:"zx";i:5;s:3:"abc";}}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, [5]string{"one", "two", "q", "a", "zx"}, c.Value)
	})
}

func TestUnmarshal_skip_value(t *testing.T) {
	type Container struct {
		Value []string `php:"value"`
	}

	var c Container
	raw := `a:3:{s:2:"bb";b:1;s:5:"value";a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}s:6:"value2";a:3:{i:0;s:1:"1";i:1;s:1:"2";i:2;s:1:"3";}}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two", "q"}, c.Value)
}

var _ phpserialize.Unmarshaler = (*unmarshaler)(nil)

type unmarshaler []byte

func (u *unmarshaler) UnmarshalPHP(bytes []byte) error {
	*u = append((*u)[0:0], bytes...)

	return nil
}

func TestUnmarshal_unmarshaler(t *testing.T) {
	t.Parallel()
	type Container struct {
		Value unmarshaler `php:"value"`
	}

	var c Container
	raw := `a:1:{s:5:"value";a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, `a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}`, string(c.Value))
}

func TestUnmarshal_string_wrapper(t *testing.T) {
	t.Parallel()

	type Container struct {
		Value int `php:"value,string"`
	}

	var c Container
	raw := `a:1:{s:5:"value";s:3:"233";}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, int(233), c.Value)
}

func TestUnmarshal_map(t *testing.T) {
	t.Parallel()

	t.Run("map[string]string", func(t *testing.T) {
		raw := `a:1:{s:5:"value";a:5:{s:3:"one";s:1:"1";s:3:"two";s:1:"2";s:5:"three";s:1:"3";s:4:"four";s:1:"4";s:4:"five";s:1:"5";}}`
		var c struct {
			Value map[string]string `php:"value"`
		}

		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"one":   "1",
			"two":   "2",
			"three": "3",
			"four":  "4",
			"five":  "5",
		}, c.Value)
	})

	t.Run("map[any]string", func(t *testing.T) {
		raw := `a:1:{s:5:"value";a:5:{i:1;s:3:"one";i:2;s:3:"two";i:3;s:5:"three";i:4;s:4:"four";i:5;s:4:"five";}}`
		var c struct {
			Value map[any]string `php:"value"`
		}

		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, map[any]string{
			int64(1): "one",
			int64(2): "two",
			int64(3): "three",
			int64(4): "four",
			int64(5): "five",
		}, c.Value)
	})

	t.Run("any", func(t *testing.T) {
		raw := `a:1:{s:5:"value";a:5:{i:1;s:3:"one";i:2;s:3:"two";i:3;s:5:"three";i:4;s:4:"four";i:5;s:4:"five";}}`
		var c struct {
			Value any `php:"value"`
		}

		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, map[any]any{
			int64(1): "one",
			int64(2): "two",
			int64(3): "three",
			int64(4): "four",
			int64(5): "five",
		}, c.Value)
	})
}

func TestUnmarshal_ptr_string(t *testing.T) {
	t.Parallel()

	t.Run("value", func(t *testing.T) {
		var c struct {
			F *string `php:"f1q"`
		}

		raw := `a:1:{s:3:"f1q";s:10:"0147852369";}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.NotNil(t, c.F)
		require.Equal(t, "0147852369", *c.F)
	})

	t.Run("empty", func(t *testing.T) {
		var c struct {
			F *string `php:"f"`
		}

		raw := `a:0:{}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Nil(t, c.F)
	})

	t.Run("nested", func(t *testing.T) {
		var c struct {
			F **string `php:"f"`
		}

		raw := `a:0:{}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.Error(t, err)
	})
}

func TestUnmarshal_anonymous_field(t *testing.T) {
	type N struct {
		A int
		B int
	}

	type M struct {
		N
		C int
	}

	var v M

	require.Error(t, phpserialize.Unmarshal([]byte(`a:4:{s:1:"A";i:3;s:1:"B";i:2;s:1:"C";i:1;}`), &v))
}

func TestUnmarshal_empty_input(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		var data []int
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
	t.Run("array", func(t *testing.T) {
		var data [5]int
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
	t.Run("map", func(t *testing.T) {
		var data map[uint]int
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
	t.Run("interface", func(t *testing.T) {
		var data any
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
	t.Run("string", func(t *testing.T) {
		var data string
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
	t.Run("int", func(t *testing.T) {
		var data int
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})

	t.Run("uint", func(t *testing.T) {
		var data uint
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})

	t.Run("bool", func(t *testing.T) {
		var data bool
		require.Error(t, phpserialize.Unmarshal([]byte(""), &data))
	})
}

func TestUnmarshal_as_string_2(t *testing.T) {
	type ID uint32
	type Type uint8
	type Item struct {
		ID   ID   `php:"eid,string"`
		Type Type `php:"type"`
	}
	type Collection = map[ID]Item

	raw := `a:7:{i:1087180;a:2:{s:3:"eid";s:7:"1087180";s:4:"type";i:2;}i:1087181;a:2:{s:3:"eid";s:7:"1087181";s:4:"type";i:2;}i:1087182;a:2:{s:3:"eid";s:7:"1087182";s:4:"type";i:2;}i:1087183;a:2:{s:3:"eid";s:7:"1087183";s:4:"type";i:2;}i:1087184;a:2:{s:3:"eid";s:7:"1087184";s:4:"type";i:2;}i:1087185;a:2:{s:3:"eid";s:7:"1087185";s:4:"type";i:2;}i:1087186;a:2:{s:3:"eid";s:7:"1087186";s:4:"type";i:2;}}`

	var data Collection

	err := phpserialize.Unmarshal([]byte(raw), &data)
	require.NoError(t, err)
}

func TestUnmarshal_null_array_1(t *testing.T) {
	raw := `a:0:{}`

	type Tag struct {
		Name  *string `php:"tag_name"`
		Count int     `php:"result,string"`
	}

	var tags []Tag

	err := phpserialize.Unmarshal([]byte(raw), &tags)
	require.NoError(t, err)
}

func TestUnmarshal_null_array_2(t *testing.T) {
	raw := `a:4:{s:1:"a";i:2;s:4:"Test";a:0:{}s:1:"b";a:0:{}s:1:"o";i:1;}`

	var data any

	err := phpserialize.Unmarshal([]byte(raw), &data)
	require.NoError(t, err)

	require.Equal(t, data, map[any]any{
		"a":    int64(2),
		"o":    int64(1),
		"Test": map[any]any{},
		"b":    map[any]any{},
	})
}

func TestUnmarshal_array_with_bool_to_map(t *testing.T) {
	data := `O:8:"stdClass":1:{s:1:"a";b:0;}`

	var actual map[string]interface{}

	err := phpserialize.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"a": false,
	}
	require.Equal(t, expected, actual)
}
