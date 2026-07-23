package phpserialize_test

import (
	"math"
	"testing"

	"github.com/trim21/go-phpserialize"
)

type rawPHPValue string

func (v *rawPHPValue) UnmarshalPHP(raw []byte) error {
	*v = rawPHPValue(raw)
	return nil
}

func TestMarshalNil(t *testing.T) {
	got, err := phpserialize.Marshal(nil)
	if err != nil {
		t.Fatalf("Marshal(nil) returned an error: %v", err)
	}
	if string(got) != "N;" {
		t.Fatalf("Marshal(nil) = %q, want %q", got, "N;")
	}
}

func TestMarshalNamedFloat32(t *testing.T) {
	type temperature float32

	got, err := phpserialize.Marshal(temperature(3.5))
	if err != nil {
		t.Fatalf("Marshal(named float32) returned an error: %v", err)
	}
	if string(got) != "d:3.5;" {
		t.Fatalf("Marshal(named float32) = %q, want %q", got, "d:3.5;")
	}

	got, err = phpserialize.Marshal(struct {
		Value temperature `php:"value,string"`
	}{Value: 3.5})
	if err != nil {
		t.Fatalf("Marshal(named float32 as string) returned an error: %v", err)
	}
	if string(got) != `a:1:{s:5:"value";s:3:"3.5";}` {
		t.Fatalf("Marshal(named float32 as string) = %q", got)
	}
}

func TestMarshalRejectsRecursiveValues(t *testing.T) {
	type recursiveMap map[string]recursiveMap
	type recursiveSlice []recursiveSlice

	t.Run("map through interface", func(t *testing.T) {
		value := map[string]any{}
		value["self"] = value
		if _, err := phpserialize.Marshal(value); err == nil {
			t.Fatal("Marshal accepted a recursive map")
		}
	})

	t.Run("slice through interface", func(t *testing.T) {
		value := make([]any, 1)
		value[0] = value
		if _, err := phpserialize.Marshal(value); err == nil {
			t.Fatal("Marshal accepted a recursive slice")
		}
	})

	t.Run("struct pointer", func(t *testing.T) {
		type node struct {
			Next *node
		}
		value := new(node)
		value.Next = value
		if _, err := phpserialize.Marshal(value); err == nil {
			t.Fatal("Marshal accepted a recursive struct")
		}
	})

	t.Run("named map", func(t *testing.T) {
		value := recursiveMap{}
		value["self"] = value
		if _, err := phpserialize.Marshal(value); err == nil {
			t.Fatal("Marshal accepted a recursive named map")
		}
	})

	t.Run("named slice", func(t *testing.T) {
		value := make(recursiveSlice, 1)
		value[0] = value
		if _, err := phpserialize.Marshal(value); err == nil {
			t.Fatal("Marshal accepted a recursive named slice")
		}
	})
}

func TestPHPFloatSpecialValues(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		serialized string
	}{
		{name: "positive infinity", value: math.Inf(1), serialized: "d:INF;"},
		{name: "negative infinity", value: math.Inf(-1), serialized: "d:-INF;"},
		{name: "not a number", value: math.NaN(), serialized: "d:NAN;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := phpserialize.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal(%v) returned an error: %v", tt.value, err)
			}
			if string(got) != tt.serialized {
				t.Fatalf("Marshal(%v) = %q, want %q", tt.value, got, tt.serialized)
			}

			var decoded float64
			if err := phpserialize.Unmarshal([]byte(tt.serialized), &decoded); err != nil {
				t.Fatalf("Unmarshal(%q) returned an error: %v", tt.serialized, err)
			}
			if math.IsNaN(tt.value) {
				if !math.IsNaN(decoded) {
					t.Fatalf("Unmarshal(%q) = %v, want NaN", tt.serialized, decoded)
				}
			} else if decoded != tt.value {
				t.Fatalf("Unmarshal(%q) = %v, want %v", tt.serialized, decoded, tt.value)
			}
		})
	}
}

func TestUnmarshalNilPointer(t *testing.T) {
	var dst *int
	if err := phpserialize.Unmarshal([]byte("i:1;"), dst); err == nil {
		t.Fatal("Unmarshal into a nil pointer returned nil error")
	}
}

func TestUnmarshalTopLevelBoundaries(t *testing.T) {
	var dst int
	if err := phpserialize.Unmarshal(nil, &dst); err == nil {
		t.Fatal("Unmarshal accepted empty input")
	}
	if err := phpserialize.Unmarshal([]byte("i:1;x"), &dst); err == nil {
		t.Fatal("Unmarshal accepted trailing data")
	}
}

func TestUnmarshalerReceivesFloat(t *testing.T) {
	for _, raw := range []string{"d:3.5;", "i:-5;"} {
		var dst rawPHPValue
		if err := phpserialize.Unmarshal([]byte(raw), &dst); err != nil {
			t.Fatalf("Unmarshal %q through Unmarshaler returned an error: %v", raw, err)
		}
		if string(dst) != raw {
			t.Fatalf("Unmarshaler received %q, want %q", dst, raw)
		}
	}
}

func TestUnmarshalNullIntoString(t *testing.T) {
	var dst string
	if err := phpserialize.Unmarshal([]byte("N;"), &dst); err != nil {
		t.Fatalf("Unmarshal null into string returned an error: %v", err)
	}
}

