package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

const lenOffset = unsafe.Offsetof(reflect.SliceHeader{}.Len)

func compileSlice(rt reflect.Type) (encoder, error) {
	elType := rt.Elem()

	offset := rt.Elem().Size()

	encoder, err := compile(elType)
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, p uintptr) error {
		dataPtr := *(*uintptr)(unsafe.Pointer(p))

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
			err = encoder(ctx, dataPtr+offset*uintptr(i))
			if err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}
