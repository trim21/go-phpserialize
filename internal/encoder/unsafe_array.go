package encoder

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

func reflectValueIndex(v rValue, i int) reflect.Value {
	s := (*reflect.SliceHeader)(unsafe.Pointer(v.ptr))
	tt := (*sliceType)(unsafe.Pointer(v.typ))
	typ := tt.elem
	val := s.Data + uintptr(i)*typ.size
	fl := flagAddr | flagIndir | ro(v.flag) | flag(typ.Kind())
	return *(*reflect.Value)(unsafe.Pointer(&rValue{typ, val, fl}))
}

// sliceType represents a slice type.
type sliceType struct {
	rtype
	elem *rtype // slice element type
}

func ro(f uintptr) flag {
	if f&flagRO != 0 {
		return flagStickyRO
	}
	return 0
}
