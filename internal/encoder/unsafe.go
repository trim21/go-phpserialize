package encoder

import (
	"unsafe"
)

func PtrDeRef(p uintptr) uintptr {
	return **(**uintptr)(unsafe.Pointer(&p))
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

type emptyInterface struct {
	typ reflect.Type   // value type
	ptr unsafe.Pointer // value address
}
