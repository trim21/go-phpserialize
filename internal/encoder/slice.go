package encoder

import (
	"reflect"
)

func compileSlice(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	var enc encoder
	var compileError error

	if rt.Elem().Kind() == reflect.Map {
		enc, compileError = compile(reflect.PointerTo(rt.Elem()), seen)
		if compileError != nil {
			return nil, compileError
		}
	} else {
		enc, compileError = compile(rt.Elem(), seen)
		if compileError != nil {
			return nil, compileError
		}
	}

	return checkRecursiveEncoder(func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		length := rv.Len()

		b = appendArrayBegin(b, int64(length))
		var err error
		for i := 0; i < length; i++ {
			b = appendIntBytes(b, int64(i))
			b, err = enc(ctx, b, rv.Index(i))
			if err != nil {
				return b, err
			}
		}
		return append(b, '}'), nil
	}), nil
}
