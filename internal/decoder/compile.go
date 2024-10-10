package decoder

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

var (
	cachedDecoderMap atomic.Pointer[map[reflect.Type]Decoder]
)

func init() {
	var m = map[reflect.Type]Decoder{}
	cachedDecoderMap.Store(&m)
}

func CompileToGetDecoder(rt reflect.Type) (Decoder, error) {
	decoderMap := *cachedDecoderMap.Load()
	if dec, exists := decoderMap[rt]; exists {
		return dec, nil
	}

	dec, err := compileHead(rt, map[reflect.Type]Decoder{})
	if err != nil {
		return nil, err
	}

	storeDecoder(rt, dec, decoderMap)

	return dec, nil
}

func storeDecoder(rt reflect.Type, dec Decoder, m map[reflect.Type]Decoder) {
	newDecoderMap := make(map[reflect.Type]Decoder, len(m)+1)
	for k, v := range m {
		newDecoderMap[k] = v
	}

	newDecoderMap[rt] = dec

	cachedDecoderMap.Store(&newDecoderMap)
}

func compileHead(rt reflect.Type, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	if reflect.PointerTo(rt).Implements(unmarshalPHPType) {
		return newUnmarshalTextDecoder(reflect.PointerTo(rt), "", ""), nil
	}
	return compile(rt.Elem(), "", "", structTypeToDecoder)
}

func compile(rt reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	switch {
	case reflect.PointerTo(rt).Implements(unmarshalPHPType):
		return newUnmarshalTextDecoder(reflect.PointerTo(rt), structName, fieldName), nil
	}

	switch rt.Kind() {
	case reflect.Ptr:
		return compilePtr(rt, structName, fieldName, structTypeToDecoder)
	case reflect.Struct:
		return compileStruct(rt, structName, fieldName, structTypeToDecoder)
	case reflect.Slice:
		elem := rt.Elem()
		if elem.Kind() == reflect.Uint8 {
			return compileBytes(elem, structName, fieldName)
		}
		return compileSlice(rt, structName, fieldName, structTypeToDecoder)
	case reflect.Array:
		return compileArray(rt, structName, fieldName, structTypeToDecoder)
	case reflect.Map:
		return compileMap(rt, structName, fieldName, structTypeToDecoder)
	case reflect.Interface:
		return compileInterface(rt, structName, fieldName)
	case reflect.Uintptr:
		return compileUint(rt, structName, fieldName)
	case reflect.Int:
		return compileInt(rt, structName, fieldName)
	case reflect.Int8:
		return compileInt8(rt, structName, fieldName)
	case reflect.Int16:
		return compileInt16(rt, structName, fieldName)
	case reflect.Int32:
		return compileInt32(rt, structName, fieldName)
	case reflect.Int64:
		return compileInt64(rt, structName, fieldName)
	case reflect.Uint:
		return compileUint(rt, structName, fieldName)
	case reflect.Uint8:
		return compileUint8(rt, structName, fieldName)
	case reflect.Uint16:
		return compileUint16(rt, structName, fieldName)
	case reflect.Uint32:
		return compileUint32(rt, structName, fieldName)
	case reflect.Uint64:
		return compileUint64(rt, structName, fieldName)
	case reflect.String:
		return compileString(rt, structName, fieldName)
	case reflect.Bool:
		return compileBool(structName, fieldName)
	case reflect.Float32:
		return compileFloat32(structName, fieldName)
	case reflect.Float64:
		return compileFloat64(structName, fieldName)
	}
	return newInvalidDecoder(rt, structName, fieldName), nil
}

func isStringTagSupportedType(typ reflect.Type) bool {
	switch {
	case reflect.PointerTo(typ).Implements(unmarshalPHPType):
		return false
	}
	switch typ.Kind() {
	case reflect.Map:
		return false
	case reflect.Slice:
		return false
	case reflect.Array:
		return false
	case reflect.Struct:
		return false
	case reflect.Interface:
		return false
	}
	return true
}

