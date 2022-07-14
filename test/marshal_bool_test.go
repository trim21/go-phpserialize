package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/test"
	"github.com/volatiletech/null/v9"
)

func TestMarshalBool_ptr_as_string(t *testing.T) {
	t.Run("direct-false", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
		}{
			Value: null.BoolFrom(false).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";s:5:"false";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("direct-true", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
		}{
			Value: null.BoolFrom(true).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:1:{s:5:"value";s:4:"true";}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("indirect-false", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
			B     *bool
		}{
			Value: null.BoolFrom(false).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:2:{s:5:"value";s:5:"false";s:1:"B";N;}`
		test.StringEqual(t, expected, string(actual))
	})

	t.Run("indirect-true", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
			B     *bool
		}{
			Value: null.BoolFrom(true).Ptr(),
		}

		actual, err := phpserialize.Marshal(&data)
		require.NoError(t, err)
		expected := `a:2:{s:5:"value";s:4:"true";s:1:"B";N;}`
		test.StringEqual(t, expected, string(actual))
	})
}
