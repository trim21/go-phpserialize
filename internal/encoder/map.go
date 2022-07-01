package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
	"github.com/gookit/goutil/dump"
)

func compileMap(typ reflect.Type, rv reflect.Value) (encoder, error) {
	// for map[int]string, keyType is int, valueType is string
	keyType := typ.Key()
	valueType := typ.Elem()

	flag := (*(*rValue)(unsafe.Pointer(&rv))).flag

	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, &UnsupportedTypeAsMapKeyError{Type: keyType}
	}

	if valueType.Kind() == reflect.Interface {
		// interface slow path
		return compileMapWithInterfaceValue(typ, rv, keyType, valueType)
	}

	keyEncoder, err := compile(keyType, reflect.New(keyType).Elem())
	if err != nil {
		return nil, err
	}

	valueEncoder, err := compile(valueType, reflect.New(valueType).Elem())
	if err != nil {
		return nil, err
	}

	return func(buf *buffer, p uintptr) error {
		rv := reflectValueMapFromPtr(typ, p, flag)
		for _, key := range rv.MapKeys() {
			err := keyEncoder(buf, reflectValueToLocal(key).ptr)
			if err != nil {
				return err
			}

			v := rv.MapIndex(key)

			err = valueEncoder(buf, reflectValueToLocal(v).ptr)
			if err != nil {
				return err
			}
		}

		if rv.IsNil() {
			appendArrayBegin(buf, 0)
			buf.b = append(buf.b, '}')

			return nil
		}

		mapLen := rv.Len()
		appendArrayBegin(buf, int64(mapLen))

		if mapLen == 0 {
			buf.b = append(buf.b, '}')
			return nil
		}

		smr := rv.MapRange()
		dump.P(smr)
		mr := *(**MapIter)(unsafe.Pointer(&smr))
		dump.P(mr)
		for mr.Next() {
			err := keyEncoder(buf, mr.Key())
			if err != nil {
				return err
			}

			err = valueEncoder(buf, mr.Value())
			if err != nil {
				return err
			}
		}
		buf.b = append(buf.b, '}')
		return nil
	}, nil
}

func compileMapWithInterfaceValue(typ reflect.Type, rv reflect.Value, keyType, valueType reflect.Type) (encoder, error) {
	keyEncoder, err := compile(keyType, rv)
	if err != nil {
		return nil, err
	}

	flag := (*(*rValue)(unsafe.Pointer(&rv))).flag

	valueEncoder, err := compileInterface(valueType)
	if err != nil {
		return nil, err
	}

	return func(buf *buffer, p uintptr) error {
		rv := reflectValueMapFromPtr(typ, p, flag)

		if rv.IsNil() {
			appendArrayBegin(buf, 0)

			return nil
		}

		appendArrayBegin(buf, int64(rv.Len()))

		mr := rv.MapRange()
		for mr.Next() {
			// HINT: This is very likely to break after new go version.
			p := mr.Key()
			v := (*rValue)(unsafe.Pointer(&p))

			err := keyEncoder(buf, v.ptr)
			if err != nil {
				return err
			}

			p = mr.Value()
			v = (*rValue)(unsafe.Pointer(&p))
			err = valueEncoder(buf, v.ptr)
			if err != nil {
				return err
			}
		}
		buf.b = append(buf.b, '}')
		return nil
	}, nil
}
