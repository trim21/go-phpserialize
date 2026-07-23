package decoder

import (
	"strconv"
	"sync"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type RuntimeContext struct {
	Buf          []byte
	MaxMapSize   uint
	MaxSliceSize uint
}

var (
	runtimeContextPool = sync.Pool{
		New: func() any {
			return &RuntimeContext{}
		},
	}
)

func TakeRuntimeContext() *RuntimeContext {
	return runtimeContextPool.Get().(*RuntimeContext)
}

func ReleaseRuntimeContext(ctx *RuntimeContext) {
	ctx.Buf = nil
	ctx.MaxSliceSize = 0
	ctx.MaxMapSize = 0
	runtimeContextPool.Put(ctx)
}

var (
	isWhiteSpace = [256]bool{}

	isInteger = [256]bool{}
)

func init() {
	isWhiteSpace[' '] = true
	isWhiteSpace['\n'] = true
	isWhiteSpace['\t'] = true
	isWhiteSpace['\r'] = true

	for i := 0; i < 10; i++ {
		isInteger[[]byte(strconv.Itoa(i))[0]] = true
	}
}

func hasByte(buf []byte, cursor int64) bool {
	return cursor >= 0 && cursor < int64(len(buf))
}

func hasBytes(buf []byte, cursor, size int64) bool {
	return cursor >= 0 && size >= 0 && cursor <= int64(len(buf))-size
}

// `:${length}:`
func skipLengthWithBothColon(buf []byte, cursor int64) (int64, error) {
	if !hasByte(buf, cursor) {
		return cursor, errors.ErrUnexpectedEnd("length", cursor)
	}
	if buf[cursor] != ':' {
		return cursor, errors.ErrUnexpected("':' before length", cursor, buf[cursor])
	}
	cursor++

	start := cursor
	for hasByte(buf, cursor) && isInteger[buf[cursor]] {
		cursor++
	}
	if cursor == start {
		return cursor, errors.ErrSyntax("php-serialize: length requires at least one digit", cursor)
	}

	if !hasByte(buf, cursor) {
		return cursor, errors.ErrUnexpectedEnd("length", cursor)
	}
	if buf[cursor] != ':' {
		return cursor, errors.ErrUnexpected("':' after length", cursor, buf[cursor])
	}

	cursor++

	return cursor, nil
}

// jump from 'O' to colon ':' before array length
func skipClassName(buf []byte, cursor int64) (int64, error) {
	// O:8:"stdClass":1:{s:1:"a";s:1:"q";}
	if !hasByte(buf, cursor) || buf[cursor] != 'O' {
		return cursor, errors.ErrUnexpectedEnd("object", cursor)
	}
	classLen, start, err := readLength(buf, cursor+1)
	if err != nil {
		return cursor, err
	}
	if !hasByte(buf, start) || buf[start] != '"' {
		return start, errors.ErrUnexpectedEnd("class name", start)
	}
	end := start + classLen + 1
	if !hasBytes(buf, end, 2) {
		return end, errors.ErrUnexpectedEnd("class name", end)
	}
	if buf[end] != '"' {
		return end, errors.ErrUnexpected(`class name quoted '"'`, end, buf[end])
	}
	if buf[end+1] != ':' {
		return end + 1, errors.ErrUnexpected("':' after class name", end+1, buf[end+1])
	}
	return end, nil
}

// i:{};
// the cursor should point to the beginning `i`
// and it will point to the byte after `;`
func skipInt(buf []byte, cursor int64) (int64, error) {
	_, end, err := readIntegerBytes(buf, cursor)
	return end, err
}

func skipString(buf []byte, cursor int64) (int64, error) {
	if !hasByte(buf, cursor) || buf[cursor] != 's' {
		return cursor, errors.ErrUnexpectedEnd("string", cursor)
	}
	_, end, err := readString(buf, cursor+1)
	return end, err
}

func skipArray(buf []byte, cursor, depth int64) (int64, error) {
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth('a', cursor)
	}
	length, end, err := readLength(buf, cursor)
	if err != nil {
		return cursor, err
	}
	cursor = end
	if !hasByte(buf, cursor) {
		return cursor, errors.ErrUnexpectedEnd("array", cursor)
	}
	if buf[cursor] != '{' {
		return cursor, errors.ErrInvalidBeginningOfArray(buf[cursor], cursor)
	}
	cursor++
	for range length {
		cursor, err = skipValue(buf, cursor, depth)
		if err != nil {
			return cursor, err
		}
		cursor, err = skipValue(buf, cursor, depth)
		if err != nil {
			return cursor, err
		}
	}
	if !hasByte(buf, cursor) {
		return cursor, errors.ErrUnexpectedEnd("array", cursor)
	}
	if buf[cursor] != '}' {
		return cursor, errors.ErrUnexpected("'}' after array", cursor, buf[cursor])
	}
	return cursor + 1, nil
}

func skipValue(buf []byte, cursor, depth int64) (int64, error) {
	if !hasByte(buf, cursor) {
		return cursor, errors.ErrUnexpectedEnd("value", cursor)
	}
	switch buf[cursor] {
	case 'O':
		end, err := skipClassName(buf, cursor)
		if err != nil {
			return cursor, err
		}
		cursor = end
		fallthrough
	case 'a':
		return skipArray(buf, cursor+1, depth+1)
	case 's':
		return skipString(buf, cursor)
	case 'd':
		_, end, err := readFloatBytes(buf, cursor)
		return end, err
	case 'i':
		return skipInt(buf, cursor)
	case 'b':
		_, err := readBool(buf, cursor)
		return cursor + 4, err

	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	default:
		return cursor, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
	}
}

