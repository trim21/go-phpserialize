package decoder

import (
	"bytes"
	"reflect"
	"unicode/utf8"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type stringDecoder struct {
	structName string
	fieldName  string
}

func newStringDecoder(structName, fieldName string) *stringDecoder {
	return &stringDecoder{
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *stringDecoder) errUnmarshalType(typeName string, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  typeName,
		Type:   reflect.TypeOf(""),
		Offset: offset,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}

func (d *stringDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	**(**string)(unsafe.Pointer(&p)) = *(*string)(unsafe.Pointer(&bytes))
	return cursor, nil
}

var (
	hexToInt = [256]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}
)

func unicodeToRune(code []byte) rune {
	var r rune
	for i := 0; i < len(code); i++ {
		r = r*16 + rune(hexToInt[code[i]])
	}
	return r
}

func (d *stringDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	switch buf[cursor] {
	case 'n':
		if err := validateNull(buf, cursor); err != nil {
			return nil, 0, err
		}
		cursor += 2
		return nil, cursor, nil
	case 'b':
		return nil, 0, d.errUnmarshalType("bool", cursor)
	case 'd':
		return nil, 0, d.errUnmarshalType("float", cursor)
	case 's':
		cursor++
		break
	case 'i':
		return nil, 0, d.errUnmarshalType("number", cursor)
		// read int as string
	default:
		return nil, 0, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
	}

	s, end, err := readString(buf, cursor)
	if err != nil {
		return nil, 0, err
	}

	return s, end, nil
}

var unescapeMap = [256]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func unsafeAdd(ptr unsafe.Pointer, offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(offset))
}

func unescapeString(buf []byte) int {
	p := (*sliceHeader)(unsafe.Pointer(&buf)).data
	end := unsafeAdd(p, len(buf))
	src := unsafeAdd(p, bytes.IndexByte(buf, '\\'))
	dst := src
	for src != end {
		c := char(src, 0)
		if c == '\\' {
			escapeChar := char(src, 1)
			if escapeChar != 'u' {
				*(*byte)(dst) = unescapeMap[escapeChar]
				src = unsafeAdd(src, 2)
				dst = unsafeAdd(dst, 1)
			} else {
				v1 := hexToInt[char(src, 2)]
				v2 := hexToInt[char(src, 3)]
				v3 := hexToInt[char(src, 4)]
				v4 := hexToInt[char(src, 5)]
				code := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
				if code >= 0xd800 && code < 0xdc00 && uintptr(unsafeAdd(src, 11)) < uintptr(end) {
					if char(src, 6) == '\\' && char(src, 7) == 'u' {
						v1 := hexToInt[char(src, 8)]
						v2 := hexToInt[char(src, 9)]
						v3 := hexToInt[char(src, 10)]
						v4 := hexToInt[char(src, 11)]
						lo := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
						if lo >= 0xdc00 && lo < 0xe000 {
							code = (code-0xd800)<<10 | (lo - 0xdc00) + 0x10000
							src = unsafeAdd(src, 6)
						}
					}
				}
				var b [utf8.UTFMax]byte
				n := utf8.EncodeRune(b[:], code)
				switch n {
				case 4:
					*(*byte)(unsafeAdd(dst, 3)) = b[3]
					fallthrough
				case 3:
					*(*byte)(unsafeAdd(dst, 2)) = b[2]
					fallthrough
				case 2:
					*(*byte)(unsafeAdd(dst, 1)) = b[1]
					fallthrough
				case 1:
					*(*byte)(unsafeAdd(dst, 0)) = b[0]
				}
				src = unsafeAdd(src, 6)
				dst = unsafeAdd(dst, n)
			}
		} else {
			*(*byte)(dst) = c
			src = unsafeAdd(src, 1)
			dst = unsafeAdd(dst, 1)
		}
	}
	return int(uintptr(dst) - uintptr(p))
}
