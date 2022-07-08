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
		enc, err = compile(runtime.PtrTo(rt.Elem()))
		if err != nil {
			return nil, err
		}
	} else {
		enc, err = compile(rt.Elem())
		if err != nil {
			return nil, err
		}
	}
	return func(ctx *Ctx, p uintptr) error {
		dataPtr := *(*uintptr)(unsafe.Pointer(p))
		// fmt.Println(unsafe.Pointer(p))

		// no data ptr, nil slice
		// even empty slice has a non-zero data ptr
		if dataPtr == 0 {
			appendNil(ctx)
			return nil
		}

		length := *(*int)(unsafe.Pointer(p + lenOffset))

		appendArrayBegin(ctx, int64(length))

		for i := 0; i < length; i++ {
			appendInt(ctx, int64(i))
			err = enc(ctx, dataPtr+offset*uintptr(i))
			if err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}
