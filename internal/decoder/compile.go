package decoder

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"unicode"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

var (
	cachedDecoderMap atomic.Pointer[map[reflect.Type]Decoder]
)

func init() {
	var m = map[reflect.Type]Decoder{}
	cachedDecoderMap.Store(&m)
}

func CompileToGetDecoder(typ reflect.Type) (Decoder, error) {
	return compileToGetDecoderSlowPath(typ)
}

func storeDecoder(rt reflect.Type, dec Decoder, m map[reflect.Type]Decoder) {
	newDecoderMap := make(map[reflect.Type]Decoder, len(m)+1)
	newDecoderMap[rt] = dec

	for k, v := range m {
		newDecoderMap[k] = v
	}

	cachedDecoderMap.Store(&newDecoderMap)
}

func compileToGetDecoderSlowPath(rt reflect.Type) (Decoder, error) {
	decoderMap := *cachedDecoderMap.Load()
	if dec, exists := decoderMap[rt]; exists {
		return dec, nil
	}

	dec, err := compileHead(rt, map[uintptr]Decoder{})
	if err != nil {
		return nil, err
	}
	storeDecoder(rt, dec, decoderMap)
	return dec, nil
}

func compileHead(rt reflect.Type, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
	switch {
	case reflect.PointerTo(rt).Implements(unmarshalPHPType):
		return newUnmarshalTextDecoder(reflect.PointerTo(rt), "", ""), nil
	}
	return compile(rt.Elem(), "", "", structTypeToDecoder)
}

func compile(rt reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
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

func compileMapKey(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
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

func compilePtr(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
	dec, err := compile(typ.Elem(), structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newPtrDecoder(dec, typ.Elem(), structName, fieldName)
}

func compileInt(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v int64) {
		*(*int)(p) = int(v)
	}), nil
}

func compileInt8(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v int64) {
		*(*int8)(p) = int8(v)
	}), nil
}

func compileInt16(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v int64) {
		*(*int16)(p) = int16(v)
	}), nil
}

func compileInt32(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v int64) {
		*(*int32)(p) = int32(v)
	}), nil
}

func compileInt64(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newIntDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v int64) {
		*(*int64)(p) = v
	}), nil
}

func compileUint(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v uint64) {
		*(*uint)(p) = uint(v)
	}), nil
}

func compileUint8(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v uint64) {
		*(*uint8)(p) = uint8(v)
	}), nil
}

func compileUint16(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v uint64) {
		*(*uint16)(p) = uint16(v)
	}), nil
}

func compileUint32(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v uint64) {
		*(*uint32)(p) = uint32(v)
	}), nil
}

func compileUint64(typ reflect.Type, structName, fieldName string) (Decoder, error) {
	return newUintDecoder(typ, structName, fieldName, func(p unsafe.Pointer, v uint64) {
		*(*uint64)(p) = v
	}), nil
}

func compileFloat32(structName, fieldName string) (Decoder, error) {
	return newFloatDecoder(structName, fieldName, func(p unsafe.Pointer, v float64) {
		*(*float32)(p) = float32(v)
	}), nil
}

