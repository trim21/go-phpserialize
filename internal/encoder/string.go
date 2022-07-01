package encoder

import (
	"strconv"
	"unsafe"

	"github.com/goccy/go-reflect"
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

func encodeStringVariable(buf *Ctx, p uintptr) error {
	s := (*reflect.StringHeader)(unsafe.Pointer(p))
	sVal := *(*string)(unsafe.Pointer(p))
	buf.b = append(buf.b, 's', ':')
	buf.b = strconv.AppendInt(buf.b, int64(s.Len), 10)
	buf.b = append(buf.b, ':')
	buf.b = strconv.AppendQuote(buf.b, sVal)
	buf.b = append(buf.b, ';')
	return nil
}