func compileMapKey(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	if reflect.PointerTo(typ).Implements(unmarshalPHPType) {
		return newUnmarshalTextDecoder(reflect.PointerTo(typ), structName, fieldName), nil
	}
	if typ.Kind() == reflect.String {
		return newStringDecoder(structName, fieldName), nil
	}
	dec, err := compile(typ, structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}

	for {
		switch t := dec.(type) {
		case *stringDecoder, *interfaceDecoder:
			return dec, nil
		case *boolDecoder, *intDecoder, *uintDecoder:
			return dec, nil
		case *ptrDecoder:
			dec = t.dec
		default:
			return newInvalidDecoder(typ, structName, fieldName), nil
		}
	}
}

func compilePtr(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	dec, err := compile(typ.Elem(), structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newPtrDecoder(dec, typ.Elem(), structName, fieldName)
}

func compileInt(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName), nil
}

func compileInt8(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName), nil
}

func compileInt16(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName), nil
}

func compileInt32(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName), nil
}

func compileInt64(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName), nil
}

func compileUint(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName), nil
}

func compileUint8(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName), nil
}

func compileUint16(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName), nil
}

func compileUint32(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName), nil
}

func compileUint64(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName), nil
}

func compileFloat32(structName, fieldName string) (Decoder, error) {
	return newFloatDecoder(structName, fieldName), nil
}

func compileFloat64(structName, fieldName string) (Decoder, error) {
	return newFloatDecoder(structName, fieldName), nil
}

func compileString(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newStringDecoder(structName, fieldName), nil
}

func compileBool(structName, fieldName string) (Decoder, error) {
	return newBoolDecoder(structName, fieldName), nil
}

func compileBytes(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newBytesDecoder(typ, structName, fieldName), nil
}

func compileSlice(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	elem := typ.Elem()
	decoder, err := compile(elem, structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newSliceDecoder(decoder, elem, elem.Size(), structName, fieldName), nil
}

func compileArray(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	elem := typ.Elem()
	decoder, err := compile(elem, structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newArrayDecoder(decoder, elem, typ.Len(), structName, fieldName), nil
}

func compileMap(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	keyDec, err := compileMapKey(typ.Key(), structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	valueDec, err := compile(typ.Elem(), structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newMapDecoder(typ, typ.Key(), keyDec, typ.Elem(), valueDec, structName, fieldName), nil
}

func compileInterface(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newInterfaceDecoder(typ, structName, fieldName), nil
}

func compileStruct(rt reflect.Type, structName, fieldName string, structTypeToDecoder map[reflect.Type]Decoder) (Decoder, error) {
	if dec, exists := structTypeToDecoder[rt]; exists {
		return dec, nil
	}
	structDec := newStructDecoder(structName, fieldName, map[string]*structFieldSet{})
	structTypeToDecoder[rt] = structDec
	structName = rt.Name()

	var allFields []*structFieldSet

	fieldNum := rt.NumField()
	for i := 0; i < fieldNum; i++ {
		field := rt.Field(i)
		if runtime.IsIgnoredStructField(field) {
			continue
		}

		if field.Anonymous {
			if (field.Type.Kind() == reflect.Struct) || (field.Type.Kind() == reflect.Ptr && (field.Type.Elem().Kind() == reflect.Struct)) {
				return nil, errors.New("anonymous struct field is not supported")
			}
		}

		tag := runtime.StructTagFromField(field)
		dec, err := compile(field.Type, structName, field.Name, structTypeToDecoder)
		if err != nil {
			return nil, err
		}

		if tag.IsString && isStringTagSupportedType(field.Type) {
			dec, err = newWrappedStringDecoder(field.Type, dec, structName, field.Name)
			if err != nil {
				return nil, err
			}
		}

		var key string
		if tag.Key != "" {
			key = tag.Key
		} else {
			key = field.Name
		}

		fieldSet := &structFieldSet{
			dec:      dec,
			fieldIdx: i,
			key:      key,
		}

		allFields = append(allFields, fieldSet)
	}

	seen := map[string]bool{}
	for _, set := range allFields {
		if seen[set.key] {
			return nil, fmt.Errorf("found duplicate keys for struct %s: %s", rt.String(), set.key)
		}

		seen[set.key] = true
		structDec.fieldMap[set.key] = set
	}

	delete(structTypeToDecoder, rt)

	return structDec, nil
}
