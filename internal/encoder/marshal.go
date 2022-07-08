package encoder

import (
	"unsafe"
)

func Marshal(v any) ([]byte, error) {
	// Get type information and pointer from interface{} value without allocation.
	header := (*emptyInterface)(unsafe.Pointer(&v))

	typeID := uintptr(unsafe.Pointer(header.typ))

	// Reuse the Ctx once allocated using sync.Pool
	ctx := newCtx()
	defer freeCtx(ctx)

	// value will have a writing barrier until we release it.
	ctx.KeepRefs = append(ctx.KeepRefs, header.ptr)

	// builds an optimized path by typeID and caches it
	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		if err := enc.(encoder)(ctx, uintptr(header.ptr)); err != nil {
			return nil, err
		}

		// allocate a new Ctx required length only
		b := make([]byte, len(ctx.b))
		copy(b, ctx.b)
		return b, nil
	}

	// First time,
	// builds an optimized path by type and caches it with typeID.
	enc, err := compile(header.typ)
	if err != nil {
		return nil, err
	}
	typeToEncoderMap.Store(typeID, enc)
	if err := enc(ctx, uintptr(header.ptr)); err != nil {
		return nil, err
	}

	// allocate a new Ctx required length only
	b := make([]byte, len(ctx.b))

	copy(b, ctx.b)
	return b, nil
}
