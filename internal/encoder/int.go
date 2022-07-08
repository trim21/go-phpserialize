package encoder

import (
	"strconv"
	"unsafe"
)

func encodeInt8(ctx *Ctx, p uintptr) error {
	value := *(*int8)(unsafe.Pointer(p))
	appendInt(ctx, int64(value))
	return nil
}

func encodeInt16(ctx *Ctx, p uintptr) error {
	value := *(*int16)(unsafe.Pointer(p))
	appendInt(ctx, int64(value))
	return nil
}

func encodeInt32(ctx *Ctx, p uintptr) error {
	value := *(*int32)(unsafe.Pointer(p))
	appendInt(ctx, int64(value))
	return nil
}

func encodeInt64(ctx *Ctx, p uintptr) error {
	value := *(*int64)(unsafe.Pointer(p))
	appendInt(ctx, int64(value))
	return nil
}

func encodeInt(ctx *Ctx, p uintptr) error {
	value := *(*int)(unsafe.Pointer(p))
	appendInt(ctx, int64(value))
	return nil
}

func appendInt(ctx *Ctx, v int64) {
	ctx.b = append(ctx.b, 'i', ':')
	ctx.b = strconv.AppendInt(ctx.b, v, 10)
	ctx.b = append(ctx.b, ';')
}
