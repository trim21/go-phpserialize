package encoder

import (
	"reflect"
	"sync"
	"unsafe"

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
	if valueType.Kind() == reflect.Map {
		valueEncoder, err = compile(runtime.PtrTo(valueType))
		if err != nil {
			return nil, err
		}
	} else {
		valueEncoder, err = compile(valueType)
		if err != nil {
			return nil, err
		}
	}

	// reflect.ValueOf(map[int]int{}).MapIndex()

	return func(ctx *Ctx, p uintptr) error {
		if p == 0 {
			// nil
			appendNil(ctx)
			return nil
		}

		ptr := ptrToUnsafePtr(p)

		mapLen := runtime.MapLen(ptr)

		if mapLen == 0 {
			appendEmptyArray(ctx)
			return nil
		}

		appendArrayBegin(ctx, int64(mapLen))

		var mapCtx = newMapCtx()
		defer freeMapCtx(mapCtx)

		ctx.KeepRefs = append(ctx.KeepRefs, unsafe.Pointer(mapCtx))

		mapIterInit(rt, ptr, &mapCtx.Iter)
		for i := 0; i < mapLen; i++ {
			err := keyEncoder(ctx, uintptr(mapIterKey(&mapCtx.Iter)))
			if err != nil {
				return err
			}

			err = valueEncoder(ctx, uintptr(mapIterValue(&mapCtx.Iter)))
			if err != nil {
				return err
			}

			mapIterNext(&mapCtx.Iter)
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}

var mapCtxPool = sync.Pool{
	New: func() any {
		return &mapIter{}
	},
}

func newMapCtx() *mapIter {
	ctx := mapCtxPool.Get().(*mapIter)
	ctx.Iter = hiter{}

	return ctx
}

func freeMapCtx(ctx *mapIter) {
	mapCtxPool.Put(ctx)
}
