package encoder

import (
	"reflect"
	"strconv"
	"unsafe"
)

// encodeString encode string "result" to `s:6:"result";`
// encode UTF-8 string "叛逆的鲁鲁修" `s:18:"叛逆的鲁鲁修";`
// str length is underling bytes length, not len(str)
// a unsafe.Pointer(&s) is actual a pointer to reflect.StringHeader
func encodeString(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	sh := **(**reflect.StringHeader)(unsafe.Pointer(&p))
	sVal := **(**string)(unsafe.Pointer(&p))
	b = append(b, 's', ':')
	b = strconv.AppendInt(b, int64(sh.Len), 10)
	b = append(b, ':', '"')
	b = append(b, sVal...)

	return append(b, '"', ';'), nil
}

func EncodeStringPtr(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	sh := ***(***reflect.StringHeader)(unsafe.Pointer(&p))
	sVal := ***(***string)(unsafe.Pointer(&p))

	b = append(b, 's', ':')
	b = strconv.AppendInt(b, int64(sh.Len), 10)
	b = append(b, ':', '"')
	b = append(b, sVal...)

	return append(b, '"', ';'), nil
}

func compileConstStringNoError(s string) func(*Ctx, []byte) []byte {
	var finalStr = "s:" + strconv.Itoa(len(s)) + `:"` + s + `";`
	return func(ctx *Ctx, b []byte) []byte {
		return append(b, finalStr...)
	}
}
