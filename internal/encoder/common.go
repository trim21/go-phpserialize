package encoder

import "strconv"

func appendArrayBegin(buf *Ctx, fieldNum int64) {
	buf.b = append(buf.b, 'a', ':')
	buf.b = strconv.AppendInt(buf.b, fieldNum, 10)
	buf.b = append(buf.b, ':', '{')
}

func appendString(buf *Ctx, s string) {
	buf.b = append(buf.b, 's', ':')
	buf.b = strconv.AppendInt(buf.b, int64(len(s)), 10)
	buf.b = append(buf.b, ':')
	buf.b = strconv.AppendQuote(buf.b, s)
	buf.b = append(buf.b, ';')
}
