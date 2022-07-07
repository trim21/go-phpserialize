package decoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type ptrDecoder struct {
	dec        Decoder
	typ        *runtime.Type
	structName string
	fieldName  string
}

func newPtrDecoder(dec Decoder, typ *runtime.Type, structName, fieldName string) *ptrDecoder {
	return &ptrDecoder{
		dec:        dec,
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *ptrDecoder) contentDecoder() Decoder {
	dec, ok := d.dec.(*ptrDecoder)
	if !ok {
		return d.dec
	}
	return dec.contentDecoder()
}

//nolint:golint
//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(*runtime.Type) unsafe.Pointer

func (d *ptrDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	if buf[cursor] == 'N' {
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		if p != nil {
			*(*unsafe.Pointer)(p) = nil
		}
		cursor += 2
		return cursor, nil
	}

	var newptr unsafe.Pointer

	if *(*unsafe.Pointer)(p) == nil {
		newptr = unsafe_New(d.typ)
		*(*unsafe.Pointer)(p) = newptr
	} else {
		newptr = *(*unsafe.Pointer)(p)
	}

	c, err := d.dec.Decode(ctx, cursor, depth, newptr)
	if err != nil {
		return 0, err
	}
	cursor = c
	return cursor, nil
}
