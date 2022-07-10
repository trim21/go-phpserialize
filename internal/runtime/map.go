package runtime

import "unsafe"

type HashIter struct {
	key         unsafe.Pointer
	elem        unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer
	bptr        unsafe.Pointer
	overflow    unsafe.Pointer
	oldoverflow unsafe.Pointer
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

// new

//go:linkname MakeMap reflect.makemap
func MakeMap(*Type, int) unsafe.Pointer

// write

//go:linkname MapAssignFastStr runtime.mapassign_faststr
//go:noescape
func MapAssignFastStr(t *Type, m unsafe.Pointer, s string) unsafe.Pointer

//go:linkname MapAssign reflect.mapassign
//go:noescape
func MapAssign(t *Type, m unsafe.Pointer, k, v unsafe.Pointer)

// read

//go:linkname MapLen reflect.maplen
//go:noescape
func MapLen(m unsafe.Pointer) int

//go:linkname MapIterInit runtime.mapiterinit
//go:noescape
func MapIterInit(mapType *Type, m unsafe.Pointer, it *HashIter)

//go:linkname MapIterNext reflect.mapiternext
//go:noescape
func MapIterNext(it *HashIter)

//go:linkname MapIterValue reflect.mapiterelem
func MapIterValue(it *HashIter) uintptr

//go:linkname MapIterKey reflect.mapiterkey
//go:noescape
func MapIterKey(it *HashIter) uintptr
