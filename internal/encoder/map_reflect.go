package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// fast array for map reflect
var mapKeyEncoder = [25]encoder{
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

	// iter keys and values
	// non interface value type, fast path with uintptr
	if rt.Elem().Kind() != reflect.Interface {
		return reflectConcreteMap(ctx, rt, rv, keyType)
	}

	mapLen := rv.Len()
	if mapLen == 0 {
		appendEmptyArray(ctx)
		return nil
	}

	appendArrayBegin(ctx, int64(mapLen))

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = newMapCtx()
	defer freeMapCtx(mr)

	valueEncoder, err := compileInterface(runtime.Type2RType(rt.Elem()))
	if err != nil {
		return err
	}

	mapIterInit(runtime.Type2RType(rt), unsafe.Pointer(rv.Pointer()), &mr.Iter)
	for i := 0; i < mapLen; i++ {
		err := keyEncoder(ctx, uintptr(mapIterKey(&mr.Iter)))
		if err != nil {
			return err
		}

		iterElem := mapIterValue(&mr.Iter)

		// value := *(*reflect.Value)(unsafe.Pointer(&rValue{m.typ, uintptr(iterElem), ro(m.flag) | flag(vtype)}))

		err = valueEncoder(ctx, uintptr(iterElem))
		if err != nil {
			return err
		}

		mapIterNext(&mr.Iter)
	}

	ctx.b = append(ctx.b, '}')

	return nil
}

func reflectConcreteMap(ctx *Ctx, rt reflect.Type, rv reflect.Value, keyType reflect.Type) error {
	mapLen := rv.Len()
	if mapLen == 0 {
		appendEmptyArray(ctx)
		return nil
	}

	appendArrayBegin(ctx, int64(mapLen))

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
			return err
		}

		typeToEncoderMap.Store(valueTypeID, valueEncoder)
	} else {
		valueEncoder = enc.(encoder)
	}

	if rt.Elem().Kind() == reflect.Map {
		originValueEncoder := valueEncoder
		valueEncoder = func(ctx *Ctx, p uintptr) error {
			return originValueEncoder(ctx, ptrOfPtr(p))
		}
	}

	keyEncoder := mapKeyEncoder[keyType.Kind()]

	var mr = newMapCtx()
	defer freeMapCtx(mr)

	mapIterInit(runtime.Type2RType(rt), unsafe.Pointer(rv.Pointer()), &mr.Iter)
	for i := 0; i < mapLen; i++ {
		err := keyEncoder(ctx, uintptr(mapIterKey(&mr.Iter)))
		if err != nil {
			return err
		}

		value := uintptr(mapIterValue(&mr.Iter))
		err = valueEncoder(ctx, value)
		if err != nil {
			return err
		}
		mapIterNext(&mr.Iter)
	}
	ctx.b = append(ctx.b, '}')
	return nil
}
