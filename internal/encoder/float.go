package encoder

import (
	"fmt"
	"math"
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileFloat(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

// https://github.com/goccy/go-json/blob/4d0a50640b999aeafd15e3b20d8ad47fe917e6e8/internal/encoder/encoder.go#L335

func encodeFloat32(buf *Ctx, p uintptr) error {
	value := *(*float32)(unsafe.Pointer(p))
	f64 := float64(value)
	abs := math.Abs(f64)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		f32 := float32(abs)
		if f32 < 1e-6 || f32 >= 1e21 {
			fmt = 'e'
		}
	}

	buf.b = append(buf.b, 'd', ':')
	buf.b = strconv.AppendFloat(buf.b, f64, fmt, -1, 32)
	buf.b = append(buf.b, ';')
	return nil
}

func encodeFloat64(buf *Ctx, p uintptr) error {
	f64 := *(*float64)(unsafe.Pointer(p))

	abs := math.Abs(f64)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}

	buf.b = append(buf.b, 'd', ':')
	buf.b = strconv.AppendFloat(buf.b, f64, fmt, -1, 64)
	buf.b = append(buf.b, ';')
	return nil
}
