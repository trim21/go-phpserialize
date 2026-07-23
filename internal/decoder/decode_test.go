package decoder

import (
	"errors"
	"math"
	"reflect"
	"testing"

	internalerrors "github.com/trim21/go-phpserialize/internal/errors"
)

func decodeTestValue(data string, dst any) (int64, error) {
	dec, err := CompileToGetDecoder(reflect.TypeOf(dst))
	if err != nil {
		return 0, err
	}
	ctx := TakeRuntimeContext()
	defer ReleaseRuntimeContext(ctx)
	ctx.Buf = []byte(data)
	return dec.Decode(ctx, 0, 0, reflect.ValueOf(dst).Elem())
}

func TestDecodeScalarTypes(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		dst  any
		want any
	}{
		{name: "bool", raw: "b:1;", dst: new(bool), want: true},
		{name: "int", raw: "i:-1;", dst: new(int), want: -1},
		{name: "int8", raw: "i:-8;", dst: new(int8), want: int8(-8)},
		{name: "int16", raw: "i:-16;", dst: new(int16), want: int16(-16)},
		{name: "int32", raw: "i:-32;", dst: new(int32), want: int32(-32)},
		{name: "int64", raw: "i:-64;", dst: new(int64), want: int64(-64)},
		{name: "uint", raw: "i:1;", dst: new(uint), want: uint(1)},
		{name: "uint8", raw: "i:8;", dst: new(uint8), want: uint8(8)},
		{name: "uint16", raw: "i:16;", dst: new(uint16), want: uint16(16)},
		{name: "uint32", raw: "i:32;", dst: new(uint32), want: uint32(32)},
		{name: "uint64", raw: "i:64;", dst: new(uint64), want: uint64(64)},
		{name: "uintptr", raw: "i:9;", dst: new(uintptr), want: uintptr(9)},
		{name: "float32", raw: "d:3.5;", dst: new(float32), want: float32(3.5)},
		{name: "float64", raw: "d:-4.5;", dst: new(float64), want: -4.5},
		{name: "string", raw: `s:3:"foo";`, dst: new(string), want: "foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			end, err := decodeTestValue(tt.raw, tt.dst)
			if err != nil {
				t.Fatalf("decode returned an error: %v", err)
			}
			if end != int64(len(tt.raw)) {
				t.Fatalf("decode ended at %d, want %d", end, len(tt.raw))
			}
			got := reflect.ValueOf(tt.dst).Elem().Interface()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("decode result = %#v, want %#v", got, tt.want)
			}
		})
	}

	for _, raw := range []string{"d:INF;", "d:-INF;", "d:NAN;"} {
		var dst float64
		if _, err := decodeTestValue(raw, &dst); err != nil {
			t.Fatalf("decode %q returned an error: %v", raw, err)
		}
		if raw == "d:NAN;" && !math.IsNaN(dst) {
			t.Fatalf("decode %q = %v", raw, dst)
		}
	}
}

func TestDecodeCompositeTypes(t *testing.T) {
	t.Run("bytes from string", func(t *testing.T) {
		var dst []byte
		if _, err := decodeTestValue(`s:3:"foo";`, &dst); err != nil || string(dst) != "foo" {
			t.Fatalf("decode bytes = %q, %v", dst, err)
		}
	})
	t.Run("bytes from array", func(t *testing.T) {
		var dst []byte
		raw := `a:2:{i:0;i:65;i:1;i:66;}`
		if _, err := decodeTestValue(raw, &dst); err != nil || string(dst) != "AB" {
			t.Fatalf("decode bytes = %q, %v", dst, err)
		}
	})
	t.Run("slice", func(t *testing.T) {
		var dst []string
		raw := `a:2:{i:0;s:1:"a";i:2;s:1:"c";}`
		if _, err := decodeTestValue(raw, &dst); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(dst, []string{"a", "", "c"}) {
			t.Fatalf("decode slice = %#v", dst)
		}
	})
	t.Run("array", func(t *testing.T) {
		var dst [1]string
		raw := `a:2:{i:0;s:1:"a";i:2;s:1:"c";}`
		if _, err := decodeTestValue(raw, &dst); err != nil || dst[0] != "a" {
			t.Fatalf("decode array = %#v, %v", dst, err)
		}
	})
	t.Run("map", func(t *testing.T) {
		var dst map[string]*int
		raw := `a:1:{s:1:"x";i:2;}`
		if _, err := decodeTestValue(raw, &dst); err != nil || dst["x"] == nil || *dst["x"] != 2 {
			t.Fatalf("decode map = %#v, %v", dst, err)
		}
	})
	t.Run("object map", func(t *testing.T) {
		var dst map[string]any
		raw := `O:8:"stdClass":1:{s:1:"x";i:2;}`
		if _, err := decodeTestValue(raw, &dst); err != nil || dst["x"] != int64(2) {
			t.Fatalf("decode object map = %#v, %v", dst, err)
		}
	})
	t.Run("pointer", func(t *testing.T) {
		var dst *int
		if _, err := decodeTestValue("i:2;", &dst); err != nil || dst == nil || *dst != 2 {
			t.Fatalf("decode pointer = %#v, %v", dst, err)
		}
		if _, err := decodeTestValue("N;", &dst); err != nil || dst != nil {
			t.Fatalf("decode null pointer = %#v, %v", dst, err)
		}
	})
}

