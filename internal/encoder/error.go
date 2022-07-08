package encoder

import (
	"reflect"
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "php serialize: can't encode type: " + e.Type.String()
}

type UnsupportedTypeAsMapKeyError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeAsMapKeyError) Error() string {
	return "php serialize: unsupported type as key of map: " + e.Type.String()
}

type UnsupportedInterfaceTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedInterfaceTypeError) Error() string {
	return "php serialize: can't encode type (as part of an interface): " + e.Type.String()
}
