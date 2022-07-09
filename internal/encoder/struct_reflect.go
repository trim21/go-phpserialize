package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func reflectStruct(ctx *Ctx, b []byte, rv reflect.Value, p uintptr) ([]byte, error) {
	rt := runtime.Type2RType(rv.Type())
	typeID := uintptr(unsafe.Pointer(rt))
	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		return enc.(encoder)(ctx, b, p)
	}

	encoder, err := compileStruct(rt)
	if err != nil {
		return nil, err
	}

	typeToEncoderMap.Store(typeID, encoder)

	return encoder(ctx, b, p)
}
