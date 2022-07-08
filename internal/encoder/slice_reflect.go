package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

func reflectSlice(ctx *Ctx, rv reflect.Value) error {
	l := rv.Len()
	rt := rv.Type()

	// not slice of interface, fast path
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteSlice(ctx, rt, rv)
	}

	// slow path with O(N) allocation.
	appendArrayBegin(ctx, int64(l))
	for i := 0; i < l; i++ {
		appendInt(ctx, int64(i))
		err := reflectInterfaceValue(ctx, rv.Index(i))
		if err != nil {
			return err
		}
	}
	ctx.b = append(ctx.b, '}')
	return nil
}

func reflectConcreteSlice(ctx *Ctx, rt reflect.Type, rv reflect.Value) error {
	typeID := uintptr(unsafe.Pointer(rt))

	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		return enc.(encoder)(ctx, reflectValueToLocal(rv).ptr)
	}

	enc, err := compile(rt)
	if err != nil {
		panic(err)
	}

	return enc(ctx, reflectValueToLocal(rv).ptr)
}
