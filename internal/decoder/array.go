package decoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type arrayDecoder struct {
	elemType     reflect.Type
	size         uintptr
	valueDecoder Decoder
	alen         int
	structName   string
	fieldName    string
	zeroValue    reflect.Value
}

func newArrayDecoder(dec Decoder, elemType reflect.Type, alen int, structName, fieldName string) *arrayDecoder {
	zeroValue := reflect.Zero(elemType)
	return &arrayDecoder{
		valueDecoder: dec,
		elemType:     elemType,
		size:         elemType.Size(),
		alen:         alen,
		structName:   structName,
		fieldName:    fieldName,
		zeroValue:    zeroValue,
	}
}

func (d *arrayDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	case 'a':
		cursor++
		if buf[cursor] != ':' {
			return cursor, errors.ErrExpected("':' before array length", cursor)
		}

		// set zero value first, php array may skip some index
		rv.SetZero()

		cursor++
		if buf[cursor] == '0' {
			err := validateEmptyArray(buf, cursor)
			if err != nil {
				return cursor, err
			}
			return cursor + 4, nil
		}

		_, end, err := readLengthInt(buf, cursor-1)
		if err != nil {
			return cursor, err
		}
		cursor = end + 1

		idx := 0
		for {
			currentIndex, end, err := readInt(buf, cursor)
			if err != nil {
				return 0, err
			}

			idx = currentIndex
			cursor = end

			if idx < d.alen {
				c, err := d.valueDecoder.Decode(ctx, cursor, depth, rv.Index(idx))
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

			if buf[cursor] == '}' {
				cursor++
				return cursor, nil
			}
		}
	default:
		return 0, errors.ErrUnexpectedEnd("array", cursor)
	}
}
