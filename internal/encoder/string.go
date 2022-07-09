package encoder

import (
	"reflect"
	"strconv"
	"unsafe"
)

// encode string "result" to `s:6:"result";`
// encode UTF-8 string "叛逆的鲁鲁修" `s:18:"叛逆的鲁鲁修";`
// str length is underling bytes length, not len(str)
func compileConstString(s string) (encoder, error) {
	var encodedStr = "s:" + strconv.Itoa(len(s)) + ":" + strconv.Quote(s) + ";"
	return func(buf *Ctx, p uintptr) error {
		buf.b = append(buf.b, encodedStr...)
		return nil
	}, nil
}

func compileConstStringNoError(s string) func(*Ctx) {
	var encodedStr = "s:" + strconv.Itoa(len(s)) + ":" + strconv.Quote(s) + ";"
	return func(ctx *Ctx) {
		ctx.b = append(ctx.b, encodedStr...)
	}
}

func encodeStringVariable(ctx *Ctx, p uintptr) error {
	s := (*reflect.StringHeader)(unsafe.Pointer(p))
	sVal := *(*string)(unsafe.Pointer(p))
	ctx.b = append(ctx.b, 's', ':')
	ctx.b = strconv.AppendInt(ctx.b, int64(s.Len), 10)
	ctx.b = append(ctx.b, ':')
	ctx.b = strconv.AppendQuote(ctx.b, sVal)
	ctx.b = append(ctx.b, ';')
	return nil
}
