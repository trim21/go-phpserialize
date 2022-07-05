package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/ifce"
)

type Decoder interface {
	Decode(*RuntimeContext, int64, int64, unsafe.Pointer) (int64, error)
	DecodeStream(*Stream, int64, unsafe.Pointer) error
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

var (
	unmarshalPHPType = reflect.TypeOf((*ifce.Unmarshaler)(nil)).Elem()
)
