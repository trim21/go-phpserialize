package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/test"
	"github.com/volatiletech/null/v9"
)

func TestMarshal_int_ptr_string(t *testing.T) {
	t.Run("ptr direct", func(t *testing.T) {
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
}
