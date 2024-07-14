package runtime

import (
	"reflect"
	"unsafe"
)

func PtrTo(rt reflect.Type) reflect.Type {
	return reflect.PointerTo(rt)
}

//go:linkname IfaceIndir reflect.ifaceIndir
//go:noescape
func IfaceIndir(any) bool

func RType2Type(t reflect.Type) reflect.Type {
	return t
}

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
