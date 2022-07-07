package decoder

import (
	"strconv"
	"sync"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

type RuntimeContext struct {
	Buf []byte
}

var (
	runtimeContextPool = sync.Pool{
		New: func() interface{} {
			return &RuntimeContext{}
		},
	}
)

func TakeRuntimeContext() *RuntimeContext {
	return runtimeContextPool.Get().(*RuntimeContext)
}

func ReleaseRuntimeContext(ctx *RuntimeContext) {
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

func char(ptr unsafe.Pointer, offset int64) byte {
	return *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(offset)))
}

// `:${length}:`
func skipLengthWithBothColon(buf []byte, cursor int64) (int64, error) {
	if buf[cursor] != ':' {
		return cursor, errors.ErrExpected("':' before length", cursor)
	}
	cursor++

	for isInteger[buf[cursor]] {
		cursor++
	}

	if buf[cursor] != ':' {
		return cursor, errors.ErrExpected("':' after length", cursor)
	}

	cursor++

	return cursor, nil
}

func skipString(buf []byte, cursor int64) (int64, error) {
	cursor++
	sLen, end, err := readLength(buf, cursor)
	if err != nil {
		return cursor, err
	}

	return end + sLen + 3, nil
}

func skipArray(buf []byte, cursor, depth int64) (int64, error) {
	bracketCount := 1
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return cursor, err
	}
	cursor = end

	for {
		switch buf[cursor] {
		// '{' and '}' may only appear in array or string,
		// we will skip value content, it's easy to only scan char '{' and '}'
		case 's':
			end, err := skipString(buf, cursor)
			if err != nil {
				return cursor, err
			}
			cursor = end
		case '}':
			bracketCount--
			depth--
			if bracketCount == 0 {
				return cursor + 1, nil
			}
		case '{':
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
			}
			cursor++
		default:
			cursor++
		}
	}
}

func skipValue(buf []byte, cursor, depth int64) (int64, error) {

	switch buf[cursor] {
	case 'a':
		return skipArray(buf, cursor+1, depth+1)
	case 's':
		return skipString(buf, cursor)
	// case 'd':

	case 'i':
		cursor++
		end, err := skipLengthWithBothColon(buf, cursor)
		if err != nil {
			return cursor, err
		}
		return end + 1, nil
	case 'b':
		cursor++
		if buf[cursor] != ':' {
			return 0, errors.ErrUnexpectedEnd("':' before bool value", cursor)
		}
		cursor++
		switch buf[cursor] {
		case '0':
		case '1':
		default:
			return 0, errors.ErrUnexpectedEnd("'0' pr '1' af bool value", cursor)
		}
		cursor++
		if buf[cursor] != ';' {
			return 0, errors.ErrUnexpectedEnd("';' end bool value", cursor)
		}
		cursor++
		return cursor, nil

	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		return cursor, nil
	default:
		return cursor, errors.ErrUnexpectedEnd("null", cursor)
	}

}

func validateTrue(buf []byte, cursor int64) error {
	if cursor+3 >= int64(len(buf)) {
		return errors.ErrUnexpectedEnd("true", cursor)
	}
	if buf[cursor+1] != 'r' {
		return errors.ErrInvalidCharacter(buf[cursor+1], "true", cursor)
	}
	if buf[cursor+2] != 'u' {
		return errors.ErrInvalidCharacter(buf[cursor+2], "true", cursor)
	}
	if buf[cursor+3] != 'e' {
		return errors.ErrInvalidCharacter(buf[cursor+3], "true", cursor)
	}
	return nil
}

func validateFalse(buf []byte, cursor int64) error {
	if cursor+4 >= int64(len(buf)) {
		return errors.ErrUnexpectedEnd("false", cursor)
	}
	if buf[cursor+1] != 'a' {
		return errors.ErrInvalidCharacter(buf[cursor+1], "false", cursor)
	}
	if buf[cursor+2] != 'l' {
		return errors.ErrInvalidCharacter(buf[cursor+2], "false", cursor)
	}
	if buf[cursor+3] != 's' {
		return errors.ErrInvalidCharacter(buf[cursor+3], "false", cursor)
	}
	if buf[cursor+4] != 'e' {
		return errors.ErrInvalidCharacter(buf[cursor+4], "false", cursor)
	}
	return nil
}

// caller should check `a:0` , this function check  `0:{};`
func validateEmptyArray(buf []byte, cursor int64) error {
	if cursor+4 >= int64(len(buf)) {
		return errors.ErrUnexpectedEnd("null", cursor)
	}

	if buf[cursor+1] != ':' {
		return errors.ErrExpected("':' before array length", cursor+1)
	}
	if buf[cursor+2] != '{' {
		return errors.ErrInvalidBeginningOfArray(buf[cursor+2], cursor+2)
	}
	if buf[cursor+3] != '}' {
		return errors.ErrExpected("empty array end with '}'", cursor+3)

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

	return parseByteStringInt64(buf[cursor+1 : end-1]), end, nil
}

// :${length}:
func readLengthInt(buf []byte, cursor int64) (int, int64, error) {
	end, err := skipLengthWithBothColon(buf, cursor)
	if err != nil {
		return 0, cursor, err
	}

	return parseByteStringInt(buf[cursor+1 : end-1]), end, nil
}

// `:${length}:"${content}";`
func readString(buf []byte, cursor int64) ([]byte, int64, error) {
	sLen, end, err := readLength(buf, cursor)
	if err != nil {
		return nil, 0, err
	}

	start := end + 1
	end = end + sLen + 1
	cursor = end + 2

	if buf[end] != '"' {
		return nil, end, errors.ErrExpected(`string quoted '"'`, end)
	}
	cursor = end + 1
	if buf[cursor] != ';' {
		return nil, end, errors.ErrExpected(`string end ';'`, cursor)
	}

	cursor++

	return buf[start:end], cursor, nil
}

func parseByteStringInt64(b []byte) int64 {
	var l int64
	for _, c := range b {
		l = l*10 + int64(c-'0')
	}

	return l
}

func parseByteStringInt(b []byte) int {
	var l int
	for _, c := range b {
		l = l*10 + int(c-'0')
	}

	return l
}
