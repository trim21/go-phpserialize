package phpserialize

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/decoder"
	"github.com/trim21/go-phpserialize/internal/errors"
)

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}

func Unmarshal(data []byte, v any) error {
	return unmarshal(data, v)
}

type emptyInterface struct {
	typ reflect.Type
	ptr unsafe.Pointer
}

func unmarshal(data []byte, v any) error {
	header := (*emptyInterface)(unsafe.Pointer(&v))

	rv := reflect.ValueOf(v)

	rt := rv.Type()

	if err := validateType(rt); err != nil {
		return err
	}

	src := make([]byte, len(data)) // append nul byte to the end
	copy(src, data)

	dec, err := decoder.CompileToGetDecoder(rt)
	if err != nil {
		return err
	}
	ctx := decoder.TakeRuntimeContext()
	ctx.Buf = src
	cursor, err := dec.Decode(ctx, 0, 0, rv)
	if err != nil {
		decoder.ReleaseRuntimeContext(ctx)
		return err
	}
	decoder.ReleaseRuntimeContext(ctx)
	return validateEndBuf(src, cursor)
}

func validateEndBuf(src []byte, cursor int64) error {
	if int64(len(src)) == cursor {
		return nil
	}

	return errors.ErrSyntax(
		fmt.Sprintf("invalid character '%c' after top-level value", src[cursor]),
		cursor+1,
	)
}

func validateType(typ reflect.Type) error {
	if typ == nil || typ.Kind() != reflect.Ptr {
		return &errors.InvalidUnmarshalError{Type: typ}
	}
	return nil
}
