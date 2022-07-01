package encoder

import (
	"sync"
	"unsafe"

	"github.com/goccy/go-reflect"
)

var mapIterPool = sync.Pool{
	New: func() interface{} {
		return &mapIter{}
	},
}

// !!! not safe to use in reflect case !!!
func compileMap(rt reflect.Type, rv reflect.Value) (encoder, error) {
	// for map[int]string, keyType is int, valueType is string
	keyType := rt.Key()
	valueType := rt.Elem()

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, &UnsupportedTypeAsMapKeyError{Type: keyType}
	}

	keyEncoder, err := compile(keyType, reflect.New(keyType).Elem())
	if err != nil {
		return nil, err
	}

	valueEncoder, err := compile(valueType, reflect.New(valueType).Elem())
	if err != nil {
		return nil, err
	}

	flag := reflectValueToLocal(rv).flag

	return func(ctx *Ctx, p uintptr) error {
		rv := reflectValueMapFromPtr(rt, p, flag)

		if rv.IsNil() {
			appendNil(ctx)
			return nil
		}

		mapLen := rv.Len()
		appendArrayBegin(ctx, int64(mapLen))

		if mapLen == 0 {
			ctx.b = append(ctx.b, '}')
			return nil
		}

		var mr = mapIterPool.Get().(*mapIter)
		defer mapIterPool.Put(mr)
		defer mr.reset()

		mr.m = *(*rValue)(unsafe.Pointer(&rv))
		for mr.Next() {
			err := keyEncoder(ctx, mr.KeyUnsafeAddress())
			if err != nil {
				return err
			}

			err = valueEncoder(ctx, mr.ValueUnsafeAddress())
			if err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}

type mapEncoder struct {
	rt           reflect.Type
	flag         uintptr
	keyEncoder   encoder
	valueEncoder encoder
}

func (e *mapEncoder) encode(ctx *Ctx, p uintptr) error {
	rv := reflectValueMapFromPtr(e.rt, p, e.flag)

	if rv.IsNil() {
		appendNil(ctx)
		return nil
	}

	mapLen := rv.Len()
	appendArrayBegin(ctx, int64(mapLen))

	if mapLen == 0 {
		ctx.b = append(ctx.b, '}')
		return nil
	}

	var mr = mapIterPool.Get().(*mapIter)
	defer mapIterPool.Put(mr)
	defer mr.reset()

	mr.m = *(*rValue)(unsafe.Pointer(&rv))
	for mr.Next() {
		err := e.keyEncoder(ctx, mr.KeyUnsafeAddress())
		if err != nil {
			return err
		}

		err = e.valueEncoder(ctx, mr.ValueUnsafeAddress())
		if err != nil {
			return err
		}
	}

	ctx.b = append(ctx.b, '}')
	return nil
}
