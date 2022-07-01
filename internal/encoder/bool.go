package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileBool(typ reflect.Type) (encoder, error) {
	return encodeBool, nil
}

func encodeBool(buf *Ctx, p uintptr) error {
	value := *(*bool)(unsafe.Pointer(p))
	appendBool(buf, value)
	return nil
}

func appendBool(buf *Ctx, v bool) {
	buf.b = append(buf.b, 'b', ':')
	if v {
		buf.b = append(buf.b, '1')
	} else {
		buf.b = append(buf.b, '0')
	}
	buf.b = append(buf.b, ';')
}
