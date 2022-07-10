//go:build race

package encoder

import (
	"sync"
)

var setsMu sync.RWMutex

func compileToGetCodeSet(typeID uintptr) (encoder, error) {
	if typeID > typeAddr.MaxTypeAddr || typeID < typeAddr.BaseTypeAddr {
		return compileToGetEncoderSlowPath(typeID)
	}
	index := (typeID - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	setsMu.RLock()
	if codeSet := cachedEncoder[index]; codeSet != nil {
		encoder, err := compileTypeID(typeID)
		if err != nil {
			setsMu.RUnlock()
			return nil, err
		}
		setsMu.RUnlock()
		return encoder, nil
	}
	setsMu.RUnlock()

	encoder, err := compileTypeID(typeID)
	if err != nil {
		return nil, err
	}
	setsMu.Lock()
	cachedEncoder[index] = encoder
	setsMu.Unlock()
	return encoder, nil
}
