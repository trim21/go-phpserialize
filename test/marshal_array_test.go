package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/test"
)

func TestMarshal_array_map(t *testing.T) {
	var data = [5]map[int]uint{
		{
			-3: 1,
			4:  8,
		},
		nil,
		{-1: 1},
	}

	actual, err := phpserialize.Marshal(data)
	require.NoError(t, err)
	expected := `a:5:{i:0;a:2:{i:-3;i:1;i:4;i:8;}i:1;N;i:2;a:1:{i:-1;i:1;}i:3;N;i:4;N;}`
	test.StringEqual(t, expected, string(actual))
}

func TestMarshal_Array_nil(t *testing.T) {
	var data [5]int

	actual, err := phpserialize.Marshal(data)
	require.NoError(t, err)
	expected := `a:5:{i:0;i:0;i:1;i:0;i:2;i:0;i:3;i:0;i:4;i:0;}`
	test.StringEqual(t, expected, string(actual))
}
