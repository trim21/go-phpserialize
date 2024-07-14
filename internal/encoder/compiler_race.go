//go:build race

package encoder

import (
	"reflect"
	"sync"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

var setsMu sync.RWMutex

func compileToGetCodeSet(rt reflect.Type) (encoder, error) {
	typeID := runtime.ToTypeID(rt)

	if typeID > typeAddr.MaxTypeAddr || typeID < typeAddr.BaseTypeAddr {
		return compileToGetEncoderSlowPath(typeID, rt)
	}
	index := (typeID - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	setsMu.RLock()
	if codeSet := cachedEncoder[index]; codeSet != nil {
		encoder, err := compileType(typeID)
		if err != nil {
			setsMu.RUnlock()
			return nil, err
		}
		setsMu.RUnlock()
		return encoder, nil
	}
	setsMu.RUnlock()

	encoder, err := compileType(rt)
	if err != nil {
		return nil, err
	}
	setsMu.Lock()
	cachedEncoder[index] = encoder
	setsMu.Unlock()
	return encoder, nil
}