func compileFloat64(structName, fieldName string) (Decoder, error) {
	return newFloatDecoder(structName, fieldName, func(p unsafe.Pointer, v float64) {
		*(*float64)(p) = v
	}), nil
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

func compileSlice(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
	elem := typ.Elem()
	decoder, err := compile(elem, structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newSliceDecoder(decoder, elem, elem.Size(), structName, fieldName), nil
}

func compileArray(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
	elem := typ.Elem()
	decoder, err := compile(elem, structName, fieldName, structTypeToDecoder)
	if err != nil {
		return nil, err
	}
	return newArrayDecoder(decoder, elem, typ.Len(), structName, fieldName), nil
}

func compileMap(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
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

func typeToStructTags(typ reflect.Type) runtime.StructTags {
	tags := runtime.StructTags{}
	fieldNum := typ.NumField()
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		if runtime.IsIgnoredStructField(field) {
			continue
		}
		tags = append(tags, runtime.StructTagFromField(field))
	}
	return tags
}

func compileStruct(typ reflect.Type, structName, fieldName string, structTypeToDecoder map[uintptr]Decoder) (Decoder, error) {
	fieldNum := typ.NumField()
	fieldMap := map[string]*structFieldSet{}
	typeptr := runtime.TypeID(typ)
	if dec, exists := structTypeToDecoder[typeptr]; exists {
		return dec, nil
	}
	structDec := newStructDecoder(structName, fieldName, fieldMap)
	structTypeToDecoder[typeptr] = structDec
	structName = typ.Name()
	tags := typeToStructTags(typ)
	allFields := []*structFieldSet{}
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		if runtime.IsIgnoredStructField(field) {
			continue
		}
		isUnexportedField := unicode.IsLower([]rune(field.Name)[0])
		tag := runtime.StructTagFromField(field)
		dec, err := compile(runtime.Type2RType(field.Type), structName, field.Name, structTypeToDecoder)
		if err != nil {
			return nil, err
		}
		if field.Anonymous && !tag.IsTaggedKey {
			if stDec, ok := dec.(*structDecoder); ok {
				if runtime.Type2RType(field.Type) == typ {
					// recursive definition
					continue
				}
				for k, v := range stDec.fieldMap {
					if tags.ExistsKey(k) {
						continue
					}
					fieldSet := &structFieldSet{
						dec:         v.dec,
						offset:      field.Offset + v.offset,
						isTaggedKey: v.isTaggedKey,
						key:         k,
						keyLen:      int64(len(k)),
					}
					allFields = append(allFields, fieldSet)
				}
			} else if pdec, ok := dec.(*ptrDecoder); ok {
				contentDec := pdec.contentDecoder()
				if pdec.typ == typ {
					// recursive definition
					continue
				}
				var fieldSetErr error
				if isUnexportedField {
					fieldSetErr = fmt.Errorf(
						"php: cannot set embedded pointer to unexported struct: %v",
						field.Type.Elem(),
					)
				}
				if dec, ok := contentDec.(*structDecoder); ok {
					for k, v := range dec.fieldMap {
						if tags.ExistsKey(k) {
							continue
						}
						fieldSet := &structFieldSet{
							dec:         newAnonymousFieldDecoder(pdec.typ, v.offset, v.dec),
							offset:      field.Offset,
							isTaggedKey: v.isTaggedKey,
							key:         k,
							keyLen:      int64(len(k)),
							err:         fieldSetErr,
						}
						allFields = append(allFields, fieldSet)
					}
				}
			} else {
				fieldSet := &structFieldSet{
					dec:         dec,
					offset:      field.Offset,
					isTaggedKey: tag.IsTaggedKey,
					key:         field.Name,
					keyLen:      int64(len(field.Name)),
				}
				allFields = append(allFields, fieldSet)
			}
		} else {
			if tag.IsString && isStringTagSupportedType(runtime.Type2RType(field.Type)) {
				dec, err = newWrappedStringDecoder(runtime.Type2RType(field.Type), dec, structName, field.Name)
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
				dec:         dec,
				offset:      field.Offset,
				isTaggedKey: tag.IsTaggedKey,
				key:         key,
				keyLen:      int64(len(key)),
			}
			allFields = append(allFields, fieldSet)
		}
	}
	for _, set := range filterDuplicatedFields(allFields) {
		fieldMap[set.key] = set
		lower := strings.ToLower(set.key)
		if _, exists := fieldMap[lower]; !exists {
			// first win
			fieldMap[lower] = set
		}
	}
	delete(structTypeToDecoder, typeptr)
	structDec.tryOptimize()
	return structDec, nil
}

func filterDuplicatedFields(allFields []*structFieldSet) []*structFieldSet {
	fieldMap := map[string][]*structFieldSet{}
	for _, field := range allFields {
		fieldMap[field.key] = append(fieldMap[field.key], field)
	}
	duplicatedFieldMap := map[string]struct{}{}
	for k, sets := range fieldMap {
		sets = filterFieldSets(sets)
		if len(sets) != 1 {
			duplicatedFieldMap[k] = struct{}{}
		}
	}

	filtered := make([]*structFieldSet, 0, len(allFields))
	for _, field := range allFields {
		if _, exists := duplicatedFieldMap[field.key]; exists {
			continue
		}
		filtered = append(filtered, field)
	}
	return filtered
}

func filterFieldSets(sets []*structFieldSet) []*structFieldSet {
	if len(sets) == 1 {
		return sets
	}
	filtered := make([]*structFieldSet, 0, len(sets))
	for _, set := range sets {
		if set.isTaggedKey {
			filtered = append(filtered, set)
		}
	}
	return filtered
}
