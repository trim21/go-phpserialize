package decoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type arrayDecoder struct {
	elemType     *runtime.Type
	size         uintptr
	valueDecoder Decoder
	alen         int
	structName   string
	fieldName    string
	zeroValue    unsafe.Pointer
}

func newArrayDecoder(dec Decoder, elemType *runtime.Type, alen int, structName, fieldName string) *arrayDecoder {
	zeroValue := *(*unsafe.Pointer)(unsafe_New(elemType))
	return &arrayDecoder{
		valueDecoder: dec,
		elemType:     elemType,
		size:         elemType.Size(),
		alen:         alen,
		structName:   structName,
		fieldName:    fieldName,
		zeroValue:    zeroValue,
	}
}

func (d *arrayDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	case 'a':
		cursor++
		if buf[cursor] != ':' {
			return cursor, errors.ErrExpected("':' before array length", cursor)
		}

		// set zero value first, php array may skip some index
		for i := 0; i < d.alen; i++ {
			*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(i)*d.size)) = d.zeroValue
		}

		cursor++
		if buf[cursor] == '0' {
			err := validateEmptyArray(buf, cursor)
			if err != nil {
				return cursor, err
			}
			return cursor + 4, nil
		}

		_, end, err := readLengthInt(buf, cursor-1)
		if err != nil {
			return cursor, err
		}
		cursor = end + 1

		idx := 0
		for {
			currentIndex, end, err := readInt(buf, cursor)
			if err != nil {
				return 0, err
			}

			idx = currentIndex
			cursor = end

			if idx < d.alen {
				c, err := d.valueDecoder.Decode(ctx, cursor, depth, unsafe.Pointer(uintptr(p)+uintptr(idx)*d.size))
				if err != nil {
					return 0, err
				}
				cursor = c
			} else {
				c, err := skipValue(buf, cursor, depth)
				if err != nil {
					return 0, err
				}
				cursor = c
			}

			if buf[cursor] == '}' {
				cursor++
				return cursor, nil
			}
		}
	default:
		return 0, errors.ErrUnexpectedEnd("array", cursor)
	}
}
