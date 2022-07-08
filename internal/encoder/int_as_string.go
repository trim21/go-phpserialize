package encoder

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileIntAsString(typ *runtime.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.Int8:
		return encodeInt8AsString, nil
	case reflect.Int16:
		return encodeInt16AsString, nil
	case reflect.Int32:
		return encodeInt32AsString, nil
	case reflect.Int64:
		return encodeInt64AsString, nil
	case reflect.Int:
		return encodeIntAsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

func encodeInt8AsString(ctx *Ctx, p uintptr) error {
	value := *(*int8)(unsafe.Pointer(p))
	appendIntAsString(ctx, int64(value))
	return nil
}

func encodeInt16AsString(ctx *Ctx, p uintptr) error {
	value := *(*int16)(unsafe.Pointer(p))
	appendIntAsString(ctx, int64(value))
	return nil
}

func encodeInt32AsString(ctx *Ctx, p uintptr) error {
	value := *(*int32)(unsafe.Pointer(p))
	appendIntAsString(ctx, int64(value))
	return nil
}

func encodeInt64AsString(ctx *Ctx, p uintptr) error {
	value := *(*int64)(unsafe.Pointer(p))
	appendIntAsString(ctx, value)
	return nil
}

func encodeIntAsString(ctx *Ctx, p uintptr) error {
	value := *(*int)(unsafe.Pointer(p))
	appendIntAsString(ctx, int64(value))
	return nil
}

func appendIntAsString(ctx *Ctx, v int64) {
	appendStringHead(ctx, iterativeDigitsCount(v))
	ctx.b = append(ctx.b, '"')
	ctx.b = strconv.AppendInt(ctx.b, v, 10)
	ctx.b = append(ctx.b, '"', ';')
}

func iterativeDigitsCount(number int64) int64 {
	var count int64
	if number <= 0 {
		count++
	}

	for number != 0 {
		number /= 10
		count += 1
	}

	return count
}
