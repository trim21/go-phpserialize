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

func compileMap(typ reflect.Type, rv reflect.Value) (encoder, error) {
	// for map[int]string, keyType is int, valueType is string
	keyType := typ.Key()
	valueType := typ.Elem()

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, &UnsupportedTypeAsMapKeyError{Type: keyType}
	}

	if valueType.Kind() == reflect.Interface {
		// interface slow path
		return compileMapWithInterfaceValue(typ, rv, keyType, valueType)
	}

	keyEncoder, err := compile(keyType, reflect.New(keyType).Elem())
	if err != nil {
		return nil, err
	}

	valueEncoder, err := compile(valueType, reflect.New(valueType).Elem())
	if err != nil {
		return nil, err
	}
	flag := (*(*rValue)(unsafe.Pointer(&rv))).flag

	return func(ctx *Ctx, p uintptr) error {
		rv := reflectValueMapFromPtr(typ, p, flag)

		if rv.IsNil() {
			appendArrayBegin(ctx, 0)
			ctx.b = append(ctx.b, '}')

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
			err := keyEncoder(ctx, mr.Key())
			if err != nil {
				return err
			}

			err = valueEncoder(ctx, mr.Value())
			if err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}

func compileMapWithInterfaceValue(typ reflect.Type, rv reflect.Value, keyType, valueType reflect.Type) (encoder, error) {
	keyEncoder, err := compile(keyType, rv)
	if err != nil {
		return nil, err
	}

	flag := (*(*rValue)(unsafe.Pointer(&rv))).flag

	valueEncoder, err := compileInterface(valueType)
	if err != nil {
		return nil, err
	}

	return func(buf *Ctx, p uintptr) error {
		rv := reflectValueMapFromPtr(typ, p, flag)

		if rv.IsNil() {
			appendArrayBegin(buf, 0)

			return nil
		}

		appendArrayBegin(buf, int64(rv.Len()))

		mr := rv.MapRange()
		for mr.Next() {
			// HINT: This is very likely to break after new go version.
			p := mr.Key()
			v := (*rValue)(unsafe.Pointer(&p))

			err := keyEncoder(buf, v.ptr)
			if err != nil {
				return err
			}

			p = mr.Value()
			v = (*rValue)(unsafe.Pointer(&p))
			err = valueEncoder(buf, v.ptr)
			if err != nil {
				return err
			}
		}
		buf.b = append(buf.b, '}')
		return nil
	}, nil
}
