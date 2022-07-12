package encoder

import (
	"fmt"
	"reflect"
	"time"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structEncoder struct {
	offset    uintptr
	encode    encoder
	fieldName string // field fieldName
	zero      emptyFunc
}

var timeType = runtime.Type2RType(reflect.TypeOf((*time.Time)(nil)).Elem())

func compileStruct(rt *runtime.Type) (encoder, error) {
	hasOmitEmpty, err := hasOmitEmptyField(runtime.RType2Type(rt))
	if err != nil {
		return nil, err
	}

	if !hasOmitEmpty {
		return compileStructNoOmitEmptyFastPath(rt)
	}

	return compileStructBufferSlowPath(rt)
}

func hasOmitEmptyField(rt reflect.Type) (bool, error) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if field.Type.Kind() == reflect.Struct {
			if cfg.IsOmitEmpty {
				return false, fmt.Errorf("can't use 'omitempty' config with struct field: %s{}.%s", rt.String(), field.Name)
			}

			v, err := hasOmitEmptyField(field.Type)
			if err != nil {
				return false, err
			}
			if v {
				return true, nil
			}
		}
		if cfg.IsOmitEmpty {
			return true, nil
		}
	}

	return false, nil
}

func compileStructBufferSlowPath(rt *runtime.Type) (encoder, error) {
	encoders, err := compileStructFieldsEncodersOmitEmpty(rt, 0)
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		buf := newBuffer()
		defer freeBuffer(buf)
		structBuffer := buf.b

		var writtenField int64
		for _, enc := range encoders {
			if enc.zero != nil {
				empty, err := enc.zero(ctx, p+enc.offset)
				if err != nil {
					return nil, err
				}
				if empty {
					continue
				}
			}

			structBuffer = appendPhpStringVariable(ctx, structBuffer, enc.fieldName)
			structBuffer, err = enc.encode(ctx, structBuffer, p+enc.offset)
			if err != nil {
				return b, err
			}

			writtenField++
		}

		b = appendArrayBeginBytes(b, writtenField)
		b = append(b, structBuffer...)
		buf.b = structBuffer

		return append(b, '}'), nil
	}, nil
}

func compileStructFieldsEncodersOmitEmpty(rt *runtime.Type, baseOffset uintptr) (encoders []structEncoder, err error) {
	indirect := runtime.IfaceIndir(rt)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if cfg.Key == "-" || !cfg.Field.IsExported() {
			continue
		}
		offset := field.Offset + baseOffset

		var isEmpty emptyFunc
		var fieldEncoder encoder
		var err error

		if !indirect {
			if field.Type.Kind() == reflect.Ptr {
				switch field.Type.Elem().Kind() {
				case reflect.Array, reflect.Slice, reflect.String:
					isEmpty = EmptyPtr
					enc, err := compile(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}
					fieldEncoder = enc
				case reflect.Map:
					isEmpty = EmptyPtr
					enc, err := compileMap(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}
					fieldEncoder = deRefEncoder(enc)
				}
			}
		}

		if fieldEncoder == nil {
			if field.Type.Kind() == reflect.Struct && field.Anonymous {
				enc, err := compileStructFieldsEncodersOmitEmpty(runtime.Type2RType(field.Type), offset)
				if err != nil {
					return nil, err
				}

				encoders = append(encoders, enc...)
				continue
			}

			fieldEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
		}

		if cfg.IsString {
			fieldEncoder, err = compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
		} else {
			if indirect && (field.Type.Kind() == reflect.Map) {
				fieldEncoder = deRefEncoder(fieldEncoder)
			} else {

			}
		}

		if cfg.IsOmitEmpty {
			if isEmpty == nil {
				isEmpty, err = compileEmptyFunc(runtime.Type2RType(field.Type))
				if err != nil {
					return nil, err
				}
			}
		}

		var key string
		if cfg.Key != "" {
			key = cfg.Key
		} else {
			key = field.Name
		}

		encoders = append(encoders, structEncoder{
			offset:    offset,
			encode:    fieldEncoder,
			fieldName: key,
			zero:      isEmpty,
		})
	}

	return encoders, nil
}

// struct don't have `omitempty` tag, fast path
func compileStructNoOmitEmptyFastPath(rt *runtime.Type) (encoder, error) {
	var fieldConfigs = make([]*runtime.StructTag, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		fieldConfigs = append(fieldConfigs, cfg)
	}

	fieldCount, encoders, err := compileStructFieldsEncodersNoOmit(rt, 0)
	if err != nil {
		return nil, err
	}

	return structEncoderNoOmitEmpty(encoders, fieldCount), nil
}

func compileStructFieldsEncodersNoOmit(rt *runtime.Type, baseOffset uintptr) (fieldCount int64, encoders []encoder, err error) {
	indirect := runtime.IfaceIndir(rt)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if cfg.Key == "-" || !field.IsExported() {
			continue
		}

		fieldCount++

		var key string
		if cfg.Key != "" {
			key = cfg.Key
		} else {
			key = field.Name
		}

		if cfg.IsString {
			fieldValueEncoder, err := compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return 0, nil, err
			}

			offset := field.Offset
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = appendPhpStringVariable(ctx, b, key)
				return fieldValueEncoder(ctx, b, p+offset+baseOffset)
			})
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			count, enc, err := compileStructFieldsEncodersNoOmit(runtime.Type2RType(field.Type), baseOffset+field.Offset)
			if err != nil {
				return 0, nil, err
			}

			fieldCount += count - 1 // remove current field
			encoders = append(encoders, enc...)

			continue
		}

		var fieldValueEncoder encoder
		var err error

		if !indirect {
			if field.Type.Kind() == reflect.Ptr {
				switch field.Type.Elem().Kind() {
				case reflect.Array:
					fieldValueEncoder, err = compileArray(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return 0, nil, err
					}
				case reflect.Slice:
					fieldValueEncoder, err = compileSlice(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return 0, nil, err
					}
				case reflect.Map, reflect.String:
					fieldValueEncoder, err = compile(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return 0, nil, err
					}
				}
			}
		}

		if fieldValueEncoder == nil {
			fieldValueEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return 0, nil, err
			}
		}

		offset := field.Offset

		if indirect && field.Type.Kind() == reflect.Map {
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = appendPhpStringVariable(ctx, b, key)
				return fieldValueEncoder(ctx, b, PtrDeRef(p+offset+baseOffset))
			})
		} else {
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = appendPhpStringVariable(ctx, b, key)
				return fieldValueEncoder(ctx, b, p+offset+baseOffset)
			})
		}
	}

	return
}

func structEncoderNoOmitEmpty(encoders []encoder, fieldCount int64) encoder {
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		b = appendArrayBeginBytes(b, fieldCount)

		var err error
		for _, enc := range encoders {
			b, err = enc(ctx, b, p)
			if err != nil {
				return b, err
			}
		}

		return append(b, '}'), nil
	}
}
