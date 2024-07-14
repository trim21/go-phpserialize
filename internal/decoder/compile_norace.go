//go:build !race
// +build !race

package decoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func CompileToGetDecoder(rt reflect.Type) (Decoder, error) {
	typeptr := runtime.TypeID(rt)
	if typeptr > typeAddr.MaxTypeAddr {
		return compileToGetDecoderSlowPath(typeptr, rt)
	}

	index := (typeptr - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	if dec := cachedDecoder[index]; dec != nil {
		return dec, nil
	}

	dec, err := compileHead(rt, map[uintptr]Decoder{})
	if err != nil {
		return nil, err
	}
	cachedDecoder[index] = dec
	return dec, nil
}
