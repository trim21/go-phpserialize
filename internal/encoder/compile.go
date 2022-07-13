package encoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type encoder func(ctx *Ctx, b []byte, p uintptr) ([]byte, error)

func compileTypeID(typeID uintptr) (encoder, error) {
	rt := *(**runtime.Type)(unsafe.Pointer(&typeID))

	return compile(rt)
}

func compile(rt *runtime.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return encodeBool, nil
	case reflect.Int8:
		return encodeInt8, nil
	case reflect.Int16:
		return encodeInt16, nil
	case reflect.Int32:
		return encodeInt32, nil
	case reflect.Int64:
		return encodeInt64, nil
	case reflect.Int:
		return encodeInt, nil
	case reflect.Uint8:
		return encodeUint8, nil
	case reflect.Uint16:
		return encodeUint16, nil
	case reflect.Uint32:
		return encodeUint32, nil
	case reflect.Uint64:
		return encodeUint64, nil
	case reflect.Uint:
		return encodeUint, nil
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	case reflect.String:
		return encodeString, nil
	case reflect.Struct:
		return compileStruct(rt)
	case reflect.Array:
		return compileArray(rt)
	case reflect.Slice:
		return compileSlice(rt)
	case reflect.Map:
		return compileMap(rt)
	case reflect.Interface:
		return compileInterface(rt)
	case reflect.Ptr:
		return compilePtr(rt)
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s)", rt.String(), rt.Kind())
}

func compileMapKey(typ *runtime.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.String:
		return encodeString, nil

	case reflect.Int8:
		return encodeInt8, nil
	case reflect.Int16:
		return encodeInt16, nil
	case reflect.Int32:
		return encodeInt32, nil
	case reflect.Int64:
		return encodeInt64, nil
	case reflect.Int:
		return encodeInt, nil
	case reflect.Uint8:
		return encodeUint8, nil
	case reflect.Uint16:
		return encodeUint16, nil
	case reflect.Uint32:
		return encodeUint32, nil
	case reflect.Uint64:
		return encodeUint64, nil
	case reflect.Uint:
		return encodeUint, nil
	}

	return nil, fmt.Errorf("failed to build encoder for map key, unsupported type %s (kind %s)", typ.String(), typ.Kind())
}

func compileAsString(rt *runtime.Type) (encoder, error) {
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
