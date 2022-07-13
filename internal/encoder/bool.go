package encoder

import (
	"unsafe"
)

func encodeBool(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	v := **(**bool)(unsafe.Pointer(&p))

	b = append(b, 'b', ':')
	if v {
		b = append(b, '1')
	} else {
		b = append(b, '0')
	}

	return append(b, ';'), nil
}
