package encoder

import (
	"reflect"
	"unsafe"
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
	header := (*emptyInterface)(unsafe.Pointer(&v))

	enc, err := compileWithCache(reflect.TypeOf(v))
	if err != nil {
		return nil, err
	}

	ptr := uintptr(header.ptr)
	ctx.KeepRefs = append(ctx.KeepRefs, header.ptr)

	return enc(ctx, b, ptr)
}
