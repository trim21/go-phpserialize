package decoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type uintDecoder struct {
	typ        reflect.Type
	kind       reflect.Kind
	op         func(unsafe.Pointer, uint64)
	structName string
	fieldName  string
}

func newUintDecoder(typ reflect.Type, structName, fieldName string, op func(unsafe.Pointer, uint64)) *uintDecoder {
	return &uintDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		op:         op,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *uintDecoder) typeError(buf []byte, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   runtime.RType2Type(d.typ),
		Offset: offset,
	}
}

var (
	pow10u64 = [...]uint64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	}
	pow10u64Len = len(pow10u64)
)

func (d *uintDecoder) parseUint(b []byte) (uint64, error) {
	maxDigit := len(b)
	if maxDigit > pow10u64Len {
		return 0, fmt.Errorf("invalid length of number")
	}
	sum := uint64(0)
	for i := 0; i < maxDigit; i++ {
		c := uint64(b[i]) - 48
		digitValue := pow10u64[maxDigit-i-1]
		sum += c * digitValue
	}
	return sum, nil
}

func (d *uintDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	if char(b, cursor) != 'i' {
		return nil, cursor, errors.ErrExpected("int", cursor)
	}

	cursor++
	if char(b, cursor) != ':' {
		return nil, cursor, errors.ErrExpected("int sep ':'", cursor)
	}
	cursor++

	switch char(b, cursor) {
	case '0':
		cursor++
		if char(b, cursor) != ';' {
			return nil, cursor, errors.ErrExpected("';' end int", cursor)
		}
		return numZeroBuf, cursor + 1, nil
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start := cursor
		cursor++
		for numTable[char(b, cursor)] {
			cursor++
		}
		if char(b, cursor) != ';' {
			return nil, cursor, errors.ErrExpected("';' end int", cursor)
		}
		num := buf[start:cursor]
		return num, cursor + 1, nil
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil
	default:
		return nil, 0, d.typeError([]byte{char(b, cursor)}, cursor)
	}
}

func (d *uintDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c

	return d.processBytes(bytes, cursor, p)
}

func (d *uintDecoder) processBytes(bytes []byte, cursor int64, p unsafe.Pointer) (int64, error) {
	u64, err := d.parseUint(bytes)
	if err != nil {
		return 0, d.typeError(bytes, cursor)
	}
	switch d.kind {
	case reflect.Uint8:
		if (1 << 8) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Uint16:
		if (1 << 16) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Uint32:
		if (1 << 32) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	}
	d.op(p, u64)
	return cursor, nil
}
