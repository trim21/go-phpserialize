package encoder

import (
	"reflect"
	"strconv"
)

func compileUintAsString(rt reflect.Type) (encoder, error) {
	return encodeUintAsString, nil
}

func encodeUintAsString(buf *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendUintAsString(b, rv.Uint())
}

func appendUintAsString(b []byte, v uint64) ([]byte, error) {
	b = appendStringHead(b, uintDigitsCount(v))
	b = append(b, '"')
	b = strconv.AppendUint(b, v, 10)
	b = append(b, '"', ';')
	return b, nil
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
