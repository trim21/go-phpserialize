package decoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type boolDecoder struct {
	structName string
	fieldName  string
}

func newBoolDecoder(structName, fieldName string) *boolDecoder {
	return &boolDecoder{structName: structName, fieldName: fieldName}
}

func (d *boolDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	if !hasByte(buf, cursor) {
		return 0, errors.ErrUnexpectedEnd("bool", cursor)
	}
	switch buf[cursor] {
	case 'b':
		value, err := readBool(buf, cursor)
		if err != nil {
			return 0, err
		}
		rv.SetBool(value)
		return cursor + 4, nil

	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	}

	return 0, errors.ErrUnexpectedEnd("bool", cursor)
}
