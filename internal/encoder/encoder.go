package encoder

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/goccy/go-reflect"
)

const DefaultStructTag = "php"

var (
	typeToEncoderMap sync.Map
	ctxPool          = sync.Pool{
		New: func() interface{} {
			return &Ctx{
				b:        make([]byte, 0, 1024),
				KeepRefs: make([]unsafe.Pointer, 0, 8),
			}
		},
	}
)

type Ctx struct {
	b        []byte
	KeepRefs []unsafe.Pointer

	// a buffer to encode float as string
	floatBuffer []byte
}

func newCtx() *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.b = ctx.b[:0]

	return ctx
}

func freeCtx(ctx *Ctx) {
	ctx.KeepRefs = ctx.KeepRefs[:0]
	ctx.floatBuffer = ctx.floatBuffer[:0]

	ctxPool.Put(ctx)
}

type encoder func(ctx *Ctx, p uintptr) error

func Marshal(v interface{}) ([]byte, error) {
	// Technique 1.
	// Get type information and pointer from interface{} rValue without allocation.
	typ, ptr := reflect.TypeAndPtrOf(v)
	// so value will have a writing barrier until we release it.
	header := (*emptyInterface)(unsafe.Pointer(&v))

	typeID := uintptr(unsafe.Pointer(header.typ))

	p := uintptr(ptr)

	// Technique 2.
	// Reuse the Ctx once allocated using sync.Pool
	ctx := newCtx()
	defer freeCtx(ctx)
	ctx.KeepRefs = append(ctx.KeepRefs, header.ptr)
	// ctx.KeepRefs = append(ctx.KeepRefs, unsafe.Pointer(&p))

	// Technique 3.
	// builds an optimized path by typeID and caches it
	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		if err := enc.(encoder)(ctx, p); err != nil {
			return nil, err
		}

		// allocate a new Ctx required length only
		b := make([]byte, len(ctx.b))
		copy(b, ctx.b)
		return b, nil
	}

	// First time,
	// builds an optimized path by type and caches it with typeID.
	rv := reflect.ValueOf(v)
	enc, err := compile(typ, rv)
	if err != nil {
		return nil, err
	}
	typeToEncoderMap.Store(typeID, enc)
	if err := enc(ctx, p); err != nil {
		return nil, err
	}

	// allocate a new Ctx required length only
	b := make([]byte, len(ctx.b))

	copy(b, ctx.b)
	return b, nil
}

func compile(typ reflect.Type, rv reflect.Value) (encoder, error) {
	switch typ.Kind() {
	case reflect.Bool:
		return compileBool(typ)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return compileInt(typ)
	case reflect.String:
		return encodeStringVariable, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return compileUint(typ)
	case reflect.Float32, reflect.Float64:
		return compileFloat(typ)
	case reflect.Struct:
		return compileStruct(typ, rv)
	case reflect.Slice:
		return compileSlice(typ, rv)
	case reflect.Map:
		return compileMap(typ, rv)
	case reflect.Interface:
		return compileInterface(typ)
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s)", typ.String(), typ.Kind())
}

func compileMapKey(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.String:
		return encodeStringVariable, nil

	case reflect.Int8:
		return encodeInt8, nil
	case reflect.Int16:
		return encodeInt16, nil
	case reflect.Int32:
		return encodeInt32, nil
	case reflect.Int64:
		return encodeInt64, nil
	case reflect.Int:
		return encodeInt, nil
	case reflect.Uint8:
		return encodeUint8, nil
	case reflect.Uint16:
		return encodeUint16, nil
	case reflect.Uint32:
		return encodeUint32, nil
	case reflect.Uint64:
		return encodeUint64, nil
	case reflect.Uint:
		return encodeUint, nil
	}

	return nil, fmt.Errorf("failed to build encoder for map key, unsupported type %s (kind %s)", typ.String(), typ.Kind())
}

func compileAsString(typ reflect.Type, rv reflect.Value) (encoder, error) {
	switch typ.Kind() {
	case reflect.Bool:
		return compileBoolAsString(typ)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return compileIntAsString(typ)
	case reflect.String:
		return encodeStringVariable, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return compileUintAsString(typ)
	case reflect.Float32, reflect.Float64:
		return compileFloatAsString(typ)
	}

	return nil, fmt.Errorf(
		"failed to build encoder for struct field (as string), unsupported type %s (kind %s)",
		typ.String(), typ.Kind())
}
