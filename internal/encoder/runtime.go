package encoder

import (
	"unsafe"
)

// without doing any allocations.
type hiter struct {
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

//go:linkname mapIterInit runtime.mapiterinit
//go:noescape
func mapIterInit(mapType reflect.Type, m unsafe.Pointer, it *hiter)

//go:linkname mapIterNext reflect.mapiternext
//go:noescape
func mapIterNext(it *hiter)

//go:linkname mapIterKey reflect.mapiterkey
//go:noescape
func mapIterKey(it *hiter) unsafe.Pointer

//go:linkname mapIterValue reflect.mapiterelem
func mapIterValue(it *hiter) unsafe.Pointer

func (h *hiter) initialized() bool {
	return h.t != nil
}

// A mapIter is an iterator for ranging over a map.
// See ValueUnsafeAddress.MapRange.
type mapIter struct {
	hiter hiter
}

func (iter *mapIter) reset() {
	iter.hiter = hiter{}
}

// mapType represents a map type.
type mapType struct {
	rtype
	key    *rtype // map key type
	elem   *rtype // map element (value) type
	bucket *rtype // internal bucket structure
	// function for hashing keys (ptr to key, seed) -> hash
	hasher     func(unsafe.Pointer, uintptr) uintptr
	keysize    uint8  // size of key slot
	valuesize  uint8  // size of value slot
	bucketsize uint16 // size of bucket
	flags      uint32
}
