//go:build !go1.20

package encoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func PtrDeRef(p uintptr) uintptr {
	return **(**uintptr)(unsafe.Pointer(&p))
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

type emptyInterface struct {
	typ *runtime.Type  // value type
	ptr unsafe.Pointer // value address
}
