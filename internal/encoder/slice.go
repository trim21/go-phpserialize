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

	return func(buf *buffer, p uintptr) error {
		dataPtr := *(*uintptr)(unsafe.Pointer(p))
		length := *(*int)(unsafe.Pointer(p + lenOffset))

		appendArrayBegin(buf, int64(length))

		for i := 0; i < length; i++ {
			buf.b = append(buf.b, 'i', ':')
			buf.b = strconv.AppendInt(buf.b, int64(i), 10)
			buf.b = append(buf.b, ';')
			err = encoder(buf, dataPtr+offset*uintptr(i))
			if err != nil {
				return err
			}
		}
		buf.b = append(buf.b, '}')
		return nil
	}, nil
}
