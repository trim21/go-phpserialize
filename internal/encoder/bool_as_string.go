package encoder

import (
	"reflect"
	"unsafe"
)

func compileBoolAsString(typ reflect.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.Bool() {
			return append(b, `s:4:"true";`...), nil
		}
		return append(b, `s:5:"false";`...), nil
	}, nil
}

func encodeBoolAsString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := **(**bool)(unsafe.Pointer(&p))
	if value {
		return append(b, `s:4:"true";`...), nil
	}
	return append(b, `s:5:"false";`...), nil
}
