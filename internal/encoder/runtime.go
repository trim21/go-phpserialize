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
	overflow    *[]unsafe.Pointer
	oldoverflow *[]unsafe.Pointer
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
func mapIterInit(mapType *rtype, m unsafe.Pointer, it *hiter)

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
// See Value.MapRange.
type mapIter struct {
	m     rValue
	hiter hiter
}

func (iter *mapIter) reset() {
	iter.hiter = hiter{}
}

// Key returns the key of iter's current map entry without alloc.
func (iter *mapIter) Key() uintptr {
	iterKey := mapIterKey(&iter.hiter)
	return uintptr(iterKey)
}

// Value returns the value of iter's current map entry.
func (iter *mapIter) Value() uintptr {
	return uintptr(mapIterValue(&iter.hiter))
}

// Next advances the map iterator and reports whether there is another
// entry. It returns false when iter is exhausted; subsequent
// calls to Key, Value, or Next will panic.
func (iter *mapIter) Next() bool {
	if !iter.m.IsValid() {
		panic("mapIter.Next called on an iterator that does not have an associated map Value")
	}
	if !iter.hiter.initialized() {
		mapIterInit(iter.m.typ, iter.m.pointer(), &iter.hiter)
	} else {
		if mapIterKey(&iter.hiter) == nil {
			panic("mapIter.Next called on exhausted iterator")
		}
		mapIterNext(&iter.hiter)
	}
	return mapIterKey(&iter.hiter) != nil
}
