package encoder

import (
	"reflect"
	"unsafe"
)

func compileBoolAsString(typ reflect.Type) (encoder, error) {
	return encodeBoolAsString, nil
}

func encodeBoolAsString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**bool)(unsafe.Pointer(&p))
	if value {
		return append(b, `s:4:"true";`...), nil
	}
	return append(b, `s:5:"false";`...), nil
}
