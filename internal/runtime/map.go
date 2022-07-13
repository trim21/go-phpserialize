package runtime

import "unsafe"

type HashIter struct {
	_ unsafe.Pointer // key
	_ unsafe.Pointer // elem
	_ unsafe.Pointer // t
	_ unsafe.Pointer // h
	_ unsafe.Pointer // buckets
	_ unsafe.Pointer // bptr
	_ unsafe.Pointer // overflow
	_ unsafe.Pointer // oldoverflow
	_ uintptr        // startBucket
	_ uint8          // offset
	_ bool           // wrapped
	_ uint8          // B
	_ uint8          // i
	_ uintptr        // bucket
	_ uintptr        // checkBucket
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
