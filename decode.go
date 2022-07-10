package phpserialize

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/decoder"
	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

func unmarshal(data []byte, v any) error {
	src := make([]byte, len(data)) // append nul byte to the end
	copy(src, data)

	header := (*emptyInterface)(unsafe.Pointer(&v))

	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}
	dec, err := decoder.CompileToGetDecoder(header.typ)
	if err != nil {
		return err
	}
	ctx := decoder.TakeRuntimeContext()
	ctx.Buf = src
	cursor, err := dec.Decode(ctx, 0, 0, header.ptr)
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

func validateType(typ *runtime.Type, p uintptr) error {
	if typ == nil || typ.Kind() != reflect.Ptr || p == 0 {
		return &errors.InvalidUnmarshalError{Type: runtime.RType2Type(typ)}
	}
	return nil
}
