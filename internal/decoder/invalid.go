package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type invalidDecoder struct {
	typ        reflect.Type
	kind       reflect.Kind
	structName string
	fieldName  string
}

func newInvalidDecoder(typ reflect.Type, structName, fieldName string) *invalidDecoder {
	return &invalidDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *invalidDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	return 0, &errors.UnmarshalTypeError{
		Value:  "object",
		Type:   runtime.RType2Type(d.typ),
		Offset: cursor,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}
