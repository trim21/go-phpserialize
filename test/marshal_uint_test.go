package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/test"
)

func TestMarshal_uint_string(t *testing.T) {
	var data = struct {
		A uint8  `php:"a,string"`
		B uint16 `php:"b,string"`
		C uint32 `php:"c,string"`
		D uint64 `php:"d,string"`
		E uint   `php:"e,string"`
	}{
		A: 2,
		B: 3,
		C: 0,
		D: 52,
		E: 110,
	}

	actual, err := phpserialize.Marshal(&data)
	require.NoError(t, err)
	expected := `a:5:{s:1:"a";s:1:"2";s:1:"b";s:1:"3";s:1:"c";s:1:"0";s:1:"d";s:2:"52";s:1:"e";s:3:"110";}`
	test.StringEqual(t, expected, string(actual))
}

func TestMarshal_uint_omitempty(t *testing.T) {
	var data = struct {
		A uint8  `php:"a,string,omitempty"`
		B uint16 `php:"b,string,omitempty"`
		C uint32 `php:"c,string,omitempty"`
		D uint64 `php:"d,string,omitempty"`
		E uint   `php:"e,string,omitempty"`
	}{}

	actual, err := phpserialize.Marshal(&data)
	require.NoError(t, err)
	expected := `a:0:{}`
	test.StringEqual(t, expected, string(actual))
}
