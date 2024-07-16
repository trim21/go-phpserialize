package decoder

import (
	"unsafe"
)

func unsafeStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
