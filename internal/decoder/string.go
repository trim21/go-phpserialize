package decoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type stringDecoder struct {
	structName string
	fieldName  string
}

func newStringDecoder(structName, fieldName string) *stringDecoder {
	return &stringDecoder{
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *stringDecoder) errUnmarshalType(typeName string, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  typeName,
		Type:   reflect.TypeOf(""),
		Offset: offset,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}

func (d *stringDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	rv.SetString(string(bytes))
	return cursor, nil
}

func (d *stringDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	switch buf[cursor] {
	case 'n':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil
	case 'b':
		return nil, 0, d.errUnmarshalType("bool", cursor)
	case 'd':
		return nil, 0, d.errUnmarshalType("float", cursor)
	case 's':
		cursor++
	case 'i':
		return nil, 0, d.errUnmarshalType("number", cursor)
		// read int as string
	default:
		return nil, 0, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
	}

	s, end, err := readString(buf, cursor)
	if err != nil {
		return nil, 0, err
	}

	return s, end, nil
}
