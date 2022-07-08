package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// TODO: why ptr work fine without an encoder?
func compilePtr(rt *runtime.Type) (encoder, error) {
	switch rt.Elem().Kind() {
	case reflect.Bool:
		return compile(rt.Elem())
	case reflect.Uint8:
		return compile(rt.Elem())
	case reflect.Uint16:
		return compile(rt.Elem())
	case reflect.Uint32:
		return compile(rt.Elem())
	case reflect.Uint64:
		return compile(rt.Elem())
	case reflect.Uint:
		return compile(rt.Elem())
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
		return encodeStringVariable, nil
	}

	enc, err := compile(rt.Elem())
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, p uintptr) error {
		return enc(ctx, ptrOfPtr(p))
	}, nil
}
