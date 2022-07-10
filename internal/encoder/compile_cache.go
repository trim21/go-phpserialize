package encoder

import (
	"sync/atomic"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

var (
	cachedEncoder    []encoder
	cachedEncoderMap unsafe.Pointer // map[uintptr]*OpcodeSet
	typeAddr         *runtime.TypeAddr
)

func init() {
	typeAddr = runtime.AnalyzeTypeAddr()
	if typeAddr == nil {
		typeAddr = &runtime.TypeAddr{}
	}

	cachedEncoder = make([]encoder, typeAddr.AddrRange>>typeAddr.AddrShift+1)
}

func compileToGetEncoderSlowPath(typeID uintptr) (encoder, error) {
	opcodeMap := loadEncoderMap()
	if codeSet, exists := opcodeMap[typeID]; exists {
		return codeSet, nil
	}
	codeSet, err := compileTypeID(typeID)
	if err != nil {
		return nil, err
	}
	storeEncoder(typeID, codeSet, opcodeMap)
	return codeSet, nil
}

func loadEncoderMap() map[uintptr]encoder {
	p := atomic.LoadPointer(&cachedEncoderMap)
	return *(*map[uintptr]encoder)(unsafe.Pointer(&p))
}

func storeEncoder(typ uintptr, set encoder, m map[uintptr]encoder) {
	newEncoderMap := make(map[uintptr]encoder, len(m)+1)
	newEncoderMap[typ] = set

	for k, v := range m {
		newEncoderMap[k] = v
	}

	atomic.StorePointer(&cachedEncoderMap, *(*unsafe.Pointer)(unsafe.Pointer(&newEncoderMap)))
}

func compileWithCache(rt *runtime.Type) (encoder, error) {
	typeID := uintptr(unsafe.Pointer(rt))
	return compileToGetCodeSet(typeID)
}

func compileTypeIDWithCache(typeID uintptr) (encoder, error) {
	return compileToGetCodeSet(typeID)
}
