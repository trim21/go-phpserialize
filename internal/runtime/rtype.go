package runtime

import (
	"reflect"
	"unsafe"
)

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

func Type2RType(t reflect.Type) reflect.Type {
	return t
}

func TypeID(rt reflect.Type) uintptr {
	p := unsafe.Pointer(&rt)
	var ef = **(**emptyInterface)(p)
	return uintptr(ef.ptr)
}

type abiType struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       uint8   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, GCData is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	GCData *byte
}
