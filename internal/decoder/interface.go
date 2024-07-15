package decoder

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type interfaceDecoder struct {
	typ              reflect.Type
	structName       string
	fieldName        string
	sliceDecoder     *sliceDecoder
	mapArrayDecoder  *mapDecoder
	mapClassDecoder  *mapDecoder
	floatDecoder     *floatDecoder
	stringDecoder    *stringDecoder
	intDecode        *intDecoder
	mapAnyKeyDecoder *mapKeyDecoder
}

func newEmptyInterfaceDecoder(structName, fieldName string) *interfaceDecoder {
	ifaceDecoder := &interfaceDecoder{
		typ:           emptyInterfaceType,
		structName:    structName,
		fieldName:     fieldName,
		floatDecoder:  newFloatDecoder(structName, fieldName),
		intDecode:     newIntDecoder(interfaceIntType, structName, fieldName),
		stringDecoder: newStringDecoder(structName, fieldName),
	}

	ifaceDecoder.mapAnyKeyDecoder = newInterfaceMapKeyDecoder(
		newIntDecoder(interfaceIntType, structName, fieldName),
		ifaceDecoder.stringDecoder)

	ifaceDecoder.sliceDecoder = newSliceDecoder(
		ifaceDecoder,
		emptyInterfaceType,
		emptyInterfaceType.Size(),
		structName, fieldName,
	)

	ifaceDecoder.mapClassDecoder = newMapDecoder(
		interfaceClassMapType,
		stringType,
		ifaceDecoder.stringDecoder,
		interfaceClassMapType.Elem(),
		ifaceDecoder,
		structName,
		fieldName,
	)

	ifaceDecoder.mapArrayDecoder = newMapDecoder(
		interfaceMapType,
		emptyInterfaceType,
		ifaceDecoder.mapAnyKeyDecoder,
		interfaceMapType.Elem(),
		ifaceDecoder,
		structName,
		fieldName,
	)
	return ifaceDecoder
}

func newInterfaceDecoder(typ reflect.Type, structName, fieldName string) *interfaceDecoder {
	emptyIfaceDecoder := newEmptyInterfaceDecoder(structName, fieldName)
	stringDecoder := newStringDecoder(structName, fieldName)
	return &interfaceDecoder{
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
		sliceDecoder: newSliceDecoder(
			emptyIfaceDecoder,
			emptyInterfaceType,
			emptyInterfaceType.Size(),
			structName, fieldName,
		),
		mapArrayDecoder: newMapDecoder(
			interfaceMapType,
			emptyInterfaceType,
			emptyIfaceDecoder.mapAnyKeyDecoder,
			interfaceMapType.Elem(),
			emptyIfaceDecoder,
			structName,
			fieldName,
		),
		floatDecoder:     emptyIfaceDecoder.floatDecoder,
		stringDecoder:    stringDecoder,
		intDecode:        emptyIfaceDecoder.intDecode,
		mapClassDecoder:  emptyIfaceDecoder.mapClassDecoder,
		mapAnyKeyDecoder: emptyIfaceDecoder.mapAnyKeyDecoder,
	}
}

var (
	stringType            = reflect.TypeOf((*string)(nil)).Elem()
	emptyInterfaceType    = reflect.TypeOf((*any)(nil)).Elem()
	interfaceMapType      = reflect.TypeOf((*map[any]any)(nil)).Elem()
	interfaceClassMapType = reflect.TypeOf((*map[string]any)(nil)).Elem()
	interfaceIntType      = reflect.TypeOf((*int64)(nil)).Elem()
)

func decodePHPUnmarshaler(buf []byte, cursor, depth int64, unmarshaler Unmarshaler, rv reflect.Value) (int64, error) {
	start := cursor
	end, err := skipValue(buf, cursor, depth)
	if err != nil {
		return 0, err
	}
	src := buf[start:end]
	if bytes.Equal(src, nullbytes) {
		rv.SetZero()
		return end, nil
	}
	if s, ok := unquoteBytes(src); ok {
		src = s
	}
	if err := unmarshaler.UnmarshalPHP(src); err != nil {
		return 0, err
	}
	return end, nil
}

type emptyInterface struct {
	typ uintptr // type ID
	ptr unsafe.Pointer
}

func (d *interfaceDecoder) errUnmarshalType(typ reflect.Type, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  typ.String(),
		Type:   typ,
		Offset: offset,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}

func (d *interfaceDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	if rv.NumMethod() > 0 && rv.CanInterface() {
		if u, ok := rv.Interface().(Unmarshaler); ok {
			return decodePHPUnmarshaler(buf, cursor, depth, u, rv)
		}
		if buf[cursor] == 'N' {
			if err := validateNull(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 2
			rv.SetZero()
			return cursor, nil
		}
		return 0, d.errUnmarshalType(rv.Type(), cursor)
	}

	if rv.Type().NumMethod() == 0 {
		// concrete type is empty interface
		return d.decodeEmptyInterface(ctx, cursor, depth, rv)
	}
	if rv.Type().Kind() == reflect.Ptr && rv.Type().Elem() == d.typ || rv.Type().Kind() != reflect.Ptr {
		return d.decodeEmptyInterface(ctx, cursor, depth, rv)
	}
	if buf[cursor] == 'N' {
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		rv.SetZero()
		return cursor, nil
	}
	decoder, err := CompileToGetDecoder(rv.Type())
	if err != nil {
		return 0, err
	}
	return decoder.Decode(ctx, cursor, depth, rv)
}

func (d *interfaceDecoder) decodeEmptyInterface(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	switch buf[cursor] {
	case 'O':
		var v map[string]any
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.mapClassDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	case 'a':
		var v map[any]any
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.mapArrayDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	case 'd':
		var v float64
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.floatDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	case 's':
		cursor++
		b, end, err := readString(buf, cursor)
		if err != nil {
			return 0, err
		}
		rv.Set(reflect.ValueOf(string(b)))
		return end, nil
	case 'i':
		var v int64
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.intDecode.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	case 'b':
		v, err := readBool(buf, cursor)
		if err != nil {
			return 0, err
		}
		if v {
			rv.Set(reflect.ValueOf(true))
		} else {
			rv.Set(reflect.ValueOf(false))
		}
		return cursor + 4, nil
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		rv.SetZero()
		return cursor, nil
	}
	return cursor, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
}

type mapKeyDecoder struct {
	strDecoder *stringDecoder
	intDecoder *intDecoder
}

func (d *mapKeyDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf

	switch buf[cursor] {
	case 's':
		var v string
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.strDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	// string key
	case 'i':
		var v int64
		ptr := reflect.ValueOf(&v).Elem()
		cursor, err := d.intDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		rv.Set(ptr)
		return cursor, nil
	default:
		return 0, errors.ErrExpected("array key", cursor)
	}
}

func newInterfaceMapKeyDecoder(intDecoder *intDecoder, stringDecoder *stringDecoder) *mapKeyDecoder {
	return &mapKeyDecoder{intDecoder: intDecoder, strDecoder: stringDecoder}
}
