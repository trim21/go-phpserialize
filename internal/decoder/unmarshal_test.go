package decoder

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsumeString(t *testing.T) {
	t.Parallel()
	raw := []byte(`s:14:"a string value";s:2:"ff";`)
	s, offset, err := consumeString(raw[2:])
	require.NoError(t, err)
	require.Equal(t, "a string value", s)
	require.Equal(t, []byte(`s:2:"ff";`), raw[2+offset:])
}

func BenchmarkConsumeString(b *testing.B) {
	raw := []byte(`s:14:"a string value";s:2:"ff";`)
	var s string
	var offset int
	var err error
	for i := 0; i < b.N; i++ {
		s, offset, err = consumeString(raw[2:])
	}
	runtime.KeepAlive(s)
	runtime.KeepAlive(offset)
	runtime.KeepAlive(err)
}

func TestConsumeBool(t *testing.T) {
	t.Parallel()
	raw := []byte(`0;`)
	s, offset, err := consumeBool(raw)
	require.NoError(t, err)
	require.Equal(t, false, s)
	require.Equal(t, 2, offset)

	s, offset, err = consumeBool([]byte(`1;`))
	require.NoError(t, err)
	require.Equal(t, true, s)
	require.Equal(t, 2, offset)

	_, _, err = consumeBool([]byte(`a;`))
	require.Error(t, err)
}
