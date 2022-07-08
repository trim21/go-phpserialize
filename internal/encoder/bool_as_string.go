package encoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileBoolAsString(typ *runtime.Type) (encoder, error) {
	return encodeBoolAsString, nil
}

func encodeBoolAsString(ctx *Ctx, p uintptr) error {
	value := *(*bool)(unsafe.Pointer(p))
	appendBoolAsString(ctx, value)
	return nil
}

func appendBoolAsString(ctx *Ctx, v bool) {
	if v {
		ctx.b = append(ctx.b, `s:4:"true";`...)
	} else {
		ctx.b = append(ctx.b, `s:5:"false";`...)
	}
}
