package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/test"
)

type TestCase struct {
	Name     string
	Data     interface{}
	Expected string `php:"-" json:"-"`
}

func MarshalExpected(t *testing.T, data interface{}, expected string) {
	t.Helper()
	actual, err := phpserialize.Marshal(data)
	require.NoError(t, err)
	if string(actual) != expected {
		t.Errorf("Result not as expected:\n%v", test.CharacterDiff(expected, string(actual)))
		t.FailNow()
	}
}
