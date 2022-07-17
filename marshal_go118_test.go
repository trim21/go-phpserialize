//go:build go1.18

package phpserialize_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/internal/test"
)

type Generic[T any] struct {
	Value T
}

type Generic2[T any] struct {
	B     bool // prevent direct
	Value T
}

var go118TestCase = []test.Case{
	{
		Name:     "generic[int]",
		Data:     Generic[int]{1},
		Expected: `a:1:{s:5:"Value";i:1;}`,
	},
	{
		Name:     "generic[struct]",
		Data:     Generic[test.User]{test.User{}},
		Expected: `a:1:{s:5:"Value";a:2:{s:2:"id";i:0;s:4:"name";s:0:"";}}`,
	},

	{
		Name:     "generic[map]",
		Data:     Generic[map[string]int]{map[string]int{"one": 1}},
		Expected: `a:1:{s:5:"Value";a:1:{s:3:"one";i:1;}}`,
	},

	{
		Name:     "generic[slice]",
		Data:     Generic[[]string]{[]string{"hello", "world"}},
		Expected: `a:1:{s:5:"Value";a:2:{i:0;s:5:"hello";i:1;s:5:"world";}}`,
	},

	{
		Name:     "generic2[slice]",
		Data:     Generic2[[]string]{Value: []string{"hello", "world"}},
		Expected: `a:2:{s:1:"B";b:0;s:5:"Value";a:2:{i:0;s:5:"hello";i:1;s:5:"world";}}`,
	},
}

func TestMarshal_go118_concrete_types(t *testing.T) {
	t.Parallel()
	for _, data := range go118TestCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			actual, err := phpserialize.Marshal(data.Data)
			require.NoError(t, err)

			test.StringEqual(t, data.Expected, string(actual))
		})
	}
}

func TestMarshal_go118_interface(t *testing.T) {
	t.Parallel()
	for _, data := range go118TestCase {
		data := data
		t.Run(data.Name, func(t *testing.T) {
			t.Parallel()
			actual, err := phpserialize.Marshal(data)
			require.NoError(t, err)

			test.StringEqual(t, data.WrappedExpected(), string(actual))
		})
	}
}
