package encoder

import (
	"reflect"
	"strconv"
)

func encodeInt(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntBytes(b, int64(rv.Int())), nil
}

func appendIntBytes(b []byte, v int64) []byte {
	b = append(b, 'i', ':')
	b = strconv.AppendInt(b, v, 10)

	return append(b, ';')
}
