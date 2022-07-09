package encoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileBoolAsString(typ *runtime.Type) (encoder, error) {
	return encodeBoolAsString, nil
}

func encodeBoolAsString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	value := *(*bool)(unsafe.Pointer(p))
	if value {
		return append(b, `s:4:"true";`...), nil
	}
	return append(b, `s:5:"false";`...), nil
}
