package encoder

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestMarshalScalarTypes(t *testing.T) {
	tests := []struct {
		value any
		want  string
	}{
		{value: nil, want: "N;"},
		{value: true, want: "b:1;"},
		{value: false, want: "b:0;"},
		{value: int8(-8), want: "i:-8;"},
		{value: int16(-16), want: "i:-16;"},
		{value: int32(-32), want: "i:-32;"},
		{value: int64(-64), want: "i:-64;"},
		{value: int(-1), want: "i:-1;"},
		{value: uint8(8), want: "i:8;"},
		{value: uint16(16), want: "i:16;"},
		{value: uint32(32), want: "i:32;"},
		{value: uint64(64), want: "i:64;"},
		{value: uint(1), want: "i:1;"},
		{value: float32(3.5), want: "d:3.5;"},
		{value: float64(-4.5), want: "d:-4.5;"},
		{value: "foo", want: `s:3:"foo";`},
		{value: []byte("foo"), want: `s:3:"foo";`},
	}
	for _, tt := range tests {
		got, err := Marshal(tt.value)
		if err != nil {
			t.Errorf("Marshal(%T) returned an error: %v", tt.value, err)
		} else if string(got) != tt.want {
			t.Errorf("Marshal(%#v) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestMarshalFloatFormats(t *testing.T) {
	tests := []struct {
		value any
		want  string
	}{
		{value: math.Inf(1), want: "d:INF;"},
		{value: math.Inf(-1), want: "d:-INF;"},
		{value: math.NaN(), want: "d:NAN;"},
		{value: float32(1e-7), want: "d:1e-07;"},
		{value: 1e-7, want: "d:1e-07;"},
		{value: 1e21, want: "d:1e+21;"},
		{value: 0.0, want: "d:0;"},
	}
	for _, tt := range tests {
		got, err := Marshal(tt.value)
		if err != nil || string(got) != tt.want {
			t.Errorf("Marshal(%v) = %q, %v; want %q", tt.value, got, err, tt.want)
		}
	}
}

func TestMarshalCompositeTypes(t *testing.T) {
	tests := []any{
		[2]int{1, 2},
		[]string{"a", "b"},
		[]map[string]int{{"x": 1}},
		map[string]int{"x": 1},
		map[int]string{1: "x"},
		map[uint]int{1: 2},
		map[string]map[string]int{"outer": {"inner": 1}},
	}
	for _, value := range tests {
		got, err := Marshal(value)
		if err != nil {
			t.Errorf("Marshal(%T) returned an error: %v", value, err)
		} else if !strings.HasPrefix(string(got), "a:") {
			t.Errorf("Marshal(%T) = %q", value, got)
		}
	}

	var nilSlice []int
	var nilMap map[string]int
	for _, value := range []any{nilSlice, nilMap, (*int)(nil)} {
		got, err := Marshal(value)
		if err != nil || string(got) != "N;" {
			t.Errorf("Marshal(%T(nil)) = %q, %v", value, got, err)
		}
	}
}

func TestMarshalPointers(t *testing.T) {
	b := true
	i := 1
	u := uint(2)
	f32 := float32(3)
	f64 := 4.0
	s := "x"
	m := map[string]int{"x": 1}
	slice := []int{1}
	for _, value := range []any{&b, &i, &u, &f32, &f64, &s, &m, &slice} {
		if _, err := Marshal(value); err != nil {
			t.Errorf("Marshal(%T) returned an error: %v", value, err)
		}
	}
}

func TestMarshalStructTags(t *testing.T) {
	type record struct {
		Bool    bool    `php:"bool,string"`
		Int8    int8    `php:"int8,string"`
		Int16   int16   `php:"int16,string"`
		Int32   int32   `php:"int32,string"`
		Int64   int64   `php:"int64,string"`
		Int     int     `php:"int,string"`
		Uint    uint    `php:"uint,string"`
		Float32 float32 `php:"float32,string"`
		Float64 float64 `php:"float64,string"`
		Any     any     `php:"any,string"`
		Ptr     *int    `php:"ptr,omitempty"`
		Omit    string  `php:"omit,omitempty"`
		Ignore  string  `php:"-"`
		private string
	}
	value := record{
		Bool: true, Int8: -8, Int16: -16, Int32: -32, Int64: -64, Int: -1,
		Uint: 2, Float32: 3.5, Float64: 4.5, Any: int64(5),
	}
	got, err := Marshal(value)
	if err != nil {
		t.Fatalf("Marshal struct returned an error: %v", err)
	}
	for _, part := range []string{`s:4:"true";`, `s:2:"-8";`, `s:3:"3.5";`, `s:1:"5";`} {
		if !strings.Contains(string(got), part) {
			t.Errorf("Marshal struct = %q, missing %q", got, part)
		}
	}
}

func TestMarshalInterfaceValues(t *testing.T) {
	type holder struct {
		Value any `php:"value"`
	}
	values := []any{
		nil, true,
		uint8(1), uint16(2), uint32(3), uint64(4), uint(5),
		int8(-1), int16(-2), int32(-3), int64(-4), int(-5),
		float32(1.5), float64(2.5), "x",
		[]any{1, "x"}, []int{1, 2}, []int{},
		map[string]any{"x": 1}, map[string]int{"x": 1}, map[string]int{},
		struct{ X int }{X: 1},
	}
	for _, value := range values {
		if _, err := Marshal(holder{Value: value}); err != nil {
			t.Errorf("Marshal interface value %T returned an error: %v", value, err)
		}
	}

	var ptr any = new(int)
	if _, err := Marshal(holder{Value: &ptr}); err != nil {
		t.Errorf("Marshal nested interface pointer returned an error: %v", err)
	}
}

func TestMarshalInterfaceValuesAsString(t *testing.T) {
	type holder struct {
		Value any `php:"value,string"`
	}
	values := []any{
		true,
		uint8(1), uint16(2), uint32(3), uint64(4), uint(5),
		int8(-1), int16(-2), int32(-3), int64(-4), int(-5),
		float32(1.5), float64(2.5),
	}
	for _, value := range values {
		got, err := Marshal(holder{Value: value})
		if err != nil {
			t.Errorf("Marshal interface value %T as string returned an error: %v", value, err)
		} else if !strings.Contains(string(got), `s:`) {
			t.Errorf("Marshal interface value %T as string = %q", value, got)
		}
	}
	if _, err := Marshal(holder{Value: []int{1}}); err == nil {
		t.Fatal("Marshal accepted an interface slice with string tag")
	}
}

type testMarshaler struct {
	err bool
}

func (m testMarshaler) MarshalPHP() ([]byte, error) {
	if m.err {
		return nil, errors.New("marshal failure")
	}
	return []byte("i:7;"), nil
}

func TestMarshalCustomMarshaler(t *testing.T) {
	got, err := Marshal(testMarshaler{})
	if err != nil || string(got) != "i:7;" {
		t.Fatalf("Marshal custom value = %q, %v", got, err)
	}
	if _, err := Marshal(testMarshaler{err: true}); err == nil {
		t.Fatal("Marshal accepted custom marshaler error")
	}
}

func TestMarshalErrors(t *testing.T) {
	badValues := []any{
		make(chan int),
		func() {},
		map[bool]int{true: 1},
		struct{ Value any }{Value: make(chan int)},
	}
	for _, value := range badValues {
		if _, err := Marshal(value); err == nil {
			t.Errorf("Marshal(%T) returned nil error", value)
		}
	}

	i := 1
	p := &i
	if _, err := Marshal(&p); err == nil {
		t.Fatal("Marshal accepted nested pointer")
	}

	type badStringTag struct {
		Value []int `php:"value,string"`
	}
	if _, err := Marshal(badStringTag{}); err == nil {
		t.Fatal("Marshal accepted string tag on slice")
	}
}

func TestEncoderHelpers(t *testing.T) {
	for _, tt := range []struct {
		value int64
		want  int64
	}{{0, 1}, {-1, 2}, {9, 1}, {10, 2}, {-10, 3}} {
		if got := iterativeDigitsCount(tt.value); got != tt.want {
			t.Errorf("iterativeDigitsCount(%d) = %d, want %d", tt.value, got, tt.want)
		}
	}
	for _, tt := range []struct {
		value uint64
		want  int64
	}{{0, 1}, {9, 1}, {10, 2}} {
		if got := uintDigitsCount(tt.value); got != tt.want {
			t.Errorf("uintDigitsCount(%d) = %d, want %d", tt.value, got, tt.want)
		}
	}

	for _, typ := range []reflect.Type{reflect.TypeOf(true), reflect.TypeOf(int(0)), reflect.TypeOf(uint(0)), reflect.TypeOf(float64(0)), reflect.TypeOf("")} {
		if _, err := compileAsString(typ); err != nil {
			t.Errorf("compileAsString(%v) returned an error: %v", typ, err)
		}
	}

	for _, err := range []error{
		&UnsupportedTypeError{Type: reflect.TypeOf(make(chan int))},
		&UnsupportedTypeAsMapKeyError{Type: reflect.TypeOf(true)},
		&UnsupportedInterfaceTypeError{Type: reflect.TypeOf(make(chan int))},
	} {
		if strings.TrimSpace(err.Error()) == "" {
			t.Errorf("%T has an empty error message", err)
		}
	}
}
