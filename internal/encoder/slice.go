package encoder

import (
	"reflect"
)

func compileSlice(rt reflect.Type, seen seenMap) (encoder, error) {
	var enc encoder
	var err error

	if rt.Elem().Kind() == reflect.Map {
		enc, err = compile(reflect.PointerTo(rt.Elem()), seen)
		if err != nil {
			return nil, err
		}
	} else {
		enc, err = compile(rt.Elem(), seen)
		if err != nil {
			return nil, err
		}
	}

	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		length := rv.Len()

		b = appendArrayBegin(b, int64(length))
		var err error // create a new error value, so shadow compiler's error
		for i := 0; i < length; i++ {
			b = appendIntBytes(b, int64(i))
			b, err = enc(ctx, b, rv.Index(i))
			if err != nil {
				return b, err
			}
		}
		return append(b, '}'), nil
	}, nil
}
