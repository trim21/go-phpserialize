package encoder

import (
	"fmt"

	"github.com/goccy/go-reflect"
)

// will need to get type message at marshal time, slow path.
// should avoid interface for performance thinking.
func compileInterface(typ reflect.Type) (encoder, error) {
	return reflectInterface, nil
}

func reflectInterface(ctx *Ctx, p uintptr) error {
	rv := reflectValueMapFromPtr(reflect.TypeOf(any(1)), p, flag(reflect.Interface))
	return reflectInterfaceValue(ctx, rv)
}

// // slow path of encoding struct
// func encodeStruct(buf *Ctx, rv reflect.Value) error {
// 	typeID := uintptr(unsafe.Pointer(rv.Type()))
//
// 	appendArrayBegin(buf, int64(rv.NumField()))
//
// 	for i := 0; i < rv.NumField(); i++ {
// 		field := rv.Field(i)
// 		p := field.UnsafeAddr()
// 		switch field.Kind() {
// 		case reflect.Uint8:
// 			appendUint(buf, uint64(*(*uint)(p)))
// 			field.Uint()
//
// 		}
// 	}
//
// 	buf.b = append(buf.b, '}')
// 	return nil
// }

func reflectInterfaceValue(ctx *Ctx, rv reflect.Value) error {
	// simple type
	switch v := rv.Interface().(type) {
	case bool:
		appendBool(ctx, v)
		return nil
	case uint8:
		appendUint(ctx, uint64(v))
		return nil
	case uint16:
		appendUint(ctx, uint64(v))
		return nil
	case uint32:
		appendUint(ctx, uint64(v))
		return nil
	case uint:
		appendUint(ctx, uint64(v))
		return nil
	case uint64:
		appendUint(ctx, v)
		return nil

	case int8:
		appendInt(ctx, int64(v))
		return nil
	case int16:
		appendInt(ctx, int64(v))
		return nil
	case int32:
		appendInt(ctx, int64(v))
		return nil
	case int:
		appendInt(ctx, int64(v))
		return nil
	case int64:
		appendInt(ctx, v)
		return nil
	case float32:
		appendFloat32(ctx, v)
		return nil
	case float64:
		appendFloat64(ctx, v)
		return nil
	case string:
		appendString(ctx, v)
		return nil
	}

	if rv.Kind() != reflect.Interface {
		switch rv.Kind() {
		case reflect.Slice:
			return reflectSlice(ctx, rv)
		case reflect.Map:
			return reflectMap(ctx, rv)
		case reflect.Struct:
			return reflectStruct(ctx, rv)
		default:
			fmt.Println("un-expected interface underlying type", rv.Elem().Kind())
		}
	}

	switch rv.Elem().Kind() {
	case reflect.Slice:
		return reflectSlice(ctx, rv.Elem())
	case reflect.Map:
		return reflectMap(ctx, rv.Elem())
	case reflect.Struct:
		return reflectStruct(ctx, rv.Elem())
	default:
		fmt.Println("un-expected interface underlying type", rv.Elem().Kind())
	}

	// slice, map and struct as interface are not supported yet.
	return &UnsupportedInterfaceTypeError{rv.Type()}
}

func reflectMap(ctx *Ctx, rv reflect.Value) error {
	if rv.Type().Key().Kind() == reflect.Interface {
		// interface key
	}

	return &UnsupportedInterfaceTypeError{rv.Type()}
}

func reflectSlice(ctx *Ctx, rv reflect.Value) error {
	l := rv.Len()
	appendArrayBegin(ctx, int64(l))
	for i := 0; i < l; i++ {
		appendInt(ctx, int64(i))
		err := reflectInterfaceValue(ctx, rv.Index(i))
		if err != nil {
			return err
		}
	}

	ctx.b = append(ctx.b, '}')

	return nil
}

func reflectStruct(ctx *Ctx, rv reflect.Value) error {
	rt := rv.Type()
	appendArrayBegin(ctx, int64(rv.NumField()))

	for i := 0; i < rv.NumField(); i++ {
		appendString(ctx, getFieldName(rt.Field(i)))
		err := reflectInterfaceValue(ctx, rv.Field(i))
		if err != nil {
			return err
		}
	}

	ctx.b = append(ctx.b, '}')

	return nil
}
