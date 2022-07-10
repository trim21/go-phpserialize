package encoder

import (
	"strconv"
	"unsafe"
)

func encodeUint8(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**uint8)(unsafe.Pointer(&p))
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, uint64(value), 10)
	return append(b, ';'), nil
}

func encodeUint16(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**uint16)(unsafe.Pointer(&p))
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, uint64(value), 10)
	return append(b, ';'), nil
}

func encodeUint32(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**uint32)(unsafe.Pointer(&p))
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, uint64(value), 10)
	return append(b, ';'), nil
}

func encodeUint64(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**uint64)(unsafe.Pointer(&p))
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, uint64(value), 10)
	return append(b, ';'), nil
}

func encodeUint(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**uint)(unsafe.Pointer(&p))
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, uint64(value), 10)
	return append(b, ';'), nil
}