func TestUnmarshalIntegerValidation(t *testing.T) {
	t.Run("signed boundaries", func(t *testing.T) {
		for _, raw := range []string{"i:9223372036854775807;", "i:-9223372036854775808;"} {
			var dst int64
			if err := phpserialize.Unmarshal([]byte(raw), &dst); err != nil {
				t.Fatalf("Unmarshal(%q) returned an error: %v", raw, err)
			}
		}
	})

	t.Run("leading zeros", func(t *testing.T) {
		tests := []struct {
			raw  string
			want int64
		}{
			{raw: "i:000000000000000000001;", want: 1},
			{raw: "i:-000000000000000000001;", want: -1},
			{raw: "i:000000000000000000000;", want: 0},
			{raw: "i:-000000000000000000000;", want: 0},
		}
		for _, tt := range tests {
			var dst int64
			if err := phpserialize.Unmarshal([]byte(tt.raw), &dst); err != nil {
				t.Fatalf("Unmarshal(%q) returned an error: %v", tt.raw, err)
			}
			if dst != tt.want {
				t.Fatalf("Unmarshal(%q) = %d, want %d", tt.raw, dst, tt.want)
			}
		}
	})

	t.Run("signed overflow", func(t *testing.T) {
		var dst int64
		if err := phpserialize.Unmarshal([]byte("i:9223372036854775808;"), &dst); err == nil {
			t.Fatal("Unmarshal accepted an overflowing int64")
		}
	})

	t.Run("unsigned overflow", func(t *testing.T) {
		var dst uint64
		if err := phpserialize.Unmarshal([]byte("i:18446744073709551616;"), &dst); err == nil {
			t.Fatal("Unmarshal accepted an overflowing uint64")
		}
	})

	t.Run("negative unsigned", func(t *testing.T) {
		var dst uint64
		if err := phpserialize.Unmarshal([]byte("i:-1;"), &dst); err == nil {
			t.Fatal("Unmarshal accepted a negative uint64")
		}
	})

	t.Run("sign without digits", func(t *testing.T) {
		var dst int64
		if err := phpserialize.Unmarshal([]byte("i:-;"), &dst); err == nil {
			t.Fatal("Unmarshal accepted an integer sign without digits")
		}
	})

	t.Run("invalid wrapped string", func(t *testing.T) {
		var dst struct {
			Value int `php:"value,string"`
		}
		raw := `a:1:{s:5:"value";s:2:"1x";}`
		if err := phpserialize.Unmarshal([]byte(raw), &dst); err == nil {
			t.Fatal("Unmarshal accepted a non-decimal integer string")
		}
	})
}

func TestUnmarshalRejectsNegativeSequenceIndex(t *testing.T) {
	for _, dst := range []any{&[]string{}, &[1]string{}} {
		if err := phpserialize.Unmarshal([]byte(`a:1:{i:-1;s:1:"x";}`), dst); err == nil {
			t.Fatalf("Unmarshal accepted a negative index into %T", dst)
		}
	}
}

func TestUnmarshalTruncatedInputDoesNotPanic(t *testing.T) {
	tests := []struct {
		serialized string
		newDst     func() any
	}{
		{serialized: "b:1;", newDst: func() any { return new(bool) }},
		{serialized: "i:-123;", newDst: func() any { return new(int64) }},
		{serialized: "d:3.14;", newDst: func() any { return new(float64) }},
		{serialized: `s:5:"hello";`, newDst: func() any { return new(string) }},
		{serialized: `a:1:{i:0;s:1:"x";}`, newDst: func() any { return new([]string) }},
		{serialized: `a:1:{s:1:"x";i:1;}`, newDst: func() any { return new(map[string]int) }},
		{serialized: `a:1:{s:1:"X";i:1;}`, newDst: func() any { return new(struct{ X int }) }},
		{serialized: `O:8:"stdClass":1:{s:1:"X";i:1;}`, newDst: func() any { return new(struct{ X int }) }},
	}

	for _, tt := range tests {
		for end := 1; end < len(tt.serialized); end++ {
			raw := tt.serialized[:end]
			t.Run(raw, func(t *testing.T) {
				defer func() {
					if recovered := recover(); recovered != nil {
						t.Fatalf("Unmarshal(%q) panicked: %v", raw, recovered)
					}
				}()
				if err := phpserialize.Unmarshal([]byte(raw), tt.newDst()); err == nil {
					t.Fatalf("Unmarshal(%q) returned nil error", raw)
				}
			})
		}
	}
}

func TestUnmarshalOversizedStringLengthDoesNotPanic(t *testing.T) {
	inputs := []string{
		`s:9223372036854775807:"x";`,
		`O:9223372036854775807:"stdClass":0:{}`,
	}
	for _, raw := range inputs {
		t.Run(raw, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("Unmarshal(%q) panicked: %v", raw, recovered)
				}
			}()
			var dst any
			if err := phpserialize.Unmarshal([]byte(raw), &dst); err == nil {
				t.Fatalf("Unmarshal(%q) returned nil error", raw)
			}
		})
	}
}

func FuzzUnmarshalDoesNotPanic(f *testing.F) {
	seeds := []string{
		"N;",
		"b:1;",
		"i:-123;",
		"d:INF;",
		`s:5:"hello";`,
		`a:1:{i:0;s:1:"x";}`,
		`O:8:"stdClass":1:{s:1:"X";i:1;}`,
	}
	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, raw []byte) {
		destinations := []any{
			new(any),
			new(bool),
			new(int64),
			new(uint64),
			new(float64),
			new(string),
			new([]any),
			new(map[any]any),
			new(struct{ X any }),
		}
		for _, dst := range destinations {
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						t.Fatalf("Unmarshal(%q, %T) panicked: %v", raw, dst, recovered)
					}
				}()
				_ = phpserialize.Unmarshal(raw, dst)
			}()
		}
	})
}
