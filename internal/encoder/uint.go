package encoder

import (
	"strconv"
	"unsafe"
)

func encodeUint8(buf *Ctx, p uintptr) error {
	value := *(*uint8)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint16(buf *Ctx, p uintptr) error {
	value := *(*uint16)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint32(buf *Ctx, p uintptr) error {
	value := *(*uint32)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func encodeUint64(buf *Ctx, p uintptr) error {
	value := *(*uint64)(unsafe.Pointer(p))
	buf.b = append(buf.b, 'i', ':')
	buf.b = strconv.AppendUint(buf.b, value, 10)
	buf.b = append(buf.b, ';')
	return nil
}

func encodeUint(buf *Ctx, p uintptr) error {
	value := *(*uint)(unsafe.Pointer(p))
	appendUint(buf, uint64(value))
	return nil
}

func appendUint(buf *Ctx, v uint64) {
	buf.b = append(buf.b, 'i', ':')
	buf.b = strconv.AppendUint(buf.b, v, 10)
	buf.b = append(buf.b, ';')
}
