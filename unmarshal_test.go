package phpserialize_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

/*
array(
  'a string value' => 'ff',
  'a int value' => 31415926,
  'a bool value' => true,
  'a float value' => 3.14,
  662 => 223,
)
*/

func TestUnmarshal_float(t *testing.T) {
	t.Parallel()

	type Container struct {
		F float64 `php:"f,string"`
	}

	t.Run("negative", func(t *testing.T) {
		var c Container
		raw := `a:1:{s:1:"f";d:3.14;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, 3.14, c.F)
	})

	t.Run("positive", func(t *testing.T) {
		var c Container
		raw := `a:1:{s:1:"f";d:1;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, 1, c.F)
	})

	t.Run("zero", func(t *testing.T) {
		var c Container
		raw := `a:1:{s:1:"f";d:-3.14;}`
		err := phpserialize.Unmarshal([]byte(raw), &c)
		require.NoError(t, err)
		require.Equal(t, -3.14, c.F)
	})
}

func TestUnmarshal_struct_empty(t *testing.T) {
	t.Parallel()

	type Container struct {
		F string `php:"f,string"`
	}

	var c Container
	raw := `a:0:{}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, "", c.F)
}

func TestUnmarshal_struct_string(t *testing.T) {
	t.Parallel()

	type Container struct {
		F string `php:"f1q,string"`
		V bool   `php:"1a9,string"`
	}

	var c Container
	raw := `a:1:{s:3:"f1q";s:10:"0147852369"}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, "0147852369", c.F)
}

func TestUnmarshal_struct_float(t *testing.T) {
	t.Parallel()

	type Container struct {
		F float64 `php:"f1q"`
	}

	var c Container
	raw := `a:1:{s:3:"f1q";d:147852369;}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	require.NoError(t, err)
	require.Equal(t, "0147852369", c.F)
}
