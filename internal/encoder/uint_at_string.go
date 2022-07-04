package encoder

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
)

func compileUintAsString(rt reflect.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Uint8:
		return encodeUint8AsString, nil
	case reflect.Uint16:
		return encodeUint16AsString, nil
	case reflect.Uint32:
		return encodeUint32AsString, nil
	case reflect.Uint64:
		return encodeUint64AsString, nil
	case reflect.Uint:
		return encodeUintAsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", rt.Kind()))
}

func encodeUint8AsString(buf *Ctx, p uintptr) error {
	value := *(*uint8)(unsafe.Pointer(p))
	appendUintAsString(buf, uint64(value))
	return nil
}

func encodeUint16AsString(buf *Ctx, p uintptr) error {
	value := *(*uint16)(unsafe.Pointer(p))
	appendUintAsString(buf, uint64(value))
	return nil
}

func encodeUint32AsString(buf *Ctx, p uintptr) error {
	value := *(*uint32)(unsafe.Pointer(p))
	appendUintAsString(buf, uint64(value))
	return nil
}

func encodeUint64AsString(buf *Ctx, p uintptr) error {
	value := *(*uint64)(unsafe.Pointer(p))
	appendUintAsString(buf, value)
	return nil
}

func encodeUintAsString(buf *Ctx, p uintptr) error {
	value := *(*uint)(unsafe.Pointer(p))
	appendUintAsString(buf, uint64(value))
	return nil
}

func appendUintAsString(ctx *Ctx, v uint64) {
	appendStringHead(ctx, uintDigitsCount(v))

	ctx.b = append(ctx.b, '"')
	ctx.b = strconv.AppendUint(ctx.b, v, 10)
	ctx.b = append(ctx.b, '"', ';')
}

func uintDigitsCount(number uint64) int64 {
	var count int64
	if number == 0 {
		return 1
	}

	for number != 0 {
		number /= 10
		count += 1
	}

	return count
}
