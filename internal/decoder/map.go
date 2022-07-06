package decoder

import (
	"reflect"
	"unsafe"

	"github.com/sirupsen/logrus"
	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

type mapDecoder struct {
	mapType                 *runtime.Type
	keyType                 *runtime.Type
	valueType               *runtime.Type
	canUseAssignFaststrType bool
	keyDecoder              Decoder
	valueDecoder            Decoder
	structName              string
	fieldName               string
}

func newMapDecoder(mapType *runtime.Type, keyType *runtime.Type, keyDec Decoder, valueType *runtime.Type, valueDec Decoder, structName, fieldName string) *mapDecoder {
	return &mapDecoder{
		mapType:                 mapType,
		keyDecoder:              keyDec,
		keyType:                 keyType,
		canUseAssignFaststrType: canUseAssignFaststrType(keyType, valueType),
		valueType:               valueType,
		valueDecoder:            valueDec,
		structName:              structName,
		fieldName:               fieldName,
	}
}

const (
	mapMaxElemSize = 128
)

// See detail: https://github.com/goccy/go-json/pull/283
func canUseAssignFaststrType(key *runtime.Type, value *runtime.Type) bool {
	indirectElem := value.Size() > mapMaxElemSize
	if indirectElem {
		return false
	}
	return key.Kind() == reflect.String
}

//go:linkname makemap reflect.makemap
func makemap(*runtime.Type, int) unsafe.Pointer

//nolint:golint
//go:linkname mapassign_faststr runtime.mapassign_faststr
//go:noescape
func mapassign_faststr(t *runtime.Type, m unsafe.Pointer, s string) unsafe.Pointer

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(t *runtime.Type, m unsafe.Pointer, k, v unsafe.Pointer)

func (d *mapDecoder) mapassign(t *runtime.Type, m, k, v unsafe.Pointer) {
	if d.canUseAssignFaststrType {
		mapV := mapassign_faststr(t, m, *(*string)(k))
		typedmemmove(d.valueType, mapV, v)
	} else {
		mapassign(t, m, k, v)
	}
}

func (d *mapDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	depth++
	if depth > maxDecodeNestingDepth {
		return errors.ErrExceededMaxDepth(s.char(), s.cursor)
	}

	switch s.skipWhiteSpace() {
	case 'n':
		if err := nullBytes(s); err != nil {
			return err
		}
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = nil
		return nil
	case '{':
	default:
		return errors.ErrExpected("{ character for map value", s.totalOffset())
	}
	mapValue := *(*unsafe.Pointer)(p)
	if mapValue == nil {
		mapValue = makemap(d.mapType, 0)
	}
	s.cursor++
	if s.equalChar('}') {
		*(*unsafe.Pointer)(p) = mapValue
		s.cursor++
		return nil
	}
	for {
		k := unsafe_New(d.keyType)
		if err := d.keyDecoder.DecodeStream(s, depth, k); err != nil {
			return err
		}
		s.skipWhiteSpace()
		if !s.equalChar(':') {
			return errors.ErrExpected("colon after object key", s.totalOffset())
		}
		s.cursor++
		v := unsafe_New(d.valueType)
		if err := d.valueDecoder.DecodeStream(s, depth, v); err != nil {
			return err
		}
		d.mapassign(d.mapType, mapValue, k, v)
		s.skipWhiteSpace()
		if s.equalChar('}') {
			**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
			s.cursor++
			return nil
		}
		if !s.equalChar(',') {
			return errors.ErrExpected("comma after object value", s.totalOffset())
		}
		s.cursor++
	}
}

// TODO
func (d *mapDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	logrus.Info("decode map")
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	buflen := int64(len(buf))
	if buflen < 2 {
		return 0, errors.ErrExpected("{} for map", cursor)
	}
	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = nil
		return cursor, nil
	case 'a':
		// array case
		cursor++
	default:
		return 0, errors.ErrExpected("{ character for map value", cursor)
	}

	l, end, err := readLength(buf, cursor)
	if err != nil {
		return 0, err
	}

	cursor = end
	if buf[cursor] != '{' {
		return 0, errors.ErrExpected("{ character for map value", cursor)
	}

	mapValue := *(*unsafe.Pointer)(p)
	if mapValue == nil {
		mapValue = makemap(d.mapType, int(l))
	}

	cursor++
	if buf[cursor] == '}' {
		**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
		cursor++
		return cursor, nil
	}

	for {
		k := unsafe_New(d.keyType)
		keyCursor, err := d.keyDecoder.Decode(ctx, cursor, depth, k)
		if err != nil {
			return 0, err
		}
		cursor = keyCursor
		v := unsafe_New(d.valueType)
		valueCursor, err := d.valueDecoder.Decode(ctx, cursor, depth, v)
		if err != nil {
			return 0, err
		}
		d.mapassign(d.mapType, mapValue, k, v)
		cursor = valueCursor
		if buf[cursor] == '}' {
			**(**unsafe.Pointer)(unsafe.Pointer(&p)) = mapValue
			cursor++
			return cursor, nil
		}
	}
}
