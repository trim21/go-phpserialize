package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileBoolAsString(typ reflect.Type) (encoder, error) {
	return encodeBoolAsString, nil
}

func encodeBoolAsString(buf *Ctx, p uintptr) error {
	value := *(*bool)(unsafe.Pointer(p))
	appendBoolAsString(buf, value)
	return nil
}

func appendBoolAsString(buf *Ctx, v bool) {
	if v {
		buf.b = append(buf.b, `s:4:"true";`...)
	} else {
		buf.b = append(buf.b, `s:5:"false";`...)
	}
}
