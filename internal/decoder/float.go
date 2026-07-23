package decoder

import (
	"math"
	"reflect"
	"strconv"

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
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("float", cursor)
	}
	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil

	case 'd':
		return readFloatBytes(buf, cursor)
	default:
		return nil, cursor, errors.ErrUnexpected("float start with 'd' or 'N'", cursor, buf[cursor])
	}
}

func (d *floatDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	bytes, cursor, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return cursor, nil
	}

	return d.processBytes(bytes, cursor, rv)
}

func (d *floatDecoder) processBytes(bytes []byte, cursor int64, rv reflect.Value) (int64, error) {
	var f64 float64
	var err error
	switch unsafeStr(bytes) {
	case "INF":
		f64 = math.Inf(1)
	case "-INF":
		f64 = math.Inf(-1)
	case "NAN":
		f64 = math.NaN()
	default:
		f64, err = strconv.ParseFloat(unsafeStr(bytes), 64)
	}
	if err != nil {
		return 0, errors.ErrSyntax(err.Error(), cursor)
	}

	if rv.OverflowFloat(f64) {
		return 0, errors.ErrOverflow(f64, rv.Type().Kind().String())
	}

	rv.SetFloat(f64)

	return cursor, nil
}
