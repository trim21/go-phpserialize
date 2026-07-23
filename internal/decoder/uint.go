package decoder

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type uintDecoder struct {
	typ        reflect.Type
	kind       reflect.Kind
	structName string
	fieldName  string
}

func newUintDecoder(typ reflect.Type, structName, fieldName string) *uintDecoder {
	return &uintDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *uintDecoder) typeError(buf []byte, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   d.typ,
		Offset: offset,
	}
}

func (d *uintDecoder) parseUint(b []byte) (uint64, error) {
	return strconv.ParseUint(unsafeStr(b), 10, 64)
}

func (d *uintDecoder) decodeBytes(buf []byte, cursor int64) ([]byte, int64, error) {
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("integer", cursor)
	}
	if buf[cursor] == 'N' {
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		return nil, cursor + 2, nil
	}
	return readIntegerBytes(buf, cursor)
}

func (d *uintDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	bytes, c, err := d.decodeBytes(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c

	return d.processBytes(bytes, cursor, rv)
}

func (d *uintDecoder) processBytes(bytes []byte, cursor int64, rv reflect.Value) (int64, error) {
	u64, err := d.parseUint(bytes)
	if err != nil {
		return 0, d.typeError(bytes, cursor)
	}

	if rv.OverflowUint(u64) {
		return 0, errors.ErrOverflow(u64, rv.Type().Kind().String())
	}

	rv.SetUint(u64)

	return cursor, nil
}