func TestDecodeStructAndWrappedStrings(t *testing.T) {
	type record struct {
		Bool  bool    `php:"bool,string"`
		Int   int     `php:"int,string"`
		Uint  uint    `php:"uint,string"`
		Float float64 `php:"float,string"`
		Ptr   *int    `php:"ptr,string"`
		Skip  string  `php:"-"`
	}
	raw := `a:5:{s:4:"bool";s:4:"true";s:3:"int";s:2:"-2";s:4:"uint";s:1:"3";s:5:"float";s:3:"4.5";s:3:"ptr";s:1:"5";}`
	var dst record
	if _, err := decodeTestValue(raw, &dst); err != nil {
		t.Fatalf("decode struct returned an error: %v", err)
	}
	if !dst.Bool || dst.Int != -2 || dst.Uint != 3 || dst.Float != 4.5 || dst.Ptr == nil || *dst.Ptr != 5 {
		t.Fatalf("decode struct = %#v", dst)
	}

	unknown := `a:2:{s:7:"unknown";a:1:{i:0;d:1.5;}s:3:"int";s:1:"7";}`
	if _, err := decodeTestValue(unknown, &dst); err != nil || dst.Int != 7 {
		t.Fatalf("decode struct with unknown field = %#v, %v", dst, err)
	}
}

func TestDecodeEmptyInterfaceValues(t *testing.T) {
	tests := []string{
		"N;", "b:1;", "i:1;", "d:1.5;", `s:1:"x";`,
		`a:1:{s:1:"x";i:1;}`,
		`O:8:"stdClass":1:{s:1:"x";i:1;}`,
	}
	for _, raw := range tests {
		var dst any
		if end, err := decodeTestValue(raw, &dst); err != nil || end != int64(len(raw)) {
			t.Errorf("decode interface %q ended at %d: %v", raw, end, err)
		}
	}
}

type testPHPValue struct {
	raw string
}

type testMethodInterface interface {
	TestMethod()
}

type testMethodPHPValue struct {
	raw string
}

func (*testMethodPHPValue) TestMethod() {}

func (v *testMethodPHPValue) UnmarshalPHP(raw []byte) error {
	v.raw = string(raw)
	return nil
}

type testMethodValue struct{}

func (*testMethodValue) TestMethod() {}

func (v *testPHPValue) UnmarshalPHP(raw []byte) error {
	if string(raw) == "bad" {
		return errors.New("bad value")
	}
	v.raw = string(raw)
	return nil
}

func TestDecodeUnmarshaler(t *testing.T) {
	for _, raw := range []string{`s:1:"x";`, "i:1;", `a:1:{i:0;d:1.5;}`, `O:8:"stdClass":0:{}`} {
		var dst testPHPValue
		if _, err := decodeTestValue(raw, &dst); err != nil || dst.raw != raw {
			t.Errorf("decode Unmarshaler %q = %q, %v", raw, dst.raw, err)
		}
	}
	var dst testPHPValue
	if _, err := decodeTestValue("N;", &dst); err != nil || dst != (testPHPValue{}) {
		t.Fatalf("decode null Unmarshaler = %#v, %v", dst, err)
	}

	var iface testMethodInterface = &testMethodPHPValue{}
	if _, err := decodeTestValue("i:3;", &iface); err != nil {
		t.Fatalf("decode interface Unmarshaler returned an error: %v", err)
	}
	if got := iface.(*testMethodPHPValue).raw; got != "i:3;" {
		t.Fatalf("interface Unmarshaler received %q", got)
	}
	if _, err := decodeTestValue("N;", &iface); err != nil || iface != nil {
		t.Fatalf("decode null interface = %#v, %v", iface, err)
	}

	iface = &testMethodValue{}
	if _, err := decodeTestValue("i:3;", &iface); err == nil {
		t.Fatal("decode into non-Unmarshaler interface returned nil error")
	}
	if _, err := decodeTestValue("N;", &iface); err != nil || iface != nil {
		t.Fatalf("decode null non-Unmarshaler interface = %#v, %v", iface, err)
	}
}

