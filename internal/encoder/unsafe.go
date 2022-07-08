//go:build !go1.19

package encoder

import (
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func ptrOfPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

type emptyInterface struct {
	typ *runtime.Type  // value type
	ptr unsafe.Pointer // value address
}

// will be used to support interface with methods
type nonEmptyInterface struct {
	itab *struct {
		ityp *runtime.Type // static interface type
		typ  *runtime.Type // dynamic concrete type
		// unused fields...
	}
	ptr unsafe.Pointer
}
