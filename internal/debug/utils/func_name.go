package utils

import (
	"runtime"

	"github.com/goccy/go-reflect"
)

func FuncName(fn any) string {
	t := reflect.ValueOf(fn).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	}
	return t.String()
}
