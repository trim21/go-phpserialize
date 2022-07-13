package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// MUST call `compilePtr` directly when compile encoder for struct field.
func compilePtr(rt *runtime.Type, indirect bool) (encoder, error) {
	switch rt.Elem().Kind() {
	case reflect.Bool:
		return wrapNilEncoder(encodeBool), nil
	case reflect.Uint8:
		return wrapNilEncoder(encodeUint8), nil
	case reflect.Uint16:
		return wrapNilEncoder(encodeUint16), nil
	case reflect.Uint32:
		return wrapNilEncoder(encodeUint32), nil
	case reflect.Uint64:
		return wrapNilEncoder(encodeUint64), nil
	case reflect.Uint:
		return wrapNilEncoder(encodeUint), nil
	case reflect.Int8:
		return wrapNilEncoder(encodeInt8), nil
	case reflect.Int16:
		return wrapNilEncoder(encodeInt16), nil
	case reflect.Int32:
		return wrapNilEncoder(encodeInt32), nil
	case reflect.Int64:
		return wrapNilEncoder(encodeInt64), nil
	case reflect.Int:
		return wrapNilEncoder(encodeInt), nil
	case reflect.Float32:
		return wrapNilEncoder(encodeFloat32), nil
	case reflect.Float64:
		return wrapNilEncoder(encodeFloat64), nil
	case reflect.String:
		if indirect {
			return EncodeStringPtr, nil
		}
		return encodeString, nil
	case reflect.Interface:
		return compileInterface(rt.Elem())
	case reflect.Struct:
		enc, err := compileStruct(rt.Elem())
		return wrapNilEncoder(enc), err
	}

	enc, err := compile(rt.Elem())
	if err != nil {
		return nil, err
	}

	if !indirect {
		return enc, nil
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

func compilePtrAsString(rt *runtime.Type) (encoder, error) {
	inner, err := compileAsString(rt.Elem())
	if err != nil {
		return nil, err
	}
	return deRefNilEncoder(inner), nil
}

func onlyDeReferEncoder(enc encoder) encoder {
	if enc == nil {
		return nil
	}

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		return enc(ctx, b, PtrDeRef(p))
	}
}
