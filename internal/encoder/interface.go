package encoder

import (
	"fmt"
	"reflect"
)

// will need to get type message at marshal time, slow path.
// should avoid interface for performance thinking.
func compileInterface(rt reflect.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		return reflectInterfaceValue(ctx, b, rv)
	}, nil
}

func reflectInterfaceValue(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				return appendNull(b), nil
			}
			rv = rv.Elem()
		default:
			break LOOP
		}
	}

	// simple type
	switch rv.Type().Kind() {
	case reflect.Bool:
		return encodeBool(ctx, b, rv)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return encodeUint(ctx, b, rv)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return encodeInt(ctx, b, rv)
	case reflect.Float32:
		return encodeFloat32(ctx, b, rv)
	case reflect.Float64:
		return encodeFloat64(ctx, b, rv)
	case reflect.String:
		return encodeString(ctx, b, rv)
	}

	switch rv.Kind() {
	case reflect.Slice:
		return reflectSlice(ctx, b, rv)
	case reflect.Map:
		return reflectMap(ctx, b, rv)
	case reflect.Struct:
		return reflectStruct(ctx, b, rv)
	}

	return b, &UnsupportedInterfaceTypeError{rv.Type()}
}

func compileInterfaceAsString(rt reflect.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		return reflectInterfaceValueAsString(ctx, b, rv)
	}, nil
}

func reflectInterfaceValueAsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				return appendNull(b), nil
			}
			rv = rv.Elem()
		default:
			break LOOP
		}
	}

	switch value := rv.Interface().(type) {
	case bool:
		return appendBoolAsString(b, value)
	case uint8:
		return appendUintAsString(b, uint64(value))
	case uint16:
		return appendUintAsString(b, uint64(value))
	case uint32:
		return appendUintAsString(b, uint64(value))
	case uint64:
		return appendUintAsString(b, uint64(value))
	case uint:
		return appendUintAsString(b, uint64(value))
	case int8:
		return appendIntAsString(b, int64(value))
	case int16:
		return appendIntAsString(b, int64(value))
	case int32:
		return appendIntAsString(b, int64(value))
	case int64:
		return appendIntAsString(b, int64(value))
	case int:
		return appendIntAsString(b, int64(value))
	case float32:
		return appendFloat32AsString(ctx.smallBuffer[:0], b, value), nil
	case float64:
		return appendFloat64AsString(ctx.smallBuffer[:0], b, value), nil
	}

	// slice, map and struct as interface are not supported yet.
	return b, fmt.Errorf("failed to encode %s as string", rv.Kind())
}
