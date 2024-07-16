package encoder

import (
	"reflect"
)

func compileBoolAsString(typ reflect.Type) (encoder, error) {
	return encodeBoolAsString, nil
}

func encodeBoolAsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendBoolAsString(b, rv.Bool())
}

func appendBoolAsString(b []byte, value bool) ([]byte, error) {
	if value {
		return append(b, `s:4:"true";`...), nil
	}
	return append(b, `s:5:"false";`...), nil
}
