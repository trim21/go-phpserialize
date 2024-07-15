package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type structFieldSet struct {
	dec      Decoder
	offset   uintptr
	fieldIdx int
	key      string
	err      error
}

type structDecoder struct {
	fieldMap      map[string]*structFieldSet
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
}

func newStructDecoder(structName, fieldName string, fieldMap map[string]*structFieldSet) *structDecoder {
	return &structDecoder{
		fieldMap:      fieldMap,
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
	}
}

// TODO: this can be optimized for small size struct
func decodeKey(d *structDecoder, buf []byte, cursor int64) (int64, *structFieldSet, error) {
	key, c, err := d.stringDecoder.decodeByte(buf, cursor)
	if err != nil {
		return 0, nil, err
	}
	cursor = c

	// go compiler will not escape key
	field, exists := d.fieldMap[string(key)]
	if !exists {
		return cursor, nil, nil
	}

	return cursor, field, nil
}

func (d *structDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}
	buflen := int64(len(buf))
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	switch char(b, cursor) {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
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
		cursor++
		if buf[cursor] != ':' {
			return 0, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
		}
	default:
		return 0, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
	}

	// skip  :${length}:
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return cursor, err
	}
	cursor = end
	if buf[cursor] != '{' {
		return 0, errors.ErrInvalidBeginningOfArray(char(b, cursor), cursor)
	}

	cursor++
	if buf[cursor] == '}' {
		cursor++
		return cursor, nil
	}

	for {
		c, field, err := decodeKey(d, buf, cursor)
		if err != nil {
			return 0, err
		}

		cursor = c

		// cursor++
		if cursor >= buflen {
			return 0, errors.ErrExpected("object value after colon", cursor)
		}
		if field != nil {
			if field.err != nil {
				return 0, field.err
			}
			c, err := field.dec.Decode(ctx, cursor, depth, rv.Field(field.fieldIdx))
			if err != nil {
				return 0, err
			}
			cursor = c
		} else {
			c, err := skipValue(buf, cursor, depth)
			if err != nil {
				return 0, err
			}
			cursor = c
		}

		if char(b, cursor) == '}' {
			cursor++
			return cursor, nil
		}
	}
}
