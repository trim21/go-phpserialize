package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func reflectStruct(ctx *Ctx, b []byte, rv reflect.Value, p uintptr) ([]byte, error) {
	enc, err := compileWithCache(runtime.Type2RType(rv.Type()))
	if err != nil {
		return nil, err
	}

	return enc(ctx, b, p)
}
