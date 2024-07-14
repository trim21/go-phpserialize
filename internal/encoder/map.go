package encoder

import (
	"reflect"
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

	valueEncoder, err := compile(valueType, seen)
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		b = appendArrayBegin(b, int64(rv.Len()))

		keys := rv.MapKeys()
		for _, key := range keys {
			b, err = keyEncoder(ctx, b, key)
			if err != nil {
				return b, err
			}
			b, err = valueEncoder(ctx, b, rv.MapIndex(key))
			if err != nil {
				return b, err
			}
		}

		return append(b, '}'), nil
	}, nil
}
