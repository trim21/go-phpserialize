package encoder

import (
	"fmt"
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// MUST call `compilePtr` directly when compile encoder for struct field.
func compilePtr(rt *runtime.Type) (encoder, error) {
	switch rt.Elem().Kind() {
	case reflect.Ptr:
		return nil, fmt.Errorf("encoding nested ptr is not supported *%s", rt.Elem().String())

	case reflect.Bool:
		return encodeBool, nil
	case reflect.Uint8:
		return encodeUint8, nil
	case reflect.Uint16:
		return encodeUint16, nil
	case reflect.Uint32:
		return encodeUint32, nil
	case reflect.Uint64:
		return encodeUint64, nil
	case reflect.Uint:
		return encodeUint, nil
	case reflect.Int8:
		return encodeInt8, nil
	case reflect.Int16:
		return encodeInt16, nil
	case reflect.Int32:
		return encodeInt32, nil
	case reflect.Int64:
		return encodeInt64, nil
	case reflect.Int:
		return encodeInt, nil
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	case reflect.String:
		return encodeString, nil
	case reflect.Interface:
		return compileInterface(rt.Elem())
	case reflect.Map:
		enc, err := compileMap(rt.Elem())
		return deRefNilEncoder(enc), err
	case reflect.Struct:
		enc, err := compileStruct(rt.Elem())
		return wrapNilEncoder(enc), err
	}

	enc, err := compile(rt.Elem())
	if err != nil {
		return nil, err
	}

	return deRefNilEncoder(enc), nil
}

func deRefNilEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		if p == 0 {
			return appendNull(b), nil
		}
		p = PtrDeRef(p)
		return enc(ctx, b, p)
	}
}

func wrapNilEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		if p == 0 {
			return appendNull(b), nil
		}
		return enc(ctx, b, p)
	}
}
