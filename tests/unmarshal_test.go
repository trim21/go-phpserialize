package tests

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
	t.Parallel()

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
		require.NoError(t, phpserialize.Unmarshal([]byte(raw), &c))
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
	t.Parallel()
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
	require.NoError(t, phpserialize.Unmarshal([]byte(raw), &c))
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	data := `O:8:"stdClass":1:{s:1:"a";b:0;}`

	var actual map[string]interface{}

	err := phpserialize.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"a": false,
	}
	require.Equal(t, expected, actual)
}

type Tag struct {
	Name       *string `php:"tag_name"`
	Count      int     `php:"result,string"`
	TotalCount int     `php:"tag_results,string"`
}

func TestUnmarshal_error_case(t *testing.T) {
	raw := `a:30:{i:0;a:2:{s:8:"tag_name";s:18:"叛逆的鲁鲁修";s:6:"result";s:4:"1645";}i:1;a:2:{s:8:"tag_name";s:7:"SUNRISE";s:6:"result";s:4:"1229";}i:2;a:2:{s:8:"tag_name";s:24:"反逆のルルーシュ";s:6:"result";s:3:"936";}i:3;a:2:{s:8:"tag_name";s:15:"还是死妹控";s:6:"result";s:3:"721";}i:4;a:2:{s:8:"tag_name";s:2:"TV";s:6:"result";s:3:"664";}i:5;a:2:{s:8:"tag_name";s:6:"妹控";s:6:"result";s:3:"603";}i:6;a:2:{s:8:"tag_name";s:9:"codegeass";s:6:"result";s:3:"569";}i:7;a:2:{s:8:"tag_name";s:12:"谷口悟朗";s:6:"result";s:3:"523";}i:8;a:2:{s:8:"tag_name";s:9:"鲁路修";s:6:"result";s:3:"453";}i:9;a:2:{s:8:"tag_name";s:2:"R2";s:6:"result";s:3:"427";}i:10;a:2:{s:8:"tag_name";s:4:"2008";s:6:"result";s:3:"409";}i:11;a:2:{s:8:"tag_name";s:6:"原创";s:6:"result";s:3:"385";}i:12;a:2:{s:8:"tag_name";s:11:"2008年4月";s:6:"result";s:3:"357";}i:13;a:2:{s:8:"tag_name";s:15:"大河内一楼";s:6:"result";s:3:"174";}i:14;a:2:{s:8:"tag_name";s:6:"日升";s:6:"result";s:3:"151";}i:15;a:2:{s:8:"tag_name";s:6:"萝卜";s:6:"result";s:3:"120";}i:16;a:2:{s:8:"tag_name";s:6:"机战";s:6:"result";s:3:"111";}i:17;a:2:{s:8:"tag_name";s:15:"狗得鸡鸭死";s:6:"result";s:3:"104";}i:18;a:2:{s:8:"tag_name";s:9:"福山润";s:6:"result";s:2:"94";}i:19;a:2:{s:8:"tag_name";s:9:"露露胸";s:6:"result";s:2:"84";}i:20;a:2:{s:8:"tag_name";s:5:"CLAMP";s:6:"result";s:2:"69";}i:21;a:2:{s:8:"tag_name";s:6:"科幻";s:6:"result";s:2:"67";}i:22;a:2:{s:8:"tag_name";s:9:"鲁鲁修";s:6:"result";s:2:"62";}i:23;a:2:{s:8:"tag_name";s:5:"GEASS";s:6:"result";s:2:"57";}i:24;a:2:{s:8:"tag_name";s:6:"神作";s:6:"result";s:2:"54";}i:25;a:2:{s:8:"tag_name";s:6:"战斗";s:6:"result";s:2:"49";}i:26;a:2:{s:8:"tag_name";s:6:"战争";s:6:"result";s:2:"41";}i:27;a:2:{s:8:"tag_name";s:25:"裸露修的跌二次KUSO";s:6:"result";s:2:"40";}i:28;a:2:{s:8:"tag_name";s:6:"中二";s:6:"result";s:2:"37";}i:29;a:2:{s:8:"tag_name";s:12:"樱井孝宏";s:6:"result";s:2:"34";}}`
	var tags []Tag
	err := phpserialize.Unmarshal([]byte(raw), &tags)
	require.NoError(t, err)
}