func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		raw string
		dst any
	}{
		{raw: "x", dst: new(bool)},
		{raw: "i:128;", dst: new(int8)},
		{raw: "i:-1;", dst: new(uint)},
		{raw: "i:256;", dst: new(uint8)},
		{raw: "d:1e100;", dst: new(float32)},
		{raw: "b:1;", dst: new(string)},
		{raw: "i:1;", dst: new(string)},
		{raw: "d:1;", dst: new(string)},
		{raw: "1", dst: new([]string)},
		{raw: "x", dst: new(map[string]int)},
		{raw: "x", dst: new(struct{ X int })},
		{raw: "x", dst: new(any)},
		{raw: `a:1:{i:0;}`, dst: new([]string)},
		{raw: `a:1:{s:1:"x";}`, dst: new(map[string]int)},
		{raw: "i:1;", dst: new(chan int)},
	}
	for _, tt := range tests {
		if _, err := decodeTestValue(tt.raw, tt.dst); err == nil {
			t.Errorf("decode %q into %T returned nil error", tt.raw, tt.dst)
		}
	}

	if _, err := CompileToGetDecoder(reflect.TypeOf((***int)(nil))); err == nil {
		t.Fatal("CompileToGetDecoder accepted a nested pointer")
	}
	type Embedded struct{ X int }
	if _, err := CompileToGetDecoder(reflect.TypeOf(&struct{ Embedded }{})); err == nil {
		t.Fatal("CompileToGetDecoder accepted an anonymous struct")
	}
	type duplicate struct {
		A int `php:"x"`
		B int `php:"x"`
	}
	if _, err := CompileToGetDecoder(reflect.TypeOf(&duplicate{})); err == nil {
		t.Fatal("CompileToGetDecoder accepted duplicate field names")
	}
}

func TestDecodeNullAndEmptyContainers(t *testing.T) {
	nullDestinations := []any{
		new(bool), new(int), new(uint), new(float64), new(string),
		new([]byte), new([]int), new([1]int), new(map[string]int),
		new(struct{ X int }),
	}
	for _, dst := range nullDestinations {
		if end, err := decodeTestValue("N;", dst); err != nil || end != 2 {
			t.Errorf("decode null into %T ended at %d: %v", dst, end, err)
		}
	}

	emptyDestinations := []any{
		new([]byte), new([]int), new([1]int), new(map[string]int), new(struct{ X int }),
	}
	for _, dst := range emptyDestinations {
		if end, err := decodeTestValue("a:0:{}", dst); err != nil || end != 6 {
			t.Errorf("decode empty array into %T ended at %d: %v", dst, end, err)
		}
	}
}

func TestDecodeDepthLimit(t *testing.T) {
	destinations := []any{new([]int), new([1]int), new(map[string]int), new(struct{ X int })}
	for _, dst := range destinations {
		dec, err := CompileToGetDecoder(reflect.TypeOf(dst))
		if err != nil {
			t.Fatal(err)
		}
		ctx := &RuntimeContext{Buf: []byte("a:0:{}")}
		if _, err := dec.Decode(ctx, 0, maxDecodeNestingDepth, reflect.ValueOf(dst).Elem()); err == nil {
			t.Errorf("decode into %T accepted excessive depth", dst)
		}
	}
}

type syntaxErrorPHPValue struct{}

func (*syntaxErrorPHPValue) UnmarshalPHP([]byte) error {
	return internalerrors.ErrSyntax("custom syntax error", 0)
}

func TestUnmarshalerErrorAnnotation(t *testing.T) {
	var dst syntaxErrorPHPValue
	if _, err := decodeTestValue("i:1;", &dst); err == nil {
		t.Fatal("decode custom syntax error returned nil")
	}
}

func TestUnquoteBytes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: `""`},
		{input: `"plain"`, want: "plain"},
		{input: `"日本語"`, want: "日本語"},
		{input: `"a\nb"`, want: "a\nb"},
		{input: `"\"\\\/\'"`, want: `"\/'`},
		{input: `"\b\f\r\t"`, want: "\b\f\r\t"},
		{input: `"\u0041"`, want: "A"},
		{input: `"\uD834\uDD1E"`, want: "𝄞"},
		{input: `"\uD800x"`, want: "�x"},
	}
	for _, tt := range tests {
		got, ok := unquoteBytes([]byte(tt.input))
		if !ok || string(got) != tt.want {
			t.Errorf("unquoteBytes(%q) = %q, %v; want %q", tt.input, got, ok, tt.want)
		}
	}

	invalid := [][]byte{
		nil,
		[]byte("x"),
		[]byte(`"\x"`),
		[]byte{'"', '\\', '"'},
		[]byte{'"', 'a', '"', 'b', '"'},
		[]byte{'"', '\n', '"'},
	}
	for _, input := range invalid {
		if _, ok := unquoteBytes(input); ok {
			t.Errorf("unquoteBytes(%q) unexpectedly succeeded", input)
		}
	}
	if got, ok := unquoteBytes([]byte{'"', 0xff, '"'}); !ok || string(got) != "�" {
		t.Errorf("unquoteBytes malformed UTF-8 = %q, %v", got, ok)
	}

	for _, input := range [][]byte{[]byte(`\u12`), []byte(`\uZZZZ`), []byte(`x`)} {
		if got := getu4(input); got >= 0 {
			t.Errorf("getu4(%q) = %d", input, got)
		}
	}
}
