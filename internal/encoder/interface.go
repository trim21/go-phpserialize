package encoder

import (
	"fmt"
	"reflect"
	"unsafe"
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

	value := rv.Interface()

	v := *(*emptyInterface)(unsafe.Pointer(&value))
	p := uintptr(v.ptr)

	// simple type
	switch rv.Type().Kind() {
	case reflect.Bool:
		return encodeBoolAsString(ctx, b, p)
	case reflect.Uint8:
		return encodeUint8AsString(ctx, b, p)
	case reflect.Uint16:
		return encodeUint16AsString(ctx, b, p)
	case reflect.Uint32:
		return encodeUint32AsString(ctx, b, p)
	case reflect.Uint64:
		return encodeUint64AsString(ctx, b, p)
	case reflect.Uint:
		return encodeUintAsString(ctx, b, p)
	case reflect.Int8:
		return encodeInt8AsString(ctx, b, p)
	case reflect.Int16:
		return encodeInt16AsString(ctx, b, p)
	case reflect.Int32:
		return encodeInt32AsString(ctx, b, p)
	case reflect.Int64:
		return encodeInt64AsString(ctx, b, p)
	case reflect.Int:
		return encodeIntAsString(ctx, b, p)
	case reflect.Float32:
		return encodeFloat32AsString(ctx, b, p)
	case reflect.Float64:
		return encodeFloat64AsString(ctx, b, p)
	}

	// slice, map and struct as interface are not supported yet.
	return b, fmt.Errorf("failed to encode %s as string", rv.Kind())
}
