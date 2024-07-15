package encoder

import (
	"reflect"
	"strconv"
)

var bytesType = reflect.TypeOf([]byte{})

func encodeBytes(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	b = append(b, 's', ':')
	b = strconv.AppendInt(b, int64(rv.Len()), 10)
	b = append(b, ':', '"')
	b = append(b, rv.Bytes()...)

	return append(b, '"', ';'), nil
}
