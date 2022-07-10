//go:build !go1.19

package encoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func PtrOfPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

type emptyInterface struct {
	typ *runtime.Type  // value type
	ptr unsafe.Pointer // value address
}
