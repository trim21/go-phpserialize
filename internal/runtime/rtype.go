package runtime

import (
	"reflect"
	"unsafe"
)

//go:nolint structcheck
type emptyInterface struct {
	_   uintptr
	ptr unsafe.Pointer
}

func TypeID(rt reflect.Type) uintptr {
	p := unsafe.Pointer(&rt)
	var ef = *(*emptyInterface)(p)
	return uintptr(ef.ptr)
}
