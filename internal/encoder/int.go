package encoder

import (
	"strconv"
	"unsafe"
)

func encodeInt8(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := *(*int8)(unsafe.Pointer(p))
	return appendIntBytes(b, int64(value)), nil
}

func encodeInt16(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := *(*int16)(unsafe.Pointer(p))
	return appendIntBytes(b, int64(value)), nil

}

func encodeInt32(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := *(*int32)(unsafe.Pointer(p))
	return appendIntBytes(b, int64(value)), nil

}

func encodeInt64(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**int64)(unsafe.Pointer(&p))
	return appendIntBytes(b, int64(value)), nil
}

func encodeInt(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := *(*int)(unsafe.Pointer(p))
	return appendIntBytes(b, int64(value)), nil
}

func appendIntBytes(b []byte, v int64) []byte {
	b = append(b, 'i', ':')
	b = strconv.AppendInt(b, v, 10)

	return append(b, ';')
}
