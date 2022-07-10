package encoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

// will need to get type message at marshal time, slow path.
// should avoid interface for performance thinking.
func compileInterface(rt *runtime.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		v := *(*any)(unsafe.Pointer(&emptyInterface{
			typ: rt,
			ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
		}))

		return reflectInterfaceValue(ctx, b, reflect.ValueOf(v), p)
	}, nil
}

func reflectInterfaceValue(ctx *Ctx, b []byte, rv reflect.Value, p uintptr) ([]byte, error) {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				return appendNilBytes(b), nil
			}
			rv = rv.Elem()
		default:
			break LOOP
		}
	}

	value := rv.Interface()

	v := *(*emptyInterface)(unsafe.Pointer(&value))
	pp := uintptr(v.ptr)

	// simple type
	switch v.typ.Kind() {
	case reflect.Bool:
		return encodeBool(ctx, b, pp)
	case reflect.Uint8:
		return encodeUint8(ctx, b, pp)
	case reflect.Uint16:
		return encodeUint16(ctx, b, pp)
	case reflect.Uint32:
		return encodeUint32(ctx, b, pp)
	case reflect.Uint64:
		return encodeUint64(ctx, b, pp)
	case reflect.Uint:
		return encodeUint(ctx, b, pp)
	case reflect.Int8:
		return encodeInt8(ctx, b, pp)
	case reflect.Int16:
		return encodeInt16(ctx, b, pp)
	case reflect.Int32:
		return encodeInt32(ctx, b, pp)
	case reflect.Int64:
		return encodeInt64(ctx, b, pp)
	case reflect.Int:
		return encodeInt(ctx, b, pp)
	case reflect.Float32:
		return encodeFloat32(ctx, b, pp)
	case reflect.Float64:
		return encodeFloat64(ctx, b, pp)
	case reflect.String:
		return encodeString(ctx, b, pp)
	}

	// if rv.Type().Kind() == reflect.Slice {
	// 	fmt.Println(orv.Type())
	// 	fmt.Println(orv.Elem().Type())
	// 	fmt.Println(v.typ.String())
	// 	value = *(*any)(unsafe.Pointer(&emptyInterface{typ: v.typ, ptr: unsafe.Pointer(p)}))
	// fmt.Println("php", uintptr(unsafe.Pointer(&v.typ)), rt, rt.Elem())
	//
	// fmt.Println("slice element", rv.Type().Elem())
	// fmt.Println(runtime.IfaceIndir(runtime.Type2RType(rv.Type().Elem())))
	// }

	// v = *(*emptyInterface)(unsafe.Pointer(&value))
	// pp = uintptr(v.ptr)

	switch rv.Kind() {
	case reflect.Slice:
		return reflectSlice(ctx, b, rv, p)
	case reflect.Map:
		return reflectMap(ctx, b, rv)
	case reflect.Struct:
		return reflectStruct(ctx, b, rv, pp)
	}

	return b, &UnsupportedInterfaceTypeError{rv.Type()}
}

func compileInterfaceAsString(rt *runtime.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		v := *(*any)(unsafe.Pointer(&emptyInterface{
			typ: rt,
			ptr: unsafe.Pointer(p),
		}))

		return reflectInterfaceValueAsString(ctx, b, reflect.ValueOf(v))
	}, nil
}

func reflectInterfaceValueAsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				return appendNilBytes(b), nil
			}
			rv = rv.Elem()
		default:
			break LOOP
		}
	}

	value := rv.Interface()

	v := *(*emptyInterface)(unsafe.Pointer(&value))
	p := uintptr(v.ptr)

	// simple type
	switch v.typ.Kind() {
	case reflect.Bool:
		return encodeBoolAsString(ctx, b, p)
	case reflect.Uint8:
		return encodeUint8AsString(ctx, b, p)
	case reflect.Uint16:
		return encodeUint16AsString(ctx, b, p)
	case reflect.Uint32:
		return encodeUint32AsString(ctx, b, p)
	case reflect.Uint64:
		return encodeUint64AsString(ctx, b, p)
	case reflect.Uint:
		return encodeUintAsString(ctx, b, p)
	case reflect.Int8:
		return encodeInt8AsString(ctx, b, p)
	case reflect.Int16:
		return encodeInt16AsString(ctx, b, p)
	case reflect.Int32:
		return encodeInt32AsString(ctx, b, p)
	case reflect.Int64:
		return encodeInt64AsString(ctx, b, p)
	case reflect.Int:
		return encodeIntAsString(ctx, b, p)
	case reflect.Float32:
		return encodeFloat32AsString(ctx, b, p)
	case reflect.Float64:
		return encodeFloat64AsString(ctx, b, p)
	}

	// slice, map and struct as interface are not supported yet.
	return b, fmt.Errorf("failed to encode %s as string", rv.Kind())
}
