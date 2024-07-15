package decoder

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type floatDecoder struct {
	structName string
	fieldName  string
}

func newFloatDecoder(structName, fieldName string) *floatDecoder {
	return &floatDecoder{structName: structName, fieldName: fieldName}
}

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

func (d *floatDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	bytes, cursor, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}

	if buf[cursor] != ';' {
		return cursor, errors.ErrExpected("float end with ';'", cursor)
	}
	cursor++

	if bytes == nil {
		return cursor, nil
	}

	return d.processBytes(bytes, cursor, rv)
}

func (d *floatDecoder) processBytes(bytes []byte, cursor int64, rv reflect.Value) (int64, error) {
	s := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.ErrSyntax(err.Error(), cursor)
	}

	if rv.OverflowFloat(f64) {
		return 0, errors.ErrOverflow(f64, rv.Type().Kind().String())
	}

	rv.SetFloat(f64)

	return cursor, nil
}
