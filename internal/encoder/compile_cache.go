package encoder

import (
	"reflect"
	"sync/atomic"
)

var (
	cachedEncoderMap atomic.Pointer[map[reflect.Type]encoder]
)

func init() {
	var m = map[reflect.Type]encoder{}
	cachedEncoderMap.Store(&m)
}

func compileToGetCodeSet(rt reflect.Type) (encoder, error) {
	return compileToGetEncoderSlowPath(rt)
}

func compileToGetEncoderSlowPath(rt reflect.Type) (encoder, error) {
	opcodeMap := *cachedEncoderMap.Load()
	if codeSet, exists := opcodeMap[rt]; exists {
		return codeSet, nil
	}
	codeSet, err := compileType(rt)
	if err != nil {
		return nil, err
	}
	storeEncoder(rt, codeSet, opcodeMap)
	return codeSet, nil
}

func storeEncoder(rt reflect.Type, set encoder, m map[reflect.Type]encoder) {
	newEncoderMap := make(map[reflect.Type]encoder, len(m)+1)
	newEncoderMap[rt] = set

	for k, v := range m {
		newEncoderMap[k] = v
	}

	cachedEncoderMap.Store(&newEncoderMap)
}

func compileWithCache(rt reflect.Type) (encoder, error) {
	return compileToGetCodeSet(rt)
}