// caller should check `a:0` , this function check  `0:{};`
func validateEmptyArray(buf []byte, cursor int64) error {
	if cursor+4 > int64(len(buf)) {
		return errors.ErrUnexpectedEnd("array", cursor)
	}

	if buf[cursor+1] != ':' {
		return errors.ErrUnexpected("':' before array length", cursor+1, buf[cursor+1])
	}
	if buf[cursor+2] != '{' {
		return errors.ErrInvalidBeginningOfArray(buf[cursor+2], cursor+2)
	}
	if buf[cursor+3] != '}' {
		return errors.ErrUnexpected("empty array end with '}'", cursor+3, buf[cursor+3])

	}

	return nil
}

func validateNull(buf []byte, cursor int64) error {
	if cursor+1 >= int64(len(buf)) {
		return errors.ErrUnexpectedEnd("null", cursor)
	}
	if buf[cursor+1] != ';' {
		return errors.ErrInvalidCharacter(buf[cursor+1], "null", cursor)
	}
	return nil
}

// :${length}:
func readLength(buf []byte, cursor int64) (int64, int64, error) {
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return 0, cursor, err
	}
	length, err := parseInt64(buf[cursor+1 : end-1])
	if err != nil || length < 0 {
		return 0, cursor, errors.ErrSyntax("php-serialize: invalid length", cursor+1)
	}
	return length, end, nil
}

// :${length}:
func readLengthInt(buf []byte, cursor int64) (int, int64, error) {
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return 0, cursor, err
	}

	length, err := parseInt64(buf[cursor+1 : end-1])
	if err != nil || length < 0 || int64(int(length)) != length {
		return 0, cursor, errors.ErrSyntax("php-serialize: invalid length", cursor+1)
	}
	return int(length), end, nil
}

// `:${length}:"${content}";`
func readString(buf []byte, cursor int64) ([]byte, int64, error) {
	sLen, end, err := readLength(buf, cursor)
	if err != nil {
		return nil, 0, err
	}
	if !hasByte(buf, end) {
		return nil, end, errors.ErrUnexpectedEnd("string", end)
	}
	if buf[end] != '"' {
		return nil, end, errors.ErrUnexpected(`string quoted '"'`, end, buf[end])
	}
	start := end + 1
	if !hasBytes(buf, start, sLen+2) {
		return nil, start, errors.ErrUnexpectedEnd("string", start)
	}
	end = start + sLen
	if buf[end] != '"' {
		return nil, end, errors.ErrUnexpected(`string quoted '"'`, end, buf[end])
	}
	cursor = end + 1
	if buf[cursor] != ';' {
		return nil, end, errors.ErrUnexpected(`string end ';'`, cursor, buf[cursor])
	}

	cursor++

	return buf[start:end], cursor, nil
}

func readIntegerBytes(buf []byte, cursor int64) ([]byte, int64, error) {
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("integer", cursor)
	}
	if buf[cursor] != 'i' {
		return nil, cursor, errors.ErrUnexpected("'i' to start an integer", cursor, buf[cursor])
	}
	if !hasBytes(buf, cursor, 3) {
		return nil, cursor, errors.ErrUnexpectedEnd("integer", cursor)
	}
	if buf[cursor+1] != ':' {
		return nil, cursor + 1, errors.ErrUnexpected("int separator ':'", cursor+1, buf[cursor+1])
	}
	cursor += 2
	start := cursor
	if buf[cursor] == '-' {
		cursor++
	}
	digitStart := cursor
	for hasByte(buf, cursor) && numTable[buf[cursor]] {
		cursor++
	}
	if cursor == digitStart {
		return nil, cursor, errors.ErrSyntax("php-serialize: integer requires at least one digit", cursor)
	}
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("integer", cursor)
	}
	if buf[cursor] != ';' {
		return nil, cursor, errors.ErrUnexpected("';' after integer", cursor, buf[cursor])
	}
	return buf[start:cursor], cursor + 1, nil
}

func readFloatBytes(buf []byte, cursor int64) ([]byte, int64, error) {
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("float", cursor)
	}
	if buf[cursor] != 'd' {
		return nil, cursor, errors.ErrUnexpected("'d' to start a float", cursor, buf[cursor])
	}
	if !hasBytes(buf, cursor, 3) {
		return nil, cursor, errors.ErrUnexpectedEnd("float", cursor)
	}
	if buf[cursor+1] != ':' {
		return nil, cursor + 1, errors.ErrUnexpected("float separator ':'", cursor+1, buf[cursor+1])
	}
	start := cursor + 2
	cursor = start
	for hasByte(buf, cursor) && buf[cursor] != ';' {
		cursor++
	}
	if cursor == start {
		return nil, cursor, errors.ErrSyntax("php-serialize: float requires a value", cursor)
	}
	if !hasByte(buf, cursor) {
		return nil, cursor, errors.ErrUnexpectedEnd("float", cursor)
	}
	return buf[start:cursor], cursor + 1, nil
}

func readBool(buf []byte, cursor int64) (bool, error) {
	if !hasBytes(buf, cursor, 4) {
		return false, errors.ErrUnexpectedEnd("bool", cursor)
	}
	if buf[cursor] != 'b' {
		return false, errors.ErrInvalidCharacter(buf[cursor], "bool", cursor)
	}

	if buf[cursor+1] != ':' {
		return false, errors.ErrInvalidCharacter(buf[cursor+1], "bool", cursor)
	}

	var v bool
	switch buf[cursor+2] {
	case '1':
		v = true
	case '0':
		v = false
	default:
		return false, errors.ErrInvalidCharacter(buf[cursor+2], "bool", cursor)
	}

	if buf[cursor+3] != ';' {
		return false, errors.ErrInvalidCharacter(buf[cursor+3], "bool", cursor+3)
	}

	return v, nil
}
