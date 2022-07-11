package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compilePtr(rt *runtime.Type) (encoder, error) {
	switch rt.Elem().Kind() {
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
		return EncodeStringPtr, nil
	case reflect.Interface:
		return compileInterface(rt.Elem())
	}

	enc, err := compile(rt.Elem())
	if err != nil {
		return nil, err
	}

	return deRefEncoder(enc), nil
}

func deRefEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		return enc(ctx, b, PtrDeRef(p))
	}
}
