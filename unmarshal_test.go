package phpserialize_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize/internal/decoder"
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

func tTestUnmarshal_struct(t *testing.T) {
	t.Parallel()

	b, err := os.ReadFile("./testdata/obj.txt")
	require.NoError(t, err)
	var o struct {
		S string  `php:"a string value"`
		I int     `php:"a int value"`
		B bool    `php:"a bool value"`
		F float64 `php:"a float value"`
	}

	err = decoder.Unmarshal(b, &o)
	require.NoError(t, err)

	require.Equal(t, "ff", o.S)
	require.Equal(t, int(31415926), o.I)
	require.Equal(t, true, o.B)
}
