package runtime

import (
	"unsafe"
)

type SliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

const (
	maxAcceptableTypeAddrRange = 1024 * 1024 * 2 // 2 Mib
)

type TypeAddr struct {
	BaseTypeAddr uintptr
	MaxTypeAddr  uintptr
	AddrRange    uintptr
	AddrShift    uintptr
}

var (
	typeAddr        *TypeAddr
	alreadyAnalyzed bool
)

func AnalyzeTypeAddr() *TypeAddr {
	defer func() {
		alreadyAnalyzed = true
	}()
	if alreadyAnalyzed {
		return typeAddr
	}
	return typeAddr
}
