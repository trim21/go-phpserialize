package encoder

import (
	"reflect"
)

// !!! not safe to use in reflect case !!!
func compileMap(rt reflect.Type, seen compileSeenMap) (encoder, error) {
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

	keyEncoder, compileErr := compileMapKey(keyType)
	if compileErr != nil {
		return nil, compileErr
	}

	valueEncoder, compileErr := compile(valueType, seen)
	if compileErr != nil {
		return nil, compileErr
	}

	return checkRecursiveEncoder(func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		size := rv.Len()
		b = appendArrayBegin(b, int64(size))

		iter := rv.MapRange()

		kv := reflect.New(keyType).Elem()
		vv := reflect.New(valueType).Elem()

		var err error
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
	}), nil
}
