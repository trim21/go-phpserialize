package decoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type wrappedStringDecoder struct {
	typ           *runtime.Type
	dec           Decoder
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
	isPtrType     bool
}

func newWrappedStringDecoder(typ *runtime.Type, dec Decoder, structName, fieldName string) *wrappedStringDecoder {
	return &wrappedStringDecoder{
		typ:           typ,
		dec:           dec,
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
		isPtrType:     typ.Kind() == reflect.Ptr,
	}
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
	bytes = append(bytes, nul)
	oldBuf := ctx.Buf
	ctx.Buf = bytes
	fmt.Println(string(ctx.Buf))
	if _, err := d.dec.Decode(ctx, 0, depth, p); err != nil {
		return 0, err
	}
	ctx.Buf = oldBuf
	return c, nil
}
