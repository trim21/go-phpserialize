package runtime

import (
	"reflect"
	"unsafe"
)

// IfaceIndir !! find a way to not using this
func IfaceIndir(rt reflect.Type) bool {
	t := (*abiType)(unsafe.Pointer(TypeID(rt)))

	return ifaceIndir(t)
}

//go:linkname ifaceIndir reflect.ifaceIndir
//go:noescape
func ifaceIndir(t *abiType) bool

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

type abiType struct {
}
