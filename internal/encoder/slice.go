package encoder

import (
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
)

const lenOffset = unsafe.Offsetof(reflect.SliceHeader{}.Len)

func compileSlice(rt reflect.Type, rv reflect.Value) (encoder, error) {
	elType := rt.Elem()

	offset := rt.Elem().Size()

	encoder, err := compile(elType, reflect.New(elType).Elem())
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, p uintptr) error {
		dataPtr := *(*uintptr)(unsafe.Pointer(p))
		length := *(*int)(unsafe.Pointer(p + lenOffset))

		appendArrayBegin(ctx, int64(length))

		for i := 0; i < length; i++ {
			ctx.b = append(ctx.b, 'i', ':')
			ctx.b = strconv.AppendInt(ctx.b, int64(i), 10)
			ctx.b = append(ctx.b, ';')
			err = encoder(ctx, dataPtr+offset*uintptr(i))
			if err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}
