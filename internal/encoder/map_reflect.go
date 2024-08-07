package encoder

import (
	"reflect"
)

// fast array for map reflect
var mapKeyEncoder = [25]encoder{
	reflect.String: encodeString,
	reflect.Int:    encodeInt,
	reflect.Int8:   encodeInt,
	reflect.Int16:  encodeInt,
	reflect.Int32:  encodeInt,
	reflect.Int64:  encodeInt,
	reflect.Uint:   encodeUint,
	reflect.Uint8:  encodeUint,
	reflect.Uint16: encodeUint,
	reflect.Uint32: encodeUint,
	reflect.Uint64: encodeUint,
}

func reflectMap(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	rt := rv.Type()
	keyType := rt.Key()

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, &UnsupportedTypeAsMapKeyError{Type: rv.Type().Key()}
	}

	if rv.IsNil() {
		return appendNull(b), nil
	}

	// iter keys and values
	// non interface value type, fast path with uintptr
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteMap(ctx, b, rt, rv, keyType)
	}

	mapLen := rv.Len()
	if mapLen == 0 {
		return appendEmptyArray(b), nil
	}

	b = appendArrayBegin(b, int64(mapLen))

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var valueType = rt.Elem()

	valueEncoder, err := compileInterface(valueType)
	if err != nil {
		return b, err
	}

	iter := rv.MapRange()

	kv := reflect.New(keyType).Elem()
	vv := reflect.New(valueType).Elem()

	for iter.Next() {
		kv.SetIterKey(iter)
		b, err = keyEncoder(ctx, b, kv)
		if err != nil {
			return b, err
		}

		vv.SetIterValue(iter)
		b, err = valueEncoder(ctx, b, vv)
		if err != nil {
			return b, err
		}
	}

	return append(b, '}'), nil
}

func reflectConcreteMap(ctx *Ctx, b []byte, rt reflect.Type, rv reflect.Value, keyType reflect.Type) ([]byte, error) {
	mapLen := rv.Len()
	if mapLen == 0 {
		return appendEmptyArray(b), nil
	}

	b = appendArrayBegin(b, int64(mapLen))

	// map has a different reflect.Value{}.flag.
	// map's address may be direct or indirect address
	var valueEncoder encoder
	var err error
	var valueType = rt.Elem()

	valueEncoder, err = compileWithCache(valueType)
	if err != nil {
		return nil, err
	}

	if rt.Elem().Kind() == reflect.Map {
		originValueEncoder := valueEncoder
		valueEncoder = originValueEncoder
	}

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	iter := rv.MapRange()

	kv := reflect.New(keyType).Elem()
	vv := reflect.New(valueType).Elem()

	for iter.Next() {
		kv.SetIterKey(iter)
		b, err = keyEncoder(ctx, b, kv)
		if err != nil {
			return b, err
		}

		vv.SetIterValue(iter)
		b, err = valueEncoder(ctx, b, vv)
		if err != nil {
			return b, err
		}
	}

	return append(b, '}'), nil
}
