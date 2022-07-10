package encoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type isEmpty func(ctx *Ctx, p uintptr) (isEmpty bool, err error)

func notIgnore(ctx *Ctx, p uintptr) (isEmpty bool, err error) {
	return false, nil
}

func compileEmptyer(rt *runtime.Type) (isEmpty, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**bool)(unsafe.Pointer(&p))
			return value == false, nil
		}, nil
	case reflect.Int8:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**int8)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Int16:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**int16)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Int32:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**int32)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Int64:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**int64)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Int:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**int)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.String:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			s := (*reflect.StringHeader)(unsafe.Pointer(p))
			return s.Len == 0, nil
		}, nil
	case reflect.Uint8:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**uint8)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Uint16:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**uint16)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Uint32:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**uint32)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Uint64:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**uint64)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Uint:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**uint)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Float32:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**float32)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Float64:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			value := **(**float64)(unsafe.Pointer(&p))
			return value == 0, nil
		}, nil
	case reflect.Struct:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			return false, nil
		}, nil
	case reflect.Slice:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			dataPtr := **(**uintptr)(unsafe.Pointer(&p))
			return dataPtr == 0, nil
		}, nil
	case reflect.Map:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			return p == 0, nil
		}, nil
	case reflect.Interface:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			p = unpackIface(p)
			return p == 0, nil
		}, nil
	case reflect.Ptr:
		return func(ctx *Ctx, p uintptr) (bool, error) {
			return p == 0, nil
		}, nil
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s) with tag `omitempty`", rt.String(), rt.Kind())
}
