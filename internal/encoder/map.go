package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// !!! not safe to use in reflect case !!!
func compileMap(rt reflect.Type, seen seenMap) (encoder, error) {
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

	keyEncoder, err := compileMapKey(keyType)
	if err != nil {
		return nil, err
	}

	var valueEncoder encoder

	// need special take care
	// fmt.Println(runtime.IfaceIndir(rt), runtime.IfaceIndir(valueType), rt.String())
	if valueType.Kind() == reflect.Map {
		enc, err := compileMap(valueType, seen)
		if err != nil {
			return nil, err
		}
		valueEncoder = deRefNilEncoder(enc)
	} else {
		valueEncoder, err = compile(valueType, seen)
		if err != nil {
			return nil, err
		}
	}

	typeID := runtime.ToTypeID(rt)

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		if p == 0 {
			// nil
			return appendNull(b), nil
		}

		ptr := ptrToUnsafePtr(p)

		rv := reflect.ValueOf(*(*any)(unsafe.Pointer(&eface{
			typ: typeID,
			ptr: ptr,
		})))

		mapLen := rv.Len()
		if mapLen == 0 {
			return appendEmptyArray(b), nil
		}

		b = appendArrayBegin(b, int64(mapLen))

		keys := rv.MapKeys()

		for _, key := range keys {
			b, err = keyEncoder(ctx, b, key.Pointer())
			if err != nil {
				return b, err
			}
			b, err = valueEncoder(ctx, b, rv.MapIndex(key).Pointer())
			if err != nil {
				return b, err
			}
		}

		return append(b, '}'), nil
	}, nil
}
