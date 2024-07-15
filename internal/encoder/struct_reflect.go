package encoder

import (
	"reflect"
)

func reflectStruct(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	enc, err := compileWithCache(rv.Type())
	if err != nil {
		return nil, err
	}

	return enc(ctx, b, rv)
}
