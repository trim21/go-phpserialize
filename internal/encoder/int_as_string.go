package encoder

import (
	"fmt"
	"reflect"
	"strconv"
)

func compileIntAsString(typ reflect.Type) (encoder, error) {
	switch typ.Kind() {
	case reflect.Int8:
		return encodeInt8AsString, nil
	case reflect.Int16:
		return encodeInt16AsString, nil
	case reflect.Int32:
		return encodeInt32AsString, nil
	case reflect.Int64:
		return encodeInt64AsString, nil
	case reflect.Int:
		return encodeIntAsString, nil
	}

	panic(fmt.Sprintf("unexpected kind %s", typ.Kind()))
}

func encodeInt8AsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntAsString(b, rv.Int())
}

func encodeInt16AsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntAsString(b, rv.Int())
}

func encodeInt32AsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntAsString(b, rv.Int())
}

func encodeInt64AsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntAsString(b, rv.Int())
}

func encodeIntAsString(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return appendIntAsString(b, rv.Int())
}

func appendIntAsString(b []byte, v int64) ([]byte, error) {
	b = appendStringHead(b, iterativeDigitsCount(v))
	b = append(b, '"')
	b = strconv.AppendInt(b, v, 10)
	return append(b, '"', ';'), nil
}

func iterativeDigitsCount(number int64) int64 {
	var count int64
	if number <= 0 {
		count++
	}

	for number != 0 {
		number /= 10
		count += 1
	}

	return count
}
