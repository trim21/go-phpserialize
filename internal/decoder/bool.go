package decoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type boolDecoder struct {
	structName string
	fieldName  string
}

func newBoolDecoder(structName, fieldName string) *boolDecoder {
	return &boolDecoder{structName: structName, fieldName: fieldName}
}

func (d *boolDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	c := s.skipWhiteSpace()
	for {
		switch c {
		case 't':
			if err := trueBytes(s); err != nil {
				return err
			}
			**(**bool)(unsafe.Pointer(&p)) = true
			return nil
		case 'f':
			if err := falseBytes(s); err != nil {
				return err
			}
			**(**bool)(unsafe.Pointer(&p)) = false
			return nil
		case 'n':
			if err := nullBytes(s); err != nil {
				return err
			}
			return nil
		case nul:
			if s.read() {
				c = s.char()
				continue
			}
			goto ERROR
		}
		break
	}
ERROR:
	return errors.ErrUnexpectedEnd("bool", s.totalOffset())
}

func (d *boolDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	switch buf[cursor] {
	case 'b':
		// b:0;
		// b:1;

		cursor++
		if buf[cursor] != ':' {
			return 0, errors.ErrUnexpectedEnd("':' before bool value", cursor)
		}
		cursor++
		switch buf[cursor] {
		case '0':
		case '1':
			**(**bool)(unsafe.Pointer(&p)) = true
		default:
			return 0, errors.ErrUnexpectedEnd("'0' pr '1' af bool value", cursor)
		}
		cursor++
		if buf[cursor] != ';' {
			return 0, errors.ErrUnexpectedEnd("';' end bool value", cursor)
		}
		cursor++
		return cursor, nil

	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	}
	return 0, errors.ErrUnexpectedEnd("bool", cursor)
}
