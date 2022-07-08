package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func reflectStruct(ctx *Ctx, rv reflect.Value, p uintptr) error {
	rt := runtime.Type2RType(rv.Type())
	typeID := uintptr(unsafe.Pointer(rt))
	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		return enc.(encoder)(ctx, p)
	}

	encoder, err := compileStruct(rt)
	if err != nil {
		return err
	}

	typeToEncoderMap.Store(typeID, encoder)

	return encoder(ctx, p)
}
