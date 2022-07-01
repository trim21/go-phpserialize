package phpserialize

import (
	"reflect"
	"unsafe"
)

func str2byte(s string) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}

func byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
