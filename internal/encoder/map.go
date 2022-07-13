package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// !!! not safe to use in reflect case !!!
func compileMap(rt *runtime.Type) (encoder, error) {
	// for map[int]string, keyType is int, valueType is string
	keyType := rt.Key()
	valueType := rt.Elem()

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, &UnsupportedTypeAsMapKeyError{Type: runtime.RType2Type(keyType)}
	}

	keyEncoder, err := compileMapKey(keyType)
	if err != nil {
		return nil, err
	}

	var valueEncoder encoder

	// need special take care
	// fmt.Println(runtime.IfaceIndir(rt), runtime.IfaceIndir(valueType), rt.String())
	if valueType.Kind() == reflect.Map {
		enc, err := compileMap(valueType)
		if err != nil {
			return nil, err
		}
		valueEncoder = deRefNilEncoder(enc)
	} else {
		valueEncoder, err = compile(valueType)
		if err != nil {
			return nil, err
		}
	}

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		if p == 0 {
			// nil
			return appendNull(b), nil
		}

		ptr := ptrToUnsafePtr(p)

		mapLen := runtime.MapLen(ptr)
		if mapLen == 0 {
			return appendEmptyArray(b), nil
		}

		b = appendArrayBegin(b, int64(mapLen))

		var mapCtx = newMapCtx()
		defer freeMapCtx(mapCtx)

		runtime.MapIterInit(rt, ptr, &mapCtx.Iter)
		var err error // create a new error value, so shadow compiler's error
		for i := 0; i < mapLen; i++ {
			b, err = keyEncoder(ctx, b, runtime.MapIterKey(&mapCtx.Iter))
			if err != nil {
				return b, err
			}

			b, err = valueEncoder(ctx, b, runtime.MapIterValue(&mapCtx.Iter))
			if err != nil {
				return b, err
			}

			runtime.MapIterNext(&mapCtx.Iter)
		}
		return append(b, '}'), nil
	}, nil
}
