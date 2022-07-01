//go:build !go1.19

package encoder

import (
	stdReflect "reflect"
	"unsafe"

	"github.com/goccy/go-reflect"
)

func reflectValueMapFromPtr(typ reflect.Type, p uintptr, flag uintptr) reflect.Value {
	return *(*reflect.Value)(unsafe.Pointer(&rValueReflectType{typ: typ, ptr: p, flag: flag}))
}

func reflectValueFromPtr(typ reflect.Type, p uintptr) reflect.Value {
	return *(*reflect.Value)(unsafe.Pointer(&rValueReflectType{typ: typ, ptr: p, flag: uintptr(typ.Kind())}))
}

func reflectValueToLocal(value reflect.Value) rValue {
	return *(*rValue)(unsafe.Pointer(&value))
}

func stdReflectValueToLocal(value stdReflect.Value) rValue {
	return *(*rValue)(unsafe.Pointer(&value))
}

// to get reflect.Value unexported ptr field.
// used in map encoder
type rValueReflectType struct {
	typ  reflect.Type
	ptr  uintptr
	flag uintptr
}

// to get reflect.Value unexported ptr field.
// used in map encoder
type rValue struct {
	typ  *rtype
	ptr  uintptr
	flag uintptr
}

const PtrSize = 4 << (^uintptr(0) >> 63)

// pointer returns the underlying pointer represented by v.
// v.Kind() must be Pointer, Map, Chan, Func, or UnsafePointer
// if v.Kind() == Pointer, the base type must not be go:notinheap.
func (v rValue) pointer() unsafe.Pointer {
	if v.typ.size != PtrSize || !v.typ.pointers() {
		panic("can't call pointer on a non-pointer Value")
	}
	if v.flag&flagIndir != 0 {
		return *(*unsafe.Pointer)(unsafe.Pointer(v.ptr))
	}
	return unsafe.Pointer(v.ptr)
}

// rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      uint8   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte // garbage collection data
	str       int32 // string form
	ptrToThis int32 // type for pointer to this type, may be zero
}

func (t *rtype) pointers() bool { return t.ptrdata != 0 }

// IsValid reports whether v represents a value.
// It returns false if v is the zero Value.
// If IsValid returns false, all other methods except String panic.
// Most functions and methods never return an invalid Value.
// If one does, its documentation states the conditions explicitly.
func (v rValue) IsValid() bool {
	return v.flag != 0
}

type flag = uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)
