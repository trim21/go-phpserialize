package decoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type ptrDecoder struct {
	dec        Decoder
	typ        reflect.Type
	structName string
	fieldName  string
}

func newPtrDecoder(dec Decoder, typ reflect.Type, structName, fieldName string) (Decoder, error) {
	if typ.Kind() == reflect.Ptr {
		return nil, &errors.UnsupportedTypeError{
			Type: reflect.PointerTo(typ),
		}
	}
	return &ptrDecoder{
		dec:        dec,
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
	}, nil
}

func (d *ptrDecoder) contentDecoder() Decoder {
	dec, ok := d.dec.(*ptrDecoder)
	if !ok {
		return d.dec
	}
	return dec.contentDecoder()
}

func (d *ptrDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	if buf[cursor] == 'N' {
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		rv.SetZero()
		cursor += 2
		return cursor, nil
	}

	if rv.IsNil() {
		np := reflect.New(d.typ)
		rv.Set(np)
	}

	c, err := d.dec.Decode(ctx, cursor, depth, rv.Elem())
	if err != nil {
		return 0, err
	}
	cursor = c

	return cursor, nil
}
