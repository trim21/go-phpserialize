package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func unpackIface(p uintptr) uintptr {
	return uintptr((*(*emptyInterface)(unsafe.Pointer(p))).ptr)
}

func reflectSlice(ctx *Ctx, b []byte, rv reflect.Value, p uintptr) ([]byte, error) {
	rt := rv.Type()

	// not slice of interface, fast path
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteSlice(ctx, b, runtime.Type2RType(rt), p)
	}

	shPtr := unpackIface(p)
	// no data ptr, nil slice
	// even empty slice has a non-zero data ptr
	if shPtr == 0 {
		return appendNilBytes(b), nil
	}

	el := runtime.Type2RType(rt.Elem())

	encoder, err := compileInterface(el)
	if err != nil {
		return nil, err
	}

	sh := **(**runtime.SliceHeader)(unsafe.Pointer(&shPtr))
	offset := rt.Elem().Size()

	dataPtr := uintptr(sh.Data)

	b = appendArrayBeginBytes(b, int64(sh.Len))

	for i := 0; i < sh.Len; i++ {
		b = appendIntBytes(b, int64(i))
		b, err = encoder(ctx, b, dataPtr+uintptr(i)*offset)
		if err != nil {
			return b, err
		}
	}

	return append(b, '}'), nil
}

func reflectConcreteSlice(ctx *Ctx, b []byte, rt *runtime.Type, p uintptr) ([]byte, error) {
	var typeID = uintptr(unsafe.Pointer(rt))

	p = unpackIface(p)

	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		return enc.(encoder)(ctx, b, p)
	}

	enc, err := compileWithCache(rt)
	if err != nil {
		return b, err
	}

	return enc(ctx, b, p)
}
