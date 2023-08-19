package decoder

import (
	"encoding"
	"reflect"
	"unsafe"
)

type Decoder interface {
	Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error)
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

var (
	unmarshalPHPType             = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	unmarshalTextUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)
