package decoder

import (
	"errors"
	"strconv"
	"testing"
)

func FuzzParseInt64MatchesStrconv(f *testing.F) {
	seeds := []string{
		"0",
		"-0",
		"7",
		"-7",
		"000000000000000000001",
		"-000000000000000000001",
		"9223372036854775807",
		"-9223372036854775808",
		"9223372036854775808",
		"-9223372036854775809",
		"",
		"-",
		"1x",
	}
	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, raw []byte) {
		got, gotErr := parseInt64(raw)

		// PHP integer payloads permit a leading minus, but not a leading plus.
		want, wantErr := strconv.ParseInt(string(raw), 10, 64)
		if len(raw) > 0 && raw[0] == '+' {
			wantErr = errors.New("leading plus is not valid PHP integer syntax")
		}

		if (gotErr != nil) != (wantErr != nil) {
			t.Fatalf("parseInt64(%q) error = %v, strconv error = %v", raw, gotErr, wantErr)
		}
		if gotErr == nil && got != want {
			t.Fatalf("parseInt64(%q) = %d, want %d", raw, got, want)
		}
	})
}
