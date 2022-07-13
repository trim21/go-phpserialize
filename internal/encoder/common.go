package encoder

import "strconv"

func appendEmptyArray(b []byte) []byte {
	return append(b, "a:0:{}"...)
}

func appendNull(b []byte) []byte {
	return append(b, 'N', ';')
}

func appendStringHead(b []byte, length int64) []byte {
	b = append(b, 's', ':')
	b = strconv.AppendInt(b, length, 10)
	return append(b, ':')
}

func appendArrayBegin(b []byte, fieldNum int64) []byte {
	b = append(b, 'a', ':')
	b = strconv.AppendInt(b, fieldNum, 10)
	return append(b, ':', '{')
}

func appendBytesAsPhpStringVariable(dst, src []byte) []byte {
	dst = append(dst, 's', ':', '"')
	dst = strconv.AppendInt(dst, int64(len(src)), 10)
	return append(dst, '"', ':')
}
