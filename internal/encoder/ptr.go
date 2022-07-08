package encoder

import (
	"github.com/goccy/go-reflect"
)

// TODO: why ptr work fine without an encoder?
func compilePtr(rt reflect.Type) (encoder, error) {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	return compile(rt)
}
