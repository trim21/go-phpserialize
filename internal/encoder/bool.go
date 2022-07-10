package encoder

import (
	"unsafe"
)

func encodeBool(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**bool)(unsafe.Pointer(&p))

	return appendBool(b, value), nil
}

func appendBool(b []byte, v bool) []byte {
	b = append(b, 'b', ':')
	if v {
		b = append(b, '1')
	} else {
		b = append(b, '0')
	}
	return append(b, ';')
}
