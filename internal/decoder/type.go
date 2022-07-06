package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/ifce"
)

type Decoder interface {
	Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error)
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

var (
	unmarshalPHPType = reflect.TypeOf((*ifce.Unmarshaler)(nil)).Elem()
)
