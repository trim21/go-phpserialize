package encoder

import (
	"reflect"
	"unsafe"
)

const lenOffset = unsafe.Offsetof(reflect.SliceHeader{}.Len)

func compileSlice(rt reflect.Type, seen seenMap) (encoder, error) {
	offset := rt.Elem().Size()
	var enc encoder
	var err error

	if rt.Elem().Kind() == reflect.Map {
		enc, err = compile(rt.Elem(), seen)
		if err != nil {
			return nil, err
		}
	} else {
		enc, err = compile(rt.Elem(), seen)
		if err != nil {
			return nil, err
		}
	}
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		if p == 0 {
			return appendNull(b), nil
		}
		dataPtr := **(**uintptr)(unsafe.Pointer(&p))
		// no data ptr, nil slice
		if dataPtr == 0 {
			return appendNull(b), nil
		}

		length := *(*int)(unsafe.Add(ptrToUnsafePtr(p), lenOffset))

		b = appendArrayBegin(b, int64(length))
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
