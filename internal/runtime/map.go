package runtime

import "unsafe"

//go:linkname MapLen reflect.maplen
//go:noescape
func MapLen(m unsafe.Pointer) int
