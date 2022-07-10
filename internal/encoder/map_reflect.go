package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// fast array for map reflect
var mapKeyEncoder = [25]encoder{
	reflect.String: encodeString,
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
		return appendNilBytes(b), nil
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

	b = appendArrayBeginBytes(b, int64(mapLen))

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = newMapCtx()
	defer freeMapCtx(mr)

	valueEncoder, err := compileInterface(runtime.Type2RType(rt.Elem()))
	if err != nil {
		return b, err
	}

	runtime.MapIterInit(runtime.Type2RType(rt), unsafe.Pointer(rv.Pointer()), &mr.Iter)
	for i := 0; i < mapLen; i++ {
		b, err = keyEncoder(ctx, b, runtime.MapIterKey(&mr.Iter))
		if err != nil {
			return b, err
		}

		b, err = valueEncoder(ctx, b, runtime.MapIterValue(&mr.Iter))
		if err != nil {
			return b, err
		}

		runtime.MapIterNext(&mr.Iter)
	}

	b = append(b, '}')

	return b, nil
}

func reflectConcreteMap(ctx *Ctx, b []byte, rt reflect.Type, rv reflect.Value, keyType reflect.Type) ([]byte, error) {
	mapLen := rv.Len()
	if mapLen == 0 {
		return appendEmptyArray(b), nil
	}

	b = appendArrayBeginBytes(b, int64(mapLen))

	// map has a different reflect.Value{}.flag.
	// map's address may be direct or indirect address
	var valueEncoder encoder
	var err error
	var valueType = rt.Elem()

	valueTypeID := uintptr(unsafe.Pointer(runtime.Type2RType(valueType)))

	enc, ok := typeToEncoderMap.Load(valueTypeID)
	if !ok {
		valueEncoder, err = compile(runtime.Type2RType(valueType))
		if err != nil {
			return b, err
		}

		typeToEncoderMap.Store(valueTypeID, valueEncoder)
	} else {
		valueEncoder = enc.(encoder)
	}

	if rt.Elem().Kind() == reflect.Map {
		originValueEncoder := valueEncoder
		valueEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
			return originValueEncoder(ctx, b, PtrDeRef(p))
		}
	}

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = newMapCtx()
	defer freeMapCtx(mr)

	runtime.MapIterInit(runtime.Type2RType(rt), unsafe.Pointer(rv.Pointer()), &mr.Iter)
	for i := 0; i < mapLen; i++ {
		b, err = keyEncoder(ctx, b, runtime.MapIterKey(&mr.Iter))
		if err != nil {
			return b, err
		}

		b, err = valueEncoder(ctx, b, runtime.MapIterValue(&mr.Iter))
		if err != nil {
			return b, err
		}

		runtime.MapIterNext(&mr.Iter)
	}

	return append(b, '}'), nil
}
