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

func encodeFloat32(buf *Ctx, p uintptr) error {
	value := *(*float32)(unsafe.Pointer(p))
	appendFloat32(buf, value)
	return nil
}

func encodeFloat64(ctx *Ctx, p uintptr) error {
	f64 := *(*float64)(unsafe.Pointer(p))
	appendFloat64(ctx, f64)
	return nil
}

// https://github.com/goccy/go-json/blob/4d0a50640b999aeafd15e3b20d8ad47fe917e6e8/internal/encoder/encoder.go#L335

func appendFloat32(ctx *Ctx, f32 float32) {
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

	ctx.b = append(ctx.b, 'd', ':')
	ctx.b = strconv.AppendFloat(ctx.b, f64, format, -1, 32)
	ctx.b = append(ctx.b, ';')
}

func appendFloat64(ctx *Ctx, f64 float64) {
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			format = 'e'
		}
	}

	ctx.b = append(ctx.b, 'd', ':')
	ctx.b = strconv.AppendFloat(ctx.b, f64, format, -1, 64)
	ctx.b = append(ctx.b, ';')
}

func compileFloatAsString(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.Float32:
		return encodeFloat32AsString, nil
	case reflect.Float64:
		return encodeFloat64AsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

func encodeFloat32AsString(buf *Ctx, p uintptr) error {
	value := *(*float32)(unsafe.Pointer(p))
	appendFloat32AsString(buf, value)
	return nil
}

func encodeFloat64AsString(ctx *Ctx, p uintptr) error {
	f64 := *(*float64)(unsafe.Pointer(p))
	appendFloat64AsString(ctx, f64)
	return nil
}

func appendFloat32AsString(ctx *Ctx, f32 float32) {
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

	ctx.floatBuffer = strconv.AppendFloat(ctx.floatBuffer, f64, format, -1, 32)
	appendStringHead(ctx, int64(len(ctx.floatBuffer)))
	ctx.b = append(ctx.b, '"')
	ctx.b = append(ctx.b, ctx.floatBuffer...)
	ctx.b = append(ctx.b, '"', ';')
	ctx.floatBuffer = ctx.floatBuffer[:]
}

func appendFloat64AsString(ctx *Ctx, f64 float64) {
	abs := math.Abs(f64)
	format := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			format = 'e'
		}
	}

	ctx.floatBuffer = strconv.AppendFloat(ctx.floatBuffer, f64, format, -1, 64)
	appendStringHead(ctx, int64(len(ctx.floatBuffer)))
	ctx.b = append(ctx.b, '"')
	ctx.b = append(ctx.b, ctx.floatBuffer...)
	ctx.b = append(ctx.b, '"', ';')
	ctx.floatBuffer = ctx.floatBuffer[:]
}
