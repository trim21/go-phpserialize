package encoder

import (
	"unsafe"
)

type emptyInterface struct {
	typ uintptr        // value type
	ptr unsafe.Pointer // value address
}
