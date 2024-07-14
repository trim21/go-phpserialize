package encoder

import (
	"reflect"
	"strconv"
)

func encodeUint(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	b = append(b, 'i', ':')
	b = strconv.AppendUint(b, rv.Uint(), 10)
	return append(b, ';'), nil
}
