package encoder

import (
	"unsafe"
)

func Marshal(v any) ([]byte, error) {
	header := (*emptyInterface)(unsafe.Pointer(&v))

	typeID := uintptr(unsafe.Pointer(header.typ))

	ptr := uintptr(header.ptr)

	enc, err := compileTypeIDWithCache(typeID)
	if err != nil {
		return nil, err
	}

	ctx := newCtx()
	defer freeCtx(ctx)

	ctx.KeepRefs = append(ctx.KeepRefs, header.ptr)

	buf := newBuffer()
	defer freeBuffer(buf)

	buf.b, err = enc(ctx, buf.b, ptr)
	if err != nil {
		return nil, err
	}

	// allocate a new Ctx required length only
	p := make([]byte, len(buf.b))

	copy(p, buf.b)
	return p, nil
}

func MarshalNoEscape(v any) ([]byte, error) {
	header := (*emptyInterface)(unsafe.Pointer(&v))

	typeID := uintptr(unsafe.Pointer(header.typ))

	ptr := uintptr(header.ptr)

	enc, err := compileTypeIDWithCache(typeID)
	if err != nil {
		return nil, err
	}

	ctx := newCtx()
	defer freeCtx(ctx)

	// ctx.KeepRefs = append(ctx.KeepRefs, header.ptr)

	buf := newBuffer()
	defer freeBuffer(buf)

	buf.b, err = enc(ctx, buf.b, ptr)
	if err != nil {
		return nil, err
	}

	// allocate a new Ctx required length only
	p := make([]byte, len(buf.b))

	copy(p, buf.b)
	return p, nil
}
