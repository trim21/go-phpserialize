package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func unpackIface(p uintptr) uintptr {
	return uintptr((*(*emptyInterface)(unsafe.Pointer(p))).ptr)
}

func reflectSlice(ctx *Ctx, rv reflect.Value, p uintptr) error {
	rt := rv.Type()

	// not slice of interface, fast path
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteSlice(ctx, runtime.Type2RType(rt), p)
	}

	shPtr := unpackIface(p)
	// no data ptr, nil slice
	// even empty slice has a non-zero data ptr
	if shPtr == 0 {
		appendNil(ctx)
		return nil
	}

	el := runtime.Type2RType(rt.Elem())

	encoder, err := compileInterface(el)
	if err != nil {
		return err
	}

	sh := *(*runtime.SliceHeader)(unsafe.Pointer(shPtr))
	offset := rt.Elem().Size()

	dataPtr := uintptr(sh.Data)
	appendArrayBegin(ctx, int64(sh.Len))
	for i := 0; i < sh.Len; i++ {
		appendInt(ctx, int64(i))
		err := encoder(ctx, dataPtr+uintptr(i)*offset)
		if err != nil {
			return err
		}
	}
	ctx.b = append(ctx.b, '}')
	return nil
}

func reflectConcreteSlice(ctx *Ctx, rt *runtime.Type, p uintptr) error {
	var typeID = uintptr(unsafe.Pointer(rt))

	p = unpackIface(p)

	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		return enc.(encoder)(ctx, p)
	}

	enc, err := compile(rt)
	if err != nil {
		panic(err)
	}

	return enc(ctx, p)
}
