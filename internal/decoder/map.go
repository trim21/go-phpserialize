package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type mapDecoder struct {
	mapType                 *runtime.Type
	keyType                 *runtime.Type
	valueType               *runtime.Type
	canUseAssignFaststrType bool
	keyDecoder              Decoder
	valueDecoder            Decoder
	structName              string
	fieldName               string
}

func newMapDecoder(mapType *runtime.Type, keyType *runtime.Type, keyDec Decoder, valueType *runtime.Type, valueDec Decoder, structName, fieldName string) *mapDecoder {
	return &mapDecoder{
		mapType:                 mapType,
		keyDecoder:              keyDec,
		keyType:                 keyType,
		canUseAssignFaststrType: canUseAssignFaststrType(keyType, valueType),
		valueType:               valueType,
		valueDecoder:            valueDec,
		structName:              structName,
		fieldName:               fieldName,
	}
}

const (
	mapMaxElemSize = 128
)

// See detail: https://github.com/goccy/go-json/pull/283
func canUseAssignFaststrType(key *runtime.Type, value *runtime.Type) bool {
	indirectElem := value.Size() > mapMaxElemSize
	if indirectElem {
		return false
	}
	return key.Kind() == reflect.String
}

func (d *mapDecoder) mapassign(t *runtime.Type, m, k, v unsafe.Pointer) {
	if d.canUseAssignFaststrType {
		mapV := runtime.MapAssignFastStr(t, m, *(*string)(k))
		typedmemmove(d.valueType, mapV, v)
	} else {
		runtime.MapAssign(t, m, k, v)
	}
}

func (d *mapDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
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
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = nil
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
		return 0, errors.ErrExpected("{ character for map value", cursor)
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
	if mapValue == nil {
		mapValue = runtime.MakeMap(d.mapType, int(l))
	}

	cursor++
	if buf[cursor] == '}' {
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
		cursor++
		return cursor, nil
	}

	for {
		k := unsafe_New(d.keyType)
		keyCursor, err := d.keyDecoder.Decode(ctx, cursor, depth, k)
		if err != nil {
			return 0, err
		}
		cursor = keyCursor
		v := unsafe_New(d.valueType)
		valueCursor, err := d.valueDecoder.Decode(ctx, cursor, depth, v)
		if err != nil {
			return 0, err
		}

		d.mapassign(d.mapType, mapValue, k, v)
		cursor = valueCursor
		if buf[cursor] == '}' {
			**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
			cursor++
			return cursor, nil
		}
	}
}
