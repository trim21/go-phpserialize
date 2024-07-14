package decoder

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type stringWrappedDecoder interface {
	DecodeString(ctx *RuntimeContext, bytes []byte, topCursor int64, p unsafe.Pointer) error
}

type wrappedStringDecoder struct {
	typ           reflect.Type
	dec           stringWrappedDecoder
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
	isPtrType     bool
}

func newWrappedStringDecoder(typ reflect.Type, dec Decoder, structName, fieldName string) (Decoder, error) {
	var innerDec stringWrappedDecoder
	switch v := dec.(type) {
	case *boolDecoder:
		innerDec = newStringBoolDecoder(structName, fieldName)
	case *floatDecoder:
		innerDec = newStringFloatDecoder(structName, fieldName, v)
	case *uintDecoder:
		innerDec = newStringUintDecoder(structName, fieldName, v)
	case *intDecoder:
		innerDec = newStringIntDecoder(structName, fieldName, v)
	case *ptrDecoder:
		err := v.wrapString()
		return dec, err
	default:
		return nil, &errors.UnsupportedTypeError{Type: runtime.RType2Type(typ)}
	}

	return &wrappedStringDecoder{
		typ:           typ,
		dec:           innerDec,
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
		isPtrType:     typ.Kind() == reflect.Ptr,
	}, nil
}

func (d *wrappedStringDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.stringDecoder.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		if d.isPtrType {
			*(*unsafe.Pointer)(p) = nil
		}
		return c, nil
	}

	if err := d.dec.DecodeString(ctx, bytes, cursor, p); err != nil {
		return 0, err
	}
	return c, nil
}

func newStringBoolDecoder(structName, fieldName string) *stringBoolDecoder {
	return &stringBoolDecoder{}
}

var _ stringWrappedDecoder = (*stringBoolDecoder)(nil)

type stringBoolDecoder struct {
}

func (s stringBoolDecoder) DecodeString(ctx *RuntimeContext, bytes []byte, topCursor int64, p unsafe.Pointer) error {
	str := *(*string)(unsafe.Pointer(&bytes))

	value, err := strconv.ParseBool(str)
	if err != nil {
		return errors.ErrSyntax(err.Error(), topCursor)
	}

	if value {
		**(**bool)(unsafe.Pointer(&p)) = true
	} else {
		**(**bool)(unsafe.Pointer(&p)) = false
	}

	return nil
}

var _ stringWrappedDecoder = (*stringFloatDecoder)(nil)

func newStringFloatDecoder(structName, fieldName string, dec *floatDecoder) *stringFloatDecoder {
	return &stringFloatDecoder{floatDecoder: dec}
}

type stringFloatDecoder struct {
	floatDecoder *floatDecoder
}

func (d stringFloatDecoder) DecodeString(ctx *RuntimeContext, bytes []byte, topCursor int64, p unsafe.Pointer) error {
	_, err := d.floatDecoder.processBytes(bytes, topCursor, p)
	return err
}

var _ stringWrappedDecoder = (*stringUintDecoder)(nil)

func newStringUintDecoder(structName, fieldName string, dec *uintDecoder) *stringUintDecoder {
	return &stringUintDecoder{uintDecoder: dec}
}

type stringUintDecoder struct {
	uintDecoder *uintDecoder
}

func (d stringUintDecoder) DecodeString(ctx *RuntimeContext, bytes []byte, topCursor int64, p unsafe.Pointer) error {
	_, err := d.uintDecoder.processBytes(bytes, topCursor, p)
	return err
}

var _ stringWrappedDecoder = (*stringIntDecoder)(nil)

func newStringIntDecoder(structName string, fieldName string, decoder *intDecoder) *stringIntDecoder {
	return &stringIntDecoder{intDecoder: decoder}
}

type stringIntDecoder struct {
	intDecoder *intDecoder
}

func (d *stringIntDecoder) DecodeString(ctx *RuntimeContext, bytes []byte, topCursor int64, p unsafe.Pointer) error {
	_, err := d.intDecoder.processBytes(bytes, topCursor, p)
	return err
}

func (d *ptrDecoder) wrapString() error {
	var err error
	d.dec, err = newWrappedStringDecoder(d.typ, d.dec, d.structName, d.fieldName)
	return err
}
