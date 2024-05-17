package decoder

import (
	"bytes"
	"encoding"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

func newUnmarshalTextUnmarshalerDecoder(typ *runtime.Type, structName, fieldName string) *unmarshalTextUnmarshalerDecoder {
	return &unmarshalTextUnmarshalerDecoder{
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
	}
}

type unmarshalTextUnmarshalerDecoder struct {
	typ        *runtime.Type
	structName string
	fieldName  string
}

func (d *unmarshalTextUnmarshalerDecoder) annotateError(cursor int64, err error) {
	switch e := err.(type) {
	case *errors.UnmarshalTypeError:
		e.Struct = d.structName
		e.Field = d.fieldName
	case *errors.SyntaxError:
		e.Offset = cursor
	}
}

func (d *unmarshalTextUnmarshalerDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	start := cursor
	end, err := skipValue(buf, cursor, depth)
	if err != nil {
		return 0, err
	}
	src := buf[start:end]
	if len(src) > 0 {
		switch src[0] {
		case 's':
			break
		case 'N':
			if bytes.Equal(src, nullbytes) {
				*(*unsafe.Pointer)(p) = nil
				return end, nil
			}
		default:
			return 0, &errors.UnmarshalTypeError{
				Value:  prefixToTypeName(src[0]),
				Type:   runtime.RType2Type(d.typ),
				Offset: start,
			}
		}
	}

	if s, ok := unquoteBytes(src); ok {
		src = s
	}
	v := *(*any)(unsafe.Pointer(&emptyInterface{
		typ: d.typ,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
	if err := v.(encoding.TextUnmarshaler).UnmarshalText(src); err != nil {
		d.annotateError(cursor, err)
		return 0, err
	}
	return end, nil
}

func prefixToTypeName(b byte) string {
	switch b {
	case 's':
		return "string"
	case 'N':
		return "null"
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return "number"
	case '[':
		return "array"
	case '{':
		return "object"
	}
	return "unknown"
}
