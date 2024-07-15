package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type mapDecoder struct {
	mapType      reflect.Type
	keyType      reflect.Type
	valueType    reflect.Type
	keyDecoder   Decoder
	valueDecoder Decoder
	structName   string
	fieldName    string
}

func newMapDecoder(mapType reflect.Type, keyType reflect.Type, keyDec Decoder, valueType reflect.Type, valueDec Decoder, structName, fieldName string) *mapDecoder {
	return &mapDecoder{
		mapType:      mapType,
		keyDecoder:   keyDec,
		keyType:      keyType,
		valueType:    valueType,
		valueDecoder: valueDec,
		structName:   structName,
		fieldName:    fieldName,
	}
}

func (d *mapDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	buflen := int64(len(buf))
	if buflen < 2 {
		return 0, errors.ErrExpected("{} for map", cursor)
	}
	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		rv.SetZero()
		return cursor, nil
	case 'O':
		// O:8:"stdClass":1:{s:1:"a";s:1:"q";}
		end, err := skipClassName(buf, cursor)
		if err != nil {
			return cursor, err
		}
		cursor = end
		fallthrough
	case 'a':
		// array case
		cursor++
	default:
		return 0, errors.ErrUnexpectedStart("map", buf, cursor)
	}

	l, end, err := readLength(buf, cursor)
	if err != nil {
		return 0, err
	}

	cursor = end
	if buf[cursor] != '{' {
		return 0, errors.ErrExpected("{ character for map value", cursor)
	}

	mapValue := *(*unsafe.Pointer)(p)

	// TODO
	rv := reflect.MakeMapWithSize(d.mapType, int(l))
	if mapValue == nil {
		mapValue = rv.UnsafePointer()
	}

	cursor++
	if buf[cursor] == '}' {
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
		cursor++
		return cursor, nil
	}

	for {
		k := reflect.New(d.keyType)
		keyCursor, err := d.keyDecoder.Decode(ctx, cursor, depth, k)
		if err != nil {
			return 0, err
		}
		cursor = keyCursor
		v := reflect.New(d.valueType)
		valueCursor, err := d.valueDecoder.Decode(ctx, cursor, depth, v)
		if err != nil {
			return 0, err
		}

		rv.SetMapIndex(k, v)
		cursor = valueCursor
		if buf[cursor] == '}' {
			cursor++
			return cursor, nil
		}
	}
}
