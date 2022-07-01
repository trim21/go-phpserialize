package encoder

import (
	"github.com/goccy/go-reflect"
)

// will need to get type message at marshal time, slow path.
// should avoid interface for performance thinking.
func compileInterface(typ reflect.Type) (encoder, error) {
	typ.Kind()
	return func(buf *buffer, p uintptr) error {
		rv := reflectValueFromPtr(typ, p)

		// simple type
		switch v := rv.Interface().(type) {
		case bool:
			appendBool(buf, v)
		case uint8:
			appendUint(buf, uint64(v))
		case uint16:
			appendUint(buf, uint64(v))
		case uint32:
			appendUint(buf, uint64(v))
		case uint:
			appendUint(buf, uint64(v))
		case uint64:
			appendUint(buf, v)

		case int8:
			appendInt(buf, int64(v))
		case int16:
			appendInt(buf, int64(v))
		case int32:
			appendInt(buf, int64(v))
		case int:
			appendInt(buf, int64(v))
		case int64:
			appendInt(buf, v)
		case string:
			appendString(buf, v)
		}

		// slice, map and struct will be handled in `compile`
		return &UnsupportedInterfaceTypeError{rv.Type()}
	}, nil
}

// // slow path of encoding struct
// func encodeStruct(buf *buffer, rv reflect.Value) error {
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
