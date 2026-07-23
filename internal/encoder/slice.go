package encoder

import (
	"reflect"
)

func compileSlice(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	if recursiveEnc, exists := seen[rt]; exists {
		return recursiveEnc.Encode, nil
	}
	recursiveEnc := &structRecEncoder{}
	seen[rt] = recursiveEnc

	enc, compileError := compile(rt.Elem(), seen)
	if compileError != nil {
		return nil, compileError
	}

	elemEnc := enc
	containerEnc := checkRecursiveEncoder(func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		length := rv.Len()

		b = appendArrayBegin(b, int64(length))
		var err error
		for i := 0; i < length; i++ {
			b = appendIntBytes(b, int64(i))
			b, err = elemEnc(ctx, b, rv.Index(i))
			if err != nil {
				return b, err
			}
		}
		return append(b, '}'), nil
	})
	recursiveEnc.enc = containerEnc
	return recursiveEnc.Encode, nil
}
