package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

// fast array for map reflect
var mapKeyEncoder = []encoder{
	reflect.String: encodeStringVariable,
	reflect.Int:    encodeInt,
	reflect.Int8:   encodeInt8,
	reflect.Int16:  encodeInt16,
	reflect.Int32:  encodeInt32,
	reflect.Int64:  encodeInt64,
	reflect.Uint:   encodeUint,
	reflect.Uint8:  encodeUint8,
	reflect.Uint16: encodeUint16,
	reflect.Uint32: encodeUint32,
	reflect.Uint64: encodeUint64,
}

func reflectMap(ctx *Ctx, rv reflect.Value) error {
	rt := rv.Type()
	keyType := rt.Key()

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return &UnsupportedTypeAsMapKeyError{Type: rv.Type().Key()}
	}

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

	// iter keys and values
	// non interface value type, fast path with uintptr
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteMap(ctx, rt, rv, keyType)
	}

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = mapIterPool.Get().(*mapIter)
	defer mapIterPool.Put(mr)
	defer mr.reset()

	// map length is checked, always > 0
	mr.m = *(*rValue)(unsafe.Pointer(&rv))
	for mr.Next() {
		err := keyEncoder(ctx, mr.KeyUnsafeAddress())
		if err != nil {
			return err
		}

		err = reflectInterfaceValue(ctx, mr.Value())
		if err != nil {
			return err
		}
	}

	ctx.b = append(ctx.b, '}')

	return nil
}

func reflectConcreteMap(ctx *Ctx, rt reflect.Type, rv reflect.Value, keyType reflect.Type) error {
	// map has a different reflect.Value{}.flag.
	// map's address may be direct or indirect address
	var valueEncoder encoder
	var err error
	var valueType = rt.Elem()

	valueTypeID := uintptr(unsafe.Pointer(valueType))

	enc, ok := typeToEncoderMap.Load(valueTypeID)
	if !ok {
		valueEncoder, err = compile(valueType, reflect.New(valueType).Elem())
		if err != nil {
			return err
		}

		typeToEncoderMap.Store(valueTypeID, valueEncoder)
	} else {
		valueEncoder = enc.(encoder)
	}

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = mapIterPool.Get().(*mapIter)
	defer mapIterPool.Put(mr)
	defer mr.reset()

	mr.m = reflectValueToLocal(rv)

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
}
