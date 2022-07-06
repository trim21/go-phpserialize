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

	for {
		switch buf[cursor] {
		case 'N':
			if err := validateNull(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 2
			return cursor, nil
		case 'a':
			// TODO
			idx := 0
			cursor++
			if buf[cursor] == ']' {
				for idx < d.alen {
					*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
					idx++
				}
				cursor++
				return cursor, nil
			}
			for {
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
				idx++
				switch buf[cursor] {
				case ']':
					for idx < d.alen {
						*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
						idx++
					}
					cursor++
					return cursor, nil
				case ',':
					cursor++
					continue
				default:
					return 0, errors.ErrInvalidCharacter(buf[cursor], "array", cursor)
				}
			}
		default:
			return 0, errors.ErrUnexpectedEnd("array", cursor)
		}
	}
}
