package phpserialize_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

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
