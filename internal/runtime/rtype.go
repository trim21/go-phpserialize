package runtime

import (
	"reflect"
	"unsafe"
)

func IfaceIndir(rt reflect.Type) bool {
	return ifaceIndir(unsafe.Pointer(TypeID(rt)))
}

//go:linkname ifaceIndir reflect.ifaceIndir
//go:noescape
func ifaceIndir(p unsafe.Pointer) bool

//go:nolint structcheck
type emptyInterface struct {
	_   uintptr
	ptr unsafe.Pointer
}

func Type2RType(t reflect.Type) reflect.Type {
	return t
}

func TypeID(rt reflect.Type) uintptr {
	p := unsafe.Pointer(&rt)
	var ef = **(**emptyInterface)(p)
	return uintptr(ef.ptr)
}
