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
	return func(ctx *Ctx, p uintptr) error {
		v := *(*any)(unsafe.Pointer(&emptyInterface{
			typ: rt,
			ptr: unsafe.Pointer(p),
		}))

		return reflectInterfaceValue(ctx, reflect.ValueOf(v), p)
	}, nil
}

func reflectInterfaceValue(ctx *Ctx, rv reflect.Value, p uintptr) error {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				appendNil(ctx)
				return nil
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
		return encodeBool(ctx, pp)
	case reflect.Uint8:
		return encodeUint8(ctx, pp)
	case reflect.Uint16:
		return encodeUint16(ctx, pp)
	case reflect.Uint32:
		return encodeUint32(ctx, pp)
	case reflect.Uint64:
		return encodeUint64(ctx, pp)
	case reflect.Uint:
		return encodeUint(ctx, pp)
	case reflect.Int8:
		return encodeInt8(ctx, pp)
	case reflect.Int16:
		return encodeInt16(ctx, pp)
	case reflect.Int32:
		return encodeInt32(ctx, pp)
	case reflect.Int64:
		return encodeInt64(ctx, pp)
	case reflect.Int:
		return encodeInt(ctx, pp)
	case reflect.Float32:
		return encodeFloat32(ctx, pp)
	case reflect.Float64:
		return encodeFloat64(ctx, pp)
	case reflect.String:
		return encodeStringVariable(ctx, pp)
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
		return reflectSlice(ctx, rv, p)
	case reflect.Map:
		return reflectMap(ctx, rv)
	case reflect.Struct:
		return reflectStruct(ctx, rv, pp)
	}

	return &UnsupportedInterfaceTypeError{rv.Type()}
}

func compileInterfaceAsString(rt *runtime.Type) (encoder, error) {
	return func(ctx *Ctx, p uintptr) error {
		v := *(*any)(unsafe.Pointer(&emptyInterface{
			typ: rt,
			ptr: unsafe.Pointer(p),
		}))

		return reflectInterfaceValueAsString(ctx, reflect.ValueOf(v))
	}, nil
}

func reflectInterfaceValueAsString(ctx *Ctx, rv reflect.Value) error {
LOOP:
	for {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() || rv.IsZero() {
				appendNil(ctx)
				return nil
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
		return encodeBoolAsString(ctx, p)
	case reflect.Uint8:
		return encodeUint8AsString(ctx, p)
	case reflect.Uint16:
		return encodeUint16AsString(ctx, p)
	case reflect.Uint32:
		return encodeUint32AsString(ctx, p)
	case reflect.Uint64:
		return encodeUint64AsString(ctx, p)
	case reflect.Uint:
		return encodeUintAsString(ctx, p)
	case reflect.Int8:
		return encodeInt8AsString(ctx, p)
	case reflect.Int16:
		return encodeInt16AsString(ctx, p)
	case reflect.Int32:
		return encodeInt32AsString(ctx, p)
	case reflect.Int64:
		return encodeInt64AsString(ctx, p)
	case reflect.Int:
		return encodeIntAsString(ctx, p)
	case reflect.Float32:
		return encodeFloat32AsString(ctx, p)
	case reflect.Float64:
		return encodeFloat64AsString(ctx, p)
	}

	// slice, map and struct as interface are not supported yet.
	return fmt.Errorf("failed to encode %s as string", rv.Kind())
}
