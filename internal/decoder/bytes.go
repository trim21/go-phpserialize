package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type bytesDecoder struct {
	typ           reflect.Type
	sliceDecoder  Decoder
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
}

func byteUnmarshalerSliceDecoder(typ reflect.Type, structName string, fieldName string) Decoder {
	var unmarshalDecoder Decoder
	switch {
	case reflect.PointerTo(typ).Implements(unmarshalPHPType):
		unmarshalDecoder = newUnmarshalTextDecoder(reflect.PointerTo(typ), structName, fieldName)
	default:
		unmarshalDecoder, _ = compileUint8(typ, structName, fieldName)
	}
	return newSliceDecoder(unmarshalDecoder, typ, 1, structName, fieldName)
}

func newBytesDecoder(typ reflect.Type, structName string, fieldName string) *bytesDecoder {
	return &bytesDecoder{
		typ:           typ,
		sliceDecoder:  byteUnmarshalerSliceDecoder(typ, structName, fieldName),
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
	}
}

func (d *bytesDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeBinary(ctx, cursor, depth, p)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	*(*[]byte)(p) = bytes
	return cursor, nil
}

func (d *bytesDecoder) decodeBinary(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) ([]byte, int64, error) {
	buf := ctx.Buf
	if buf[cursor] == 'a' {
		if d.sliceDecoder == nil {
			return nil, 0, &errors.UnmarshalTypeError{
				Type:   runtime.RType2Type(d.typ),
				Offset: cursor,
			}
		}
		c, err := d.sliceDecoder.Decode(ctx, cursor, depth, p)
		if err != nil {
			return nil, 0, err
		}
		return nil, c, nil
	}
	return d.stringDecoder.decodeByte(buf, cursor)
}
