package decoder

import (
	"fmt"
	"reflect"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type intDecoder struct {
	typ        reflect.Type
	kind       reflect.Kind
	structName string
	fieldName  string
}

func newIntDecoder(typ reflect.Type, structName, fieldName string) *intDecoder {
	return &intDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *intDecoder) typeError(buf []byte, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   d.typ,
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

func (d *intDecoder) parseInt(b []byte) (int64, error) {
	return parseInt64(b)
}

var pow10i64 = [...]uint64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
}

func parseInt64(buf []byte) (int64, error) {
	if len(buf) == 0 {
		return 0, fmt.Errorf("invalid integer")
	}

	negative := buf[0] == '-'
	if negative {
		buf = buf[1:]
		if len(buf) == 0 {
			return 0, fmt.Errorf("invalid integer")
		}
	}

	for len(buf) > 0 && buf[0] == '0' {
		buf = buf[1:]
	}
	if len(buf) == 0 {
		return 0, nil
	}
	if len(buf) > len(pow10i64) {
		return 0, fmt.Errorf("integer overflow")
	}

	var value uint64
	for i, c := range buf {
		if !numTable[c] {
			return 0, fmt.Errorf("invalid integer")
		}
		value += uint64(c-'0') * pow10i64[len(buf)-i-1]
	}

	limit := uint64(^uint64(0) >> 1)
	if negative {
		limit++
	}
	if value > limit {
		return 0, fmt.Errorf("integer overflow")
	}
	if negative {
		if value == uint64(1)<<63 {
			return -1 << 63, nil
		}
		return -int64(value), nil
	}
	return int64(value), nil
}

var (
	numTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
)

var (
	numZeroBuf = []byte{'0'}
)

func (d *intDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
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

func (d *intDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c

	return d.processBytes(bytes, cursor, rv)
}

func (d *intDecoder) processBytes(bytes []byte, cursor int64, rv reflect.Value) (int64, error) {
	i64, err := d.parseInt(bytes)
	if err != nil {
		return 0, d.typeError(bytes, cursor)
	}

	if rv.OverflowInt(i64) {
		return 0, errors.ErrOverflow(i64, rv.Type().Kind().String())
	}

	rv.SetInt(i64)

	return cursor, nil
}

func readInt(buf []byte, cursor int64) (int, int64, error) {
	start := cursor + 2
	bytes, end, err := readIntegerBytes(buf, cursor)
	if err != nil {
		return 0, cursor, err
	}
	value64, err := parseInt64(bytes)
	if err != nil || value64 < 0 || int64(int(value64)) != value64 {
		return 0, cursor, errors.ErrSyntax("php-serialize: invalid array index", start)
	}
	return int(value64), end, nil
}
