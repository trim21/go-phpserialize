//go:build !race

package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileToGetCodeSet(rt reflect.Type) (encoder, error) {
	typeID := runtime.ToTypeID(rt)

	if typeID > typeAddr.MaxTypeAddr || typeID < typeAddr.BaseTypeAddr {
		return compileToGetEncoderSlowPath(typeID, rt)
	}

	index := (typeID - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	if enc := cachedEncoder[index]; enc != nil {
		return enc, nil
	}
	enc, err := compileType(rt)
	if err != nil {
		return nil, err
	}
	cachedEncoder[index] = enc
	return enc, nil
}
