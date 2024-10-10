package encoder

import (
	"fmt"
	"reflect"
)

func compilePtr(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	switch rt.Elem().Kind() {
	case reflect.Ptr:
		return nil, fmt.Errorf("encoding nested ptr is not supported *%s", rt.Elem().String())
	case reflect.Bool:
		return deRefNilEncoder(encodeBool), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return deRefNilEncoder(encodeUint), nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return deRefNilEncoder(encodeInt), nil
	case reflect.Float32:
		return deRefNilEncoder(encodeFloat32), nil
	case reflect.Float64:
		return deRefNilEncoder(encodeFloat64), nil
	case reflect.String:
		return deRefNilEncoder(encodeString), nil
	case reflect.Interface:
		return compileInterface(rt.Elem())
	case reflect.Map:
		enc, err := compileMap(rt.Elem(), seen)
		return deRefNilEncoder(enc), err
	case reflect.Struct:
		enc, err := compileStruct(rt.Elem(), seen)
		return checkStructRecursiveEncoder(enc), err
	}

	enc, err := compile(rt.Elem(), seen)
	if err != nil {
		return nil, err
	}

	return deRefNilEncoder(enc), nil
}

func deRefNilEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		return enc(ctx, b, rv.Elem())
	}
}

func checkStructRecursiveEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		ctx.StackDepth++
		if ctx.StackDepth > 1000 {
			_, seen := ctx.Seen[rv.UnsafePointer()]
			if seen {
				return b, fmt.Errorf("php: try to encode recursive object %v", rv.Interface())
			}
		}

		b, err := enc(ctx, b, rv.Elem())
		if err != nil {
			return b, err
		}

		if ctx.StackDepth > 1000 {
			delete(ctx.Seen, rv.UnsafePointer())
		}
		ctx.StackDepth--

		return b, nil
	}
}

func checkRecursiveEncoder(enc encoder) encoder {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		if rv.IsNil() {
			return appendNull(b), nil
		}

		ctx.StackDepth++
		if ctx.StackDepth > 1000 {
			_, seen := ctx.Seen[rv.UnsafePointer()]
			if seen {
				return b, fmt.Errorf("php: try to encode recursive object %v", rv.Interface())
			}

			ctx.Seen[rv.UnsafePointer()] = empty{}
		}

		b, err := enc(ctx, b, rv)
		if err != nil {
			return b, err
		}

		if ctx.StackDepth > 1000 {
			delete(ctx.Seen, rv.UnsafePointer())
		}
		ctx.StackDepth--

		return b, nil
	}
}
