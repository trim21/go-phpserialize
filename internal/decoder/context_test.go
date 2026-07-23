package decoder

import "testing"

func TestReadLength(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantEnd int64
		wantErr bool
	}{
		{input: ":0:", want: 0, wantEnd: 3},
		{input: ":12:", want: 12, wantEnd: 4},
		{input: "", wantErr: true},
		{input: "x", wantErr: true},
		{input: "::", wantErr: true},
		{input: ":1", wantErr: true},
		{input: ":1x", wantErr: true},
		{input: ":9223372036854775808:", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, end, err := readLength([]byte(tt.input), 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("readLength(%q) error = %v", tt.input, err)
			}
			if err == nil && (got != tt.want || end != tt.wantEnd) {
				t.Fatalf("readLength(%q) = (%d, %d), want (%d, %d)", tt.input, got, end, tt.want, tt.wantEnd)
			}
		})
	}

	if _, _, err := readLengthInt([]byte(":12:"), 0); err != nil {
		t.Fatalf("readLengthInt returned an error: %v", err)
	}
	if _, _, err := readLengthInt([]byte(":9223372036854775808:"), 0); err == nil {
		t.Fatal("readLengthInt accepted an overflowing length")
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantEnd int64
		wantErr bool
	}{
		{input: `:0:"";`, wantEnd: 6},
		{input: `:3:"foo";`, want: "foo", wantEnd: 9},
		{input: `:3:foo";`, wantErr: true},
		{input: `:3:"fo`, wantErr: true},
		{input: `:3:"foo!;`, wantErr: true},
		{input: `:3:"foo"!`, wantErr: true},
		{input: `:9223372036854775807:"x";`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, end, err := readString([]byte(tt.input), 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("readString(%q) error = %v", tt.input, err)
			}
			if err == nil && (string(got) != tt.want || end != tt.wantEnd) {
				t.Fatalf("readString(%q) = (%q, %d), want (%q, %d)", tt.input, got, end, tt.want, tt.wantEnd)
			}
		})
	}
}

func TestReadScalarTokens(t *testing.T) {
	integerTests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{input: "i:0;", want: "0"},
		{input: "i:-12;", want: "-12"},
		{input: "", wantErr: true},
		{input: "x:1;", wantErr: true},
		{input: "i", wantErr: true},
		{input: "i!1;", wantErr: true},
		{input: "i:;", wantErr: true},
		{input: "i:-;", wantErr: true},
		{input: "i:12", wantErr: true},
		{input: "i:1x", wantErr: true},
	}
	for _, tt := range integerTests {
		got, _, err := readIntegerBytes([]byte(tt.input), 0)
		if (err != nil) != tt.wantErr || err == nil && string(got) != tt.want {
			t.Errorf("readIntegerBytes(%q) = %q, %v", tt.input, got, err)
		}
	}

	floatTests := []struct {
		input   string
		wantErr bool
	}{
		{input: "d:3.5;"},
		{input: "", wantErr: true},
		{input: "x:1;", wantErr: true},
		{input: "d", wantErr: true},
		{input: "d!1;", wantErr: true},
		{input: "d:;", wantErr: true},
		{input: "d:1", wantErr: true},
	}
	for _, tt := range floatTests {
		_, _, err := readFloatBytes([]byte(tt.input), 0)
		if (err != nil) != tt.wantErr {
			t.Errorf("readFloatBytes(%q) error = %v", tt.input, err)
		}
	}

	for _, tt := range []struct {
		input   string
		want    bool
		wantErr bool
	}{
		{input: "b:0;"},
		{input: "b:1;", want: true},
		{input: "b:2;", wantErr: true},
		{input: "b!1;", wantErr: true},
		{input: "x:1;", wantErr: true},
		{input: "b:1", wantErr: true},
	} {
		got, err := readBool([]byte(tt.input), 0)
		if (err != nil) != tt.wantErr || err == nil && got != tt.want {
			t.Errorf("readBool(%q) = %v, %v", tt.input, got, err)
		}
	}
}

func TestSkipValue(t *testing.T) {
	valid := []string{
		"N;",
		"b:1;",
		"i:-1;",
		"d:INF;",
		`s:3:"foo";`,
		`a:1:{s:1:"x";i:1;}`,
		`O:8:"stdClass":1:{s:1:"x";i:1;}`,
	}
	for _, input := range valid {
		end, err := skipValue([]byte(input), 0, 0)
		if err != nil {
			t.Errorf("skipValue(%q) returned an error: %v", input, err)
		} else if end != int64(len(input)) {
			t.Errorf("skipValue(%q) ended at %d", input, end)
		}
	}

	invalid := []string{
		"",
		"x",
		"N",
		`s:2:"x";`,
		`a:1:{s:1:"x";}`,
		`a:0:{i:0;i:1;}`,
		`O:8:"stdClass":1:{s:1:"x";}`,
	}
	for _, input := range invalid {
		if _, err := skipValue([]byte(input), 0, 0); err == nil {
			t.Errorf("skipValue(%q) returned nil error", input)
		}
	}

	if _, err := skipArray([]byte(":0:{}"), 0, maxDecodeNestingDepth+1); err == nil {
		t.Fatal("skipArray accepted excessive nesting")
	}
}

func TestSmallTokenValidators(t *testing.T) {
	for _, input := range []string{"", "N", "Nx"} {
		if err := validateNull([]byte(input), 0); err == nil {
			t.Errorf("validateNull(%q) returned nil error", input)
		}
	}
	if err := validateNull([]byte("N;"), 0); err != nil {
		t.Fatalf("validateNull returned an error: %v", err)
	}

	for _, input := range []string{"0", "00{}", "0:x}", "0:{x"} {
		if err := validateEmptyArray([]byte(input), 0); err == nil {
			t.Errorf("validateEmptyArray(%q) returned nil error", input)
		}
	}
	if err := validateEmptyArray([]byte("0:{}"), 0); err != nil {
		t.Fatalf("validateEmptyArray returned an error: %v", err)
	}

	if hasByte(nil, 0) || hasBytes([]byte("x"), -1, 1) || hasBytes([]byte("x"), 0, -1) {
		t.Fatal("bounds helpers accepted invalid bounds")
	}
}
