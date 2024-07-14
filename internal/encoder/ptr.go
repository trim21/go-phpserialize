package encoder

import (
	"fmt"
	"reflect"
)

func compilePtr(rt reflect.Type, seen seenMap) (encoder, error) {
	switch rt.Elem().Kind() {
	case reflect.Ptr:
		return nil, fmt.Errorf("encoding nested ptr is not supported *%s", rt.Elem().String())
	case reflect.Bool:
		return deRefNilEncoder(encodeBool), nil
	case reflect.Uint8:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Uint16:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Uint32:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Uint64:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Uint:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Int8:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Int16:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Int32:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Int64:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Int:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Float32:
		return deRefNilEncoder(encodeFloat32), nil
	case reflect.Float64:
		return deRefNilEncoder(encodeFloat64), nil
	case reflect.String:
		return deRefNilEncoder(encodeString), nil
	case reflect.Interface:
		return compileInterface(rt.Elem())
	case reflect.Map:
		enc, err := compileMap(rt.Elem(), seen)
		return deRefNilEncoder(enc), err
	case reflect.Struct:
		enc, err := compileStruct(rt.Elem(), seen)
		return deRefNilEncoder(enc), err
	}

	enc, err := compile(rt.Elem(), seen)
	if err != nil {
		return nil, err
	}

	return deRefNilEncoder(enc), nil
}

func deRefNilEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		return enc(ctx, b, rv.Elem())
	}
}
