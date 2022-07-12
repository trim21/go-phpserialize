package encoder

import "strconv"

func appendEmptyArray(b []byte) []byte {
	return append(b, `a:0:{}`...)
}

func appendNilBytes(b []byte) []byte {
	return append(b, 'N', ';')
}

func appendStringHeadBytes(b []byte, length int64) []byte {
	b = append(b, 's', ':')
	b = strconv.AppendInt(b, length, 10)
	return append(b, ':')
}

func appendEmptyArrayBytes(b []byte) []byte {
	return append(b, 'a', ':', '0', ':', '{', '}')
}

func appendArrayBeginBytes(b []byte, fieldNum int64) []byte {
	b = append(b, 'a', ':')
	b = strconv.AppendInt(b, fieldNum, 10)
	return append(b, ':', '{')
}

func appendByteString(dst, src []byte) []byte {
	dst = append(dst, 's', ':', '"')
	dst = strconv.AppendInt(dst, int64(len(src)), 10)
	return append(dst, '"', ':')
}
