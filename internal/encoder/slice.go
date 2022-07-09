package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

const lenOffset = unsafe.Offsetof(reflect.SliceHeader{}.Len)

func compileSlice(rt *runtime.Type) (encoder, error) {
	offset := rt.Elem().Size()
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
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		dataPtr := *(*uintptr)(unsafe.Pointer(p))
		// no data ptr, nil slice
		if dataPtr == 0 {
			return appendNilBytes(b), nil
		}

		length := *(*int)(unsafe.Pointer(p + lenOffset))

		b = appendArrayBeginBytes(b, int64(length))
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
