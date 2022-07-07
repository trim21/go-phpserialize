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

func skipObject(buf []byte, cursor, depth int64) (int64, error) {
	braceCount := 1
	for {
		switch buf[cursor] {
		case '{':
			braceCount++
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
			}
		case '}':
			depth--
			braceCount--
			if braceCount == 0 {
				return cursor + 1, nil
			}
		case '[':
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
			}
		case ']':
			depth--
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '\\':
					cursor++
					if buf[cursor] == nul {
						return 0, errors.ErrUnexpectedEnd("string of object", cursor)
					}
				case '"':
					goto SWITCH_OUT
				case nul:
					return 0, errors.ErrUnexpectedEnd("string of object", cursor)
				}
			}
		case nul:
			return 0, errors.ErrUnexpectedEnd("object of object", cursor)
		}
	SWITCH_OUT:
		cursor++
	}
}

func skipArray(buf []byte, cursor, depth int64) (int64, error) {
	bracketCount := 1
	for {
		switch buf[cursor] {
		case '[':
			bracketCount++
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
			}
		case ']':
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
		case '}':
			depth--
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '\\':
					cursor++
					if buf[cursor] == nul {
						return 0, errors.ErrUnexpectedEnd("string of object", cursor)
					}
				case '"':
					goto SWITCH_OUT
				case nul:
					return 0, errors.ErrUnexpectedEnd("string of object", cursor)
				}
			}
		case nul:
			return 0, errors.ErrUnexpectedEnd("array of object", cursor)
		}
	SWITCH_OUT:
		cursor++
	}
}

func skipValue(buf []byte, cursor, depth int64) (int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\t', '\n', '\r':
			cursor++
			continue
		case '{':
			return skipObject(buf, cursor+1, depth+1)
		case '[':
			return skipArray(buf, cursor+1, depth+1)
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '\\':
					cursor++
					if buf[cursor] == nul {
						return 0, errors.ErrUnexpectedEnd("string of object", cursor)
					}
				case '"':
					return cursor + 1, nil
				case nul:
					return 0, errors.ErrUnexpectedEnd("string of object", cursor)
				}
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			for {
				cursor++
				if floatTable[buf[cursor]] {
					continue
				}
				break
			}
			return cursor, nil
		case 't':
			if err := validateTrue(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 4
			return cursor, nil
		case 'f':
			if err := validateFalse(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 5
			return cursor, nil
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 4
			return cursor, nil
		default:
			return cursor, errors.ErrUnexpectedEnd("null", cursor)
		}
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
