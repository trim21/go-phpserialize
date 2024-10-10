package encoder

import (
	"fmt"
	"reflect"
)

type encoder func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error)

func compileType(rt reflect.Type) (encoder, error) {
	return compile(rt, compileSeenMap{})
}

func compile(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	if rt.Implements(marshalerType) {
		return compileMarshaler(rt)
	}

	if rt == bytesType {
		return encodeBytes, nil
	}

	switch rt.Kind() {
	case reflect.Bool:
		return encodeBool, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return encodeInt, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return encodeUint, nil
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	case reflect.String:
		return encodeString, nil
	case reflect.Struct:
		return compileStruct(rt, seen)
	case reflect.Array:
		return compileArray(rt)
	case reflect.Slice:
		return compileSlice(rt, seen)
	case reflect.Map:
		return compileMap(rt, seen)
	case reflect.Interface:
		return compileInterface(rt)
	case reflect.Ptr:
		return compilePtr(rt, seen)
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s)", rt.String(), rt.Kind())
}

func compileMapKey(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.String:
		return encodeString, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return encodeInt, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return encodeUint, nil
	}

	return nil, fmt.Errorf("failed to build encoder for map key, unsupported type %s (kind %s)", typ.String(), typ.Kind())
}

func compileAsString(rt reflect.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return compileBoolAsString(rt)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return compileIntAsString(rt)
	case reflect.String:
		return encodeString, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return compileUintAsString(rt)
	case reflect.Float32, reflect.Float64:
		return compileFloatAsString(rt)
	case reflect.Interface:
		return compileInterfaceAsString(rt)
	}

	return nil, fmt.Errorf(
		"failed to build encoder for struct field (as string), unsupported type %s (kind %s)",
		rt.String(), rt.Kind())
}
