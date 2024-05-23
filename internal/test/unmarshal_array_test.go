package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

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
