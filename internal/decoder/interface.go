package decoder

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type interfaceDecoder struct {
	typ           *runtime.Type
	structName    string
	fieldName     string
	sliceDecoder  *sliceDecoder
	mapDecoder    *mapDecoder
	floatDecoder  *floatDecoder
	stringDecoder *stringDecoder
	intDecode     *intDecoder
	mapKeyDecoder *mapKeyDecoder
}

func newEmptyInterfaceDecoder(structName, fieldName string) *interfaceDecoder {
	ifaceDecoder := &interfaceDecoder{
		typ:           emptyInterfaceType,
		structName:    structName,
		fieldName:     fieldName,
		floatDecoder:  newFloatDecoder(structName, fieldName, func(p unsafe.Pointer, v float64) { *(*any)(p) = v }),
		intDecode:     newIntDecoder(interfaceIntType, structName, fieldName, func(p unsafe.Pointer, v int64) { *(*any)(p) = v }),
		stringDecoder: newStringDecoder(structName, fieldName),
	}

	ifaceDecoder.mapKeyDecoder = newInterfaceMapKeyDecoder(
		newIntDecoder(interfaceIntType, structName, fieldName, func(p unsafe.Pointer, v int64) { *(*int64)(p) = v }),
		ifaceDecoder.stringDecoder)

	ifaceDecoder.sliceDecoder = newSliceDecoder(
		ifaceDecoder,
		emptyInterfaceType,
		emptyInterfaceType.Size(),
		structName, fieldName,
	)

	ifaceDecoder.mapDecoder = newMapDecoder(
		interfaceMapType,
		emptyInterfaceType,
		ifaceDecoder.mapKeyDecoder,
		interfaceMapType.Elem(),
		ifaceDecoder,
		structName,
		fieldName,
	)
	return ifaceDecoder
}

func newInterfaceDecoder(typ *runtime.Type, structName, fieldName string) *interfaceDecoder {
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
		mapDecoder: newMapDecoder(
			interfaceMapType,
			emptyInterfaceType,
			emptyIfaceDecoder.mapKeyDecoder,
			interfaceMapType.Elem(),
			emptyIfaceDecoder,
			structName,
			fieldName,
		),
		floatDecoder:  emptyIfaceDecoder.floatDecoder,
		stringDecoder: stringDecoder,
		intDecode:     emptyIfaceDecoder.intDecode,
		mapKeyDecoder: emptyIfaceDecoder.mapKeyDecoder,
	}
}

var (
	emptyInterfaceType = runtime.Type2RType(reflect.TypeOf((*any)(nil)).Elem())
	interfaceMapType   = runtime.Type2RType(reflect.TypeOf((*map[any]any)(nil)).Elem())
	interfaceIntType   = runtime.Type2RType(reflect.TypeOf((*int64)(nil))).Elem()
)

func decodePHPUnmarshaler(buf []byte, cursor, depth int64, unmarshaler Unmarshaler, p unsafe.Pointer) (int64, error) {
	start := cursor
	end, err := skipValue(buf, cursor, depth)
	if err != nil {
		return 0, err
	}
	src := buf[start:end]
	if bytes.Equal(src, nullbytes) {
		*(*unsafe.Pointer)(p) = nil
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
	typ *runtime.Type
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

func (d *interfaceDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf

	runtimeInterfaceValue := *(*any)(unsafe.Pointer(&emptyInterface{typ: d.typ, ptr: p}))
	rv := reflect.ValueOf(runtimeInterfaceValue)
	if rv.NumMethod() > 0 && rv.CanInterface() {
		if u, ok := rv.Interface().(Unmarshaler); ok {
			return decodePHPUnmarshaler(buf, cursor, depth, u, p)
		}
		if buf[cursor] == 'N' {
			if err := validateNull(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 2
			**(**any)(unsafe.Pointer(&p)) = nil
			return cursor, nil
		}
		return 0, d.errUnmarshalType(rv.Type(), cursor)
	}

	iface := rv.Interface()
	ifaceHeader := (*emptyInterface)(unsafe.Pointer(&iface))
	typ := ifaceHeader.typ
	if ifaceHeader.ptr == nil || d.typ == typ || typ == nil {
		// concrete type is empty interface
		return d.decodeEmptyInterface(ctx, cursor, depth, p)
	}
	if typ.Kind() == reflect.Ptr && typ.Elem() == d.typ || typ.Kind() != reflect.Ptr {
		return d.decodeEmptyInterface(ctx, cursor, depth, p)
	}
	if buf[cursor] == 'N' {
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		**(**any)(unsafe.Pointer(&p)) = nil
		return cursor, nil
	}
	decoder, err := CompileToGetDecoder(typ)
	if err != nil {
		return 0, err
	}
	return decoder.Decode(ctx, cursor, depth, ifaceHeader.ptr)
}

func (d *interfaceDecoder) decodeEmptyInterface(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	switch buf[cursor] {
	case 'a':
		var v map[any]any
		ptr := unsafe.Pointer(&v)
		cursor, err := d.mapDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		**(**any)(unsafe.Pointer(&p)) = v
		return cursor, nil
	case 'd': // with op
		return d.floatDecoder.Decode(ctx, cursor, depth, p)
	case 's':
		cursor++
		b, end, err := readString(buf, cursor)
		if err != nil {
			return 0, err
		}
		*(*any)(p) = string(b)
		return end, nil
	case 'i': // with op
		return d.intDecode.Decode(ctx, cursor, depth, p)
	case 'b':
		v, err := readBool(buf, cursor)
		if err != nil {
			return 0, err
		}

		if v {
			**(**any)(unsafe.Pointer(&p)) = true
		} else {
			**(**any)(unsafe.Pointer(&p)) = false
		}

		return cursor + 3, nil
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		**(**any)(unsafe.Pointer(&p)) = nil
		return cursor, nil
	}
	return cursor, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
}

type mapKeyDecoder struct {
	strDecoder *stringDecoder
	intDecoder *intDecoder
}

func (d *mapKeyDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf

	switch buf[cursor] {
	case 's':
		var v string
		ptr := unsafe.Pointer(&v)
		cursor, err := d.strDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		*(*any)(p) = v
		return cursor, nil
	// string key
	case 'i':
		var v int64
		ptr := unsafe.Pointer(&v)
		cursor, err := d.intDecoder.Decode(ctx, cursor, depth, ptr)
		if err != nil {
			return 0, err
		}
		*(*any)(p) = v
		return cursor, nil
	default:
		return 0, errors.ErrExpected("array key", cursor)
	}
}

func newInterfaceMapKeyDecoder(intDecoder *intDecoder, stringDecoder *stringDecoder) *mapKeyDecoder {
	return &mapKeyDecoder{intDecoder: intDecoder, strDecoder: stringDecoder}
}
