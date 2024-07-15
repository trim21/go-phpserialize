package decoder

import (
	"reflect"
)

type Decoder interface {
	Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error)
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

var (
	unmarshalPHPType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)
