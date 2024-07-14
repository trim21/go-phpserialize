package decoder

import (
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"sort"
	"strings"
	"unicode"
	"unicode/utf16"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type structFieldSet struct {
	dec         Decoder
	offset      uintptr
	isTaggedKey bool
	fieldIdx    int
	key         string
	keyLen      int64
	err         error
}

type structDecoder struct {
	fieldMap           map[string]*structFieldSet
	fieldUniqueNameNum int
	stringDecoder      *stringDecoder
	structName         string
	fieldName          string
	isTriedOptimize    bool
	keyBitmapUint8     [][256]uint8
	keyBitmapUint16    [][256]uint16
	sortedFieldSets    []*structFieldSet
	keyDecoder         func(*structDecoder, []byte, int64) (int64, *structFieldSet, error)
	// keyStreamDecoder   func(*structDecoder, *Stream) (*structFieldSet, string, error)
}

var (
	largeToSmallTable [256]byte
)

func init() {
	for i := 0; i < 256; i++ {
		c := i
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		largeToSmallTable[i] = byte(c)
	}
}

func newStructDecoder(structName, fieldName string, fieldMap map[string]*structFieldSet) *structDecoder {
	return &structDecoder{
		fieldMap:      fieldMap,
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
		keyDecoder:    decodeKey,
		// keyFStreamDecoder: decodeKeyStream,
	}
}

const (
	allowOptimizeMaxKeyLen   = 64
	allowOptimizeMaxFieldLen = 16
)

func (d *structDecoder) tryOptimize() {
	fieldUniqueNameMap := map[string]int{}
	fieldIdx := -1
	for k, v := range d.fieldMap {
		lower := strings.ToLower(k)
		idx, exists := fieldUniqueNameMap[lower]
		if exists {
			v.fieldIdx = idx
		} else {
			fieldIdx++
			v.fieldIdx = fieldIdx
		}
		fieldUniqueNameMap[lower] = fieldIdx
	}
	d.fieldUniqueNameNum = len(fieldUniqueNameMap)

	if d.isTriedOptimize {
		return
	}
	fieldMap := map[string]*structFieldSet{}
	conflicted := map[string]struct{}{}
	for k, v := range d.fieldMap {
		key := strings.ToLower(k)
		if key != k {
			// already exists same key (e.g. Hello and HELLO has same lower case key
			if _, exists := conflicted[key]; exists {
				d.isTriedOptimize = true
				return
			}
			conflicted[key] = struct{}{}
		}
		if field, exists := fieldMap[key]; exists {
			if field != v {
				d.isTriedOptimize = true
				return
			}
		}
		fieldMap[key] = v
	}

	if len(fieldMap) > allowOptimizeMaxFieldLen {
		d.isTriedOptimize = true
		return
	}

	var maxKeyLen int
	sortedKeys := []string{}
	for key := range fieldMap {
		keyLen := len(key)
		if keyLen > allowOptimizeMaxKeyLen {
			d.isTriedOptimize = true
			return
		}
		if maxKeyLen < keyLen {
			maxKeyLen = keyLen
		}
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// By allocating one extra capacity than `maxKeyLen`,
	// it is possible to avoid the process of comparing the index of the key with the length of the bitmap each time.
	bitmapLen := maxKeyLen + 1
	// TODO:  this
	// if len(sortedKeys) <= 8 {
	// 	keyBitmap := make([][256]uint8, bitmapLen)
	// 	for i, key := range sortedKeys {
	// 		for j := 0; j < len(key); j++ {
	// 			c := key[j]
	// 			keyBitmap[j][c] |= (1 << uint(i))
	// 		}
	// 		d.sortedFieldSets = append(d.sortedFieldSets, fieldMap[key])
	// 	}
	// 	d.keyBitmapUint8 = keyBitmap
	// d.keyDecoder = decodeKeyByBitmapUint8
	// d.keyStreamDecoder = decodeKeyByBitmapUint8Stream
	// } else {
	keyBitmap := make([][256]uint16, bitmapLen)
	for i, key := range sortedKeys {
		for j := 0; j < len(key); j++ {
			c := key[j]
			keyBitmap[j][c] |= 1 << uint(i)
		}
		d.sortedFieldSets = append(d.sortedFieldSets, fieldMap[key])
	}
	d.keyBitmapUint16 = keyBitmap
	d.keyDecoder = decodeKeyByBitmapUint16
	// }
}

// decode from '\uXXXX'
func decodeKeyCharByUnicodeRune(buf []byte, cursor int64) ([]byte, int64) {
	const defaultOffset = 4
	const surrogateOffset = 6

	r := unicodeToRune(buf[cursor : cursor+defaultOffset])
	if utf16.IsSurrogate(r) {
		cursor += defaultOffset
		if cursor+surrogateOffset >= int64(len(buf)) || buf[cursor] != '\\' || buf[cursor+1] != 'u' {
			return []byte(string(unicode.ReplacementChar)), cursor + defaultOffset - 1
		}
		cursor += 2
		r2 := unicodeToRune(buf[cursor : cursor+defaultOffset])
		if r := utf16.DecodeRune(r, r2); r != unicode.ReplacementChar {
			return []byte(string(r)), cursor + defaultOffset - 1
		}
	}
	return []byte(string(r)), cursor + defaultOffset - 1
}

func decodeKeyCharByEscapedChar(buf []byte, cursor int64) ([]byte, int64) {
	c := buf[cursor]
	cursor++
	switch c {
	case '"':
		return []byte{'"'}, cursor
	case '\\':
		return []byte{'\\'}, cursor
	case '/':
		return []byte{'/'}, cursor
	case 'b':
		return []byte{'\b'}, cursor
	case 'f':
		return []byte{'\f'}, cursor
	case 'n':
		return []byte{'\n'}, cursor
	case 'r':
		return []byte{'\r'}, cursor
	case 't':
		return []byte{'\t'}, cursor
	case 'u':
		return decodeKeyCharByUnicodeRune(buf, cursor)
	}
	return nil, cursor
}

// TODO: not finished
func decodeKeyByBitmapUint8(d *structDecoder, buf []byte, cursor int64) (int64, *structFieldSet, error) {
	var (
		curBit uint8 = math.MaxUint8
	)
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	for {
		switch char(b, cursor) {
		case 'i':
			// array with int key, should we skip or just omit?

		// case '"':
		case 's':
			cursor++
			c := char(b, cursor)
			if c != ':' {
				return 0, nil, errors.ErrSyntax(fmt.Sprintf("unexpected chat (%c) before str length", c), cursor)
			}

			cursor++
			sLen, end, err := readLength(buf, cursor)
			if err != nil {
				return 0, nil, err
			}
			cursor = end

			c = char(b, cursor)
			if c != ':' {
				return 0, nil, errors.ErrSyntax(fmt.Sprintf("unexpected chat (%c) before str length", c), cursor)
			}

			runtime.KeepAlive(sLen)
			cursor++
			c = char(b, cursor)
			switch c {
			case '"':
				cursor++
				return cursor, nil, nil
			case nul:
				return 0, nil, errors.ErrUnexpectedEnd("string", cursor)
			}
			keyIdx := 0
			bitmap := d.keyBitmapUint8
			start := cursor
			for {
				c := char(b, cursor)
				switch c {
				case '"':
					fieldSetIndex := bits.TrailingZeros8(curBit)
					field := d.sortedFieldSets[fieldSetIndex]
					keyLen := cursor - start
					cursor++
					if keyLen < field.keyLen {
						// early match
						return cursor, nil, nil
					}
					return cursor, field, nil
				case nul:
					return 0, nil, errors.ErrUnexpectedEnd("string", cursor)
				case '\\':
					cursor++
					chars, nextCursor := decodeKeyCharByEscapedChar(buf, cursor)
					for _, c := range chars {
						curBit &= bitmap[keyIdx][largeToSmallTable[c]]
						if curBit == 0 {
							return decodeKeyNotFound(b, cursor)
						}
						keyIdx++
					}
					cursor = nextCursor
				default:
					curBit &= bitmap[keyIdx][largeToSmallTable[c]]
					if curBit == 0 {
						return decodeKeyNotFound(b, cursor)
					}
					keyIdx++
				}
				cursor++
			}
		default:
			return cursor, nil, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
		}
	}
}

func decodeKeyByBitmapUint16(d *structDecoder, buf []byte, cursor int64) (int64, *structFieldSet, error) {
	var (
		curBit uint16 = math.MaxUint16
	)
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data

	switch char(b, cursor) {
	case 'i':
	// TODO array with int key
	// array with int key, should we skip or just omit?
	case 's':
		cursor++
		sLen, end, err := readLength(buf, cursor)
		if err != nil {
			return 0, nil, err
		}
		cursor = end
		cursor++ // '"'

		keyIdx := 0
		bitmap := d.keyBitmapUint16
		start := cursor

		if char(b, start+sLen) != '"' {
			return 0, nil, errors.ErrExpected("string should be quoted", cursor)
		}

		if char(b, start+sLen+1) != ';' {
			return 0, nil, errors.ErrExpected("string end with semi", cursor)
		}

		for i := start; i < start+sLen; i++ {
			cursor = i
			c := char(b, cursor)
			curBit &= bitmap[keyIdx][largeToSmallTable[c]]
			if curBit == 0 {
				return decodeKeyNotFound(b, cursor)
			}
			keyIdx++
		}

		fieldSetIndex := bits.TrailingZeros16(curBit)
		field := d.sortedFieldSets[fieldSetIndex]
		cursor++
		if sLen < field.keyLen {
			// early match
			return cursor, nil, nil
		}
		cursor++ // '"'
		cursor++ // ';'
		return cursor, field, nil
	}

	return cursor, nil, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
}

func decodeKeyNotFound(b unsafe.Pointer, cursor int64) (int64, *structFieldSet, error) {
	for {
		cursor++
		switch char(b, cursor) {
		case '"':
			cursor += 2
			return cursor, nil, nil
		case nul:
			return 0, nil, errors.ErrUnexpectedEnd("string", cursor)
		}
	}
}

func decodeKey(d *structDecoder, buf []byte, cursor int64) (int64, *structFieldSet, error) {
	key, c, err := d.stringDecoder.decodeByte(buf, cursor)
	if err != nil {
		return 0, nil, err
	}
	cursor = c
	k := *(*string)(unsafe.Pointer(&key))
	field, exists := d.fieldMap[k]
	if !exists {
		return cursor, nil, nil
	}
	return cursor, field, nil
}

func (d *structDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}
	buflen := int64(len(buf))
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	switch char(b, cursor) {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	case 'O':
		// O:8:"stdClass":1:{s:1:"a";s:1:"q";}
		end, err := skipClassName(buf, cursor)
		if err != nil {
			return cursor, err
		}
		cursor = end
		fallthrough
	case 'a':
		cursor++
		if buf[cursor] != ':' {
			return 0, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
		}
	default:
		return 0, errors.ErrInvalidBeginningOfValue(char(b, cursor), cursor)
	}

	// skip  :${length}:
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return cursor, err
	}
	cursor = end
	if buf[cursor] != '{' {
		return 0, errors.ErrInvalidBeginningOfArray(char(b, cursor), cursor)
	}

	cursor++
	if buf[cursor] == '}' {
		cursor++
		return cursor, nil
	}

	for {
		c, field, err := d.keyDecoder(d, buf, cursor)
		if err != nil {
			return 0, err
		}

		cursor = c

		// cursor++
		if cursor >= buflen {
			return 0, errors.ErrExpected("object value after colon", cursor)
		}
		if field != nil {
			if field.err != nil {
				return 0, field.err
			}
			c, err := field.dec.Decode(ctx, cursor, depth, unsafe.Pointer(uintptr(p)+field.offset))
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

		if char(b, cursor) == '}' {
			cursor++
			return cursor, nil
		}
	}
}
