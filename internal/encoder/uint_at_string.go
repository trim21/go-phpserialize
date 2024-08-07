package encoder

import (
	"fmt"
	"reflect"
	"strconv"
)

func compileUintAsString(rt reflect.Type) (encoder, error) {
	switch rt.Kind() {
	case reflect.Uint8:
		return encodeUint8AsString, nil
	case reflect.Uint16:
		return encodeUint16AsString, nil
	case reflect.Uint32:
		return encodeUint32AsString, nil
	case reflect.Uint64:
		return encodeUint64AsString, nil
	case reflect.Uint:
		return encodeUintAsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", rt.Kind()))
}

func encodeUint8AsString(buf *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendUintAsString(b, rv.Uint())

}

func encodeUint16AsString(buf *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendUintAsString(b, rv.Uint())

}

func encodeUint32AsString(buf *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendUintAsString(b, rv.Uint())

}

func encodeUint64AsString(buf *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendUintAsString(b, rv.Uint())

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
