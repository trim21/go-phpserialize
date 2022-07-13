package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileArray(rt *runtime.Type) (encoder, error) {
	offset := rt.Elem().Size()
	length := rt.Len()
	i64 := int64(length)
	var enc encoder
	var err error

	if rt.Elem().Kind() == reflect.Map {
		enc, err = compileWithCache(runtime.PtrTo(rt.Elem()))
		if err != nil {
			return nil, err
		}
	} else {
		enc, err = compileWithCache(rt.Elem())
		if err != nil {
			return nil, err
		}
	}

	return func(ctx *Ctx, b []byte, dataPtr uintptr) ([]byte, error) {
		// no data ptr, nil slice
		if dataPtr == 0 {
			return appendNull(b), nil
		}

		b = appendArrayBegin(b, i64)
		var err error // create a new error value, so shadow compiler's error
		for i := 0; i < length; i++ {
			b = appendIntBytes(b, int64(i))
			b, err = enc(ctx, b, dataPtr+offset*uintptr(i))
			if err != nil {
				return b, err
			}
		}
		return append(b, '}'), nil
	}, nil
}
