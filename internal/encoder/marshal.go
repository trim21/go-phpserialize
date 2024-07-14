package encoder

import (
	"reflect"
)

func Marshal(v any) ([]byte, error) {
	ctx := newCtx()

	b, err := encode(ctx, ctx.Buf[:0], v)
	if err != nil {
		freeCtx(ctx)
		return nil, err
	}

	dst := make([]byte, len(b))
	copy(dst, b)

	ctx.Buf = b
	freeCtx(ctx)

	return dst, nil
}

func encode(ctx *Ctx, b []byte, v any) ([]byte, error) {
	rv := reflect.ValueOf(v)

	enc, err := compileWithCache(rv.Type())
	if err != nil {
		return nil, err
	}

	return enc(ctx, b, rv)
}
