//go:build !go1.19

package encoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func ptrOfPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

func reflectValueToLocal(value reflect.Value) rValue {
	return *(*rValue)(unsafe.Pointer(&value))
}

type emptyInterface struct {
	typ *runtime.Type  // value type
	ptr unsafe.Pointer // value address
}

type nonEmptyInterface struct {
	itab *struct {
		ityp *runtime.Type // static interface type
		typ  *runtime.Type // dynamic concrete type
		// unused fields...
	}
	ptr unsafe.Pointer
}

// to get reflect.Value unexported ptr field.
// used in map encoder

// to get reflect.Value unexported ptr field.
// used in map encoder
type rValue struct {
	typ  *rtype
	ptr  uintptr
	flag uintptr
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

func ifcePtrToValuePtr(p uintptr) unsafe.Pointer {
	return ((*emptyInterface)(unsafe.Pointer(p))).ptr
}
