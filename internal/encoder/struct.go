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
	encoders, err := compileStructFieldsEncoders(rt, 0)
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

		b = appendArrayBegin(b, writtenField)
		b = append(b, structBuffer...)
		buf.b = structBuffer

		return append(b, '}'), nil
	}, nil
}

func compileStructFieldsEncoders(rt *runtime.Type, baseOffset uintptr) (encoders []structEncoder, err error) {
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

		if field.Type.Kind() == reflect.Ptr {
			switch field.Type.Elem().Kind() {
			case reflect.Map:
				isEmpty = func(ctx *Ctx, p uintptr) (isEmpty bool, err error) {
					if p == 0 {
						return true, nil
					}
					p = PtrDeRef(p)
					return p == 0, nil
				}
				enc, err := compilePtr(runtime.Type2RType(field.Type), indirect)
				if err != nil {
					return nil, err
				}
				fieldEncoder = deRefNilEncoder(enc)
			case reflect.String, reflect.Array, reflect.Slice:
				if indirect {
					isEmpty = func(ctx *Ctx, p uintptr) (isEmpty bool, err error) {
						if p == 0 {
							return true, nil
						}
						p = PtrDeRef(p)
						return p == 0, nil
					}
				} else {
					isEmpty = EmptyPtr
				}

				enc, err := compilePtr(runtime.Type2RType(field.Type), indirect)
				if err != nil {
					return nil, err
				}
				fieldEncoder = enc
			case reflect.Struct:
				enc, err := compilePtr(runtime.Type2RType(field.Type), false)
				if err != nil {
					return nil, err
				}

				if indirect {
					fieldEncoder = onlyDeReferEncoder(enc)
				} else {
					fieldEncoder = enc
				}
			}
		}

		if fieldEncoder == nil {
			if field.Type.Kind() == reflect.Struct && field.Anonymous {
				enc, err := compileStructFieldsEncoders(runtime.Type2RType(field.Type), offset)
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
				fieldEncoder = deRefNilEncoder(fieldEncoder)
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

		encoders = append(encoders, structEncoder{
			offset:    offset,
			encode:    fieldEncoder,
			fieldName: cfg.Name(),
			zero:      isEmpty,
		})
	}

	return encoders, nil
}

// struct don't have `omitempty` tag, fast path
func compileStructNoOmitEmptyFastPath(rt *runtime.Type) (encoder, error) {
	fields, err := compileStructFieldsEncoders(rt, 0)
	if err != nil {
		return nil, err
	}

	var fieldCount int64 = int64(len(fields))
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		b = appendArrayBegin(b, fieldCount)

		var err error
		for _, field := range fields {
			b = appendPhpStringVariable(ctx, b, field.fieldName)

			b, err = field.encode(ctx, b, field.offset+p)
			if err != nil {
				return b, err
			}
		}

		return append(b, '}'), nil
	}, nil
}
