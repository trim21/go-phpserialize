package encoder

import (
	"reflect"
	"strconv"
)

// encodeString encode string "result" to `s:6:"result";`
// encode UTF-8 string "叛逆的鲁鲁修" `s:18:"叛逆的鲁鲁修";`
// str length is underling bytes length, not len(str)
func encodeString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendPhpStringVariable(ctx, b, rv.String()), nil
}

func appendPhpStringVariable(ctx *Ctx, b []byte, s string) []byte {
	b = append(b, 's', ':')
	b = strconv.AppendInt(b, int64(len(s)), 10)
	b = append(b, ':', '"')
	b = append(b, s...)

	return append(b, '"', ';')
}
