package encoder

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func encodeFloat32(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**float32)(unsafe.Pointer(&p))
	return appendFloat32(b, value), nil
}

func encodeFloat64(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**float64)(unsafe.Pointer(&p))
	return appendFloat64(b, value), nil
}

// https://github.com/goccy/go-json/blob/4d0a50640b999aeafd15e3b20d8ad47fe917e6e8/internal/encoder/encoder.go#L335

func appendFloat32(b []byte, f32 float32) []byte {
	f64 := float64(f32)
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		f32 := float32(abs)
		if f32 < 1e-6 || f32 >= 1e21 {
			format = 'e'
		}
	}

	b = append(b, 'd', ':')
	b = strconv.AppendFloat(b, f64, format, -1, 32)
	return append(b, ';')
}

func appendFloat64(b []byte, f64 float64) []byte {
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			format = 'e'
		}
	}

	b = append(b, 'd', ':')
	b = strconv.AppendFloat(b, f64, format, -1, 64)
	return append(b, ';')
}

func compileFloatAsString(typ *runtime.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.Float32:
		return encodeFloat32AsString, nil
	case reflect.Float64:
		return encodeFloat64AsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

func encodeFloat32AsString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**float32)(unsafe.Pointer(&p))
	return appendFloat32AsString(ctx.smallBuffer[:0], b, value), nil
}

func encodeFloat64AsString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	f64 := **(**float64)(unsafe.Pointer(&p))
	return appendFloat64AsString(ctx.smallBuffer[:0], b, f64), nil
}

func appendFloat32AsString(buf []byte, b []byte, f32 float32) []byte {
	f64 := float64(f32)
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		f32 := float32(abs)
		if f32 < 1e-6 || f32 >= 1e21 {
			format = 'e'
		}
	}

	buf = strconv.AppendFloat(buf, f64, format, -1, 32)
	b = appendStringHead(b, int64(len(buf)))
	b = append(b, '"')
	b = append(b, buf...)
	b = append(b, '"', ';')

	return b
}

func appendFloat64AsString(buf []byte, b []byte, f64 float64) []byte {
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			format = 'e'
		}
	}

	buf = strconv.AppendFloat(buf, f64, format, -1, 64)
	b = appendStringHead(b, int64(len(buf)))
	b = append(b, '"')
	b = append(b, buf...)
	b = append(b, '"', ';')

	return b
}
