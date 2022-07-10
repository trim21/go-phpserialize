//go:build !race

package encoder

func compileToGetCodeSet(typeID uintptr) (encoder, error) {
	if typeID > typeAddr.MaxTypeAddr || typeID < typeAddr.BaseTypeAddr {
		return compileToGetEncoderSlowPath(typeID)
	}

	index := (typeID - typeAddr.BaseTypeAddr) >> typeAddr.AddrShift
	if enc := cachedEncoder[index]; enc != nil {
		return enc, nil
	}
	enc, err := compileTypeID(typeID)
	if err != nil {
		return nil, err
	}
	cachedEncoder[index] = enc
	return enc, nil
}
