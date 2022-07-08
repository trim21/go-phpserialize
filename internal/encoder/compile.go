package encoder

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

const DefaultStructTag = "php"

var (
	typeToEncoderMap sync.Map
	ctxPool          = sync.Pool{
		New: func() any {
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

func compile(rt *runtime.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return compileBool(rt)
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
	case reflect.String:
		return encodeStringVariable, nil
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
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	case reflect.Struct:
		return compileStruct(rt)
	case reflect.Slice:
		return compileSlice(rt)
	case reflect.Map:
		return compileMap(rt)
	case reflect.Interface:
		return compileInterface(rt)
	case reflect.Ptr:
		return compilePtr(rt)
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s)", rt.String(), rt.Kind())
}

func compileMapKey(typ *runtime.Type) (encoder, error) {
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

func compileAsString(rt *runtime.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return compileBoolAsString(rt)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return compileIntAsString(rt)
	case reflect.String:
		return encodeStringVariable, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return compileUintAsString(rt)
	case reflect.Float32, reflect.Float64:
		return compileFloatAsString(rt)
	case reflect.Interface:
		return compileInterfaceAsString(rt)

	}

	return nil, fmt.Errorf(
		"failed to build encoder for struct field (as string), unsupported type %s (kind %s)",
		rt.String(), rt.Kind())
}