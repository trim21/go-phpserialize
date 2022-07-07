package decoder

import (
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type floatDecoder struct {
	op         func(unsafe.Pointer, float64)
	structName string
	fieldName  string
}

func newFloatDecoder(structName, fieldName string, op func(unsafe.Pointer, float64)) *floatDecoder {
	return &floatDecoder{op: op, structName: structName, fieldName: fieldName}
}

var (
	floatTable = [256]bool{
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
		'.': true,
		'e': true,
		'E': true,
		'+': true,
		'-': true,
	}

	validEndNumberChar = [256]bool{
		nul:  true,
		' ':  true,
		'\t': true,
		'\r': true,
		'\n': true,
		',':  true,
		':':  true,
		'}':  true,
		']':  true,
	}
)

func (d *floatDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil

	case 'd':
		break
	default:
		return nil, cursor, errors.ErrExpected("float start with 'd' or 'N'", cursor)
	}

	cursor++
	if buf[cursor] != ':' {
		return nil, cursor, errors.ErrExpected("float start with 'd:'", cursor)
	}
	// cursor++

	start := cursor + 1
	for {
		cursor++
		if buf[cursor] == ';' {
			break
		}
	}

	num := buf[start:cursor]
	return num, cursor, nil
}

func (d *floatDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	bytes, cursor, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}

	if buf[cursor] != ';' {
		return cursor, errors.ErrExpected("float end with ';'", cursor)
	}

	if bytes == nil {
		return cursor, nil
	}

	s := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.ErrSyntax(err.Error(), cursor)
	}

	d.op(p, f64)

	return cursor + 1, nil
}
