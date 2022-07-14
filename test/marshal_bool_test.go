package test_test

import (
	"testing"

	"github.com/volatiletech/null/v9"
)

func TestMarshalBool_ptr_as_string(t *testing.T) {
	t.Run("direct-false", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
		}{
			Value: null.BoolFrom(false).Ptr(),
		}
		expected := `a:1:{s:5:"value";s:5:"false";}`

		MarshalExpected(t, data, expected)
	})

	t.Run("direct-true", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
		}{
			Value: null.BoolFrom(true).Ptr(),
		}
		expected := `a:1:{s:5:"value";s:4:"true";}`

		MarshalExpected(t, data, expected)
	})

	t.Run("indirect-false", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
			B     *bool
		}{
			Value: null.BoolFrom(false).Ptr(),
		}

		expected := `a:2:{s:5:"value";s:5:"false";s:1:"B";N;}`

		MarshalExpected(t, data, expected)
	})

	t.Run("indirect-true", func(t *testing.T) {
		var data = struct {
			Value *bool `php:"value,string"`
			B     *bool
		}{
			Value: null.BoolFrom(true).Ptr(),
		}

		expected := `a:2:{s:5:"value";s:4:"true";s:1:"B";N;}`

		MarshalExpected(t, data, expected)
	})
}
