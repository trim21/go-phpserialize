package encoder

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileUint(rt reflect.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Uint8:
		return encodeUint8, nil
	case reflect.Uint16:
		return encodeUint16, nil
	case reflect.Uint32:
		return encodeUint32, nil
	case reflect.Uint64:
		return encodeUint64, nil
	case reflect.Uint:
		return encodeUint, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", rt.Kind()))
}

func encodeUint8(buf *buffer, p uintptr) error {
	value := *(*uint8)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint16(buf *buffer, p uintptr) error {
	value := *(*uint16)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint32(buf *buffer, p uintptr) error {
	value := *(*uint32)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint64(buf *buffer, p uintptr) error {
	value := *(*uint64)(unsafe.Pointer(p))
	buf.b = append(buf.b, 'i', ':')
	buf.b = strconv.AppendUint(buf.b, value, 10)
	buf.b = append(buf.b, ';')
	return nil
}

func encodeUint(buf *buffer, p uintptr) error {
	value := *(*uint)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func appendUint(buf *buffer, v uint64) {
	buf.b = append(buf.b, 'i', ':')
	buf.b = strconv.AppendUint(buf.b, v, 10)
	buf.b = append(buf.b, ';')
}
