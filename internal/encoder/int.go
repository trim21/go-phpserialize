package encoder

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileInt(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
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
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

func encodeInt8(buf *Ctx, p uintptr) error {
	value := *(*int8)(unsafe.Pointer(p))
	appendInt(buf, int64(value))
	return nil
}

func encodeInt16(buf *Ctx, p uintptr) error {
	value := *(*int16)(unsafe.Pointer(p))
	appendInt(buf, int64(value))
	return nil
}

func encodeInt32(buf *Ctx, p uintptr) error {
	value := *(*int32)(unsafe.Pointer(p))
	appendInt(buf, int64(value))
	return nil
}

func encodeInt64(buf *Ctx, p uintptr) error {
	value := *(*int64)(unsafe.Pointer(p))
	appendInt(buf, int64(value))
	return nil
}

func encodeInt(buf *Ctx, p uintptr) error {
	value := *(*int)(unsafe.Pointer(p))
	appendInt(buf, int64(value))
	return nil
}

func appendInt(buf *Ctx, v int64) {
	buf.b = append(buf.b, 'i', ':')
	buf.b = strconv.AppendInt(buf.b, v, 10)
	buf.b = append(buf.b, ';')
}
