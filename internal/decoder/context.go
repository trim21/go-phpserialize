package decoder

import (
	"strconv"
	"sync"
	"unsafe"

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

func char(ptr unsafe.Pointer, offset int64) byte {
	return *(*byte)(unsafe.Add(ptr, offset))
}

// `:${length}:`
func skipLengthWithBothColon(buf []byte, cursor int64) (int64, error) {
	if buf[cursor] != ':' {
		return cursor, errors.ErrUnexpected("':' before length", cursor, buf[cursor])
	}
	cursor++

	for isInteger[buf[cursor]] {
		cursor++
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
	end, err := skipString(buf, cursor)
	return end - 2, err
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

func readBool(buf []byte, cursor int64) (bool, error) {
	if cursor+1 >= int64(len(buf)) {
		return false, errors.ErrUnexpectedEnd("null", cursor)
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
