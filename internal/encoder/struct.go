package encoder

import (
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structFieldEncoder func(ctx *Ctx, sc *structCtx, b []byte, p uintptr) ([]byte, error)

func compileStruct(rt *runtime.Type) (encoder, error) {
	var fieldConfigs = make([]*runtime.StructTag, 0, rt.NumField())
	var hasOmitEmptyField = false

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		fieldConfigs = append(fieldConfigs, cfg)
		if cfg.IsOmitEmpty {
			hasOmitEmptyField = true
		}
	}

	if !hasOmitEmptyField {
		return compileStructNoOmitEmpty(rt, fieldConfigs)
	}

	indirect := runtime.IfaceIndir(rt)
	var encoders []structFieldEncoder
	for i := 0; i < rt.NumField(); i++ {
		cfg := fieldConfigs[i]
		if cfg.Key == "-" || !cfg.Field.IsExported() {
			continue
		}

		field := cfg.Field
		offset := field.Offset

		filedNameEncoder := compileConstStringNoError(cfg.Key)

		var isEmpty emptyFunc
		var fieldEncoder encoder
		var err error

		if !indirect {
			if field.Type.Kind() == reflect.Ptr {
				switch field.Type.Elem().Kind() {
				case reflect.Array:
					isEmpty = EmptyPtr
					enc, err := compileArray(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}
					// fieldEncoder = deRefEncoder(enc)
					fieldEncoder = enc
				case reflect.String:
					isEmpty = EmptyPtr
					fieldEncoder = encodeString
				case reflect.Slice:
					isEmpty = EmptyPtr
					enc, err := compileSlice(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}
					fieldEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
						return enc(ctx, b, p)
					}
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
			fieldEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
		}

		var wrappedEncoder encoder
		if cfg.IsString {
			fieldEncoder, err = compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldEncoder(ctx, b, p+offset)
			}
		} else {
			if indirect && (field.Type.Kind() == reflect.Map) {
				wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
					b = filedNameEncoder(ctx, b)
					return fieldEncoder(ctx, b, PtrDeRef(p+offset))
				}
			} else {
				wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
					b = filedNameEncoder(ctx, b)
					return fieldEncoder(ctx, b, p+offset)
				}
			}
		}

		if cfg.IsOmitEmpty {
			if isEmpty == nil {
				isEmpty, err = compileEmptyFunc(runtime.Type2RType(field.Type))
				if err != nil {
					return nil, err
				}
			}
			encoders = append(encoders, fieldEncoderWithEmpty(wrappedEncoder, offset, isEmpty))
		} else {
			encoders = append(encoders, func(ctx *Ctx, sc *structCtx, structBuffer []byte, p uintptr) ([]byte, error) {
				sc.writtenField++
				return wrappedEncoder(ctx, structBuffer, p)
			})
		}
	}

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		sc := newStructCtx()
		defer freeStructCtx(sc)

		buf := newBuffer()
		defer freeBuffer(buf)

		var err error
		for _, enc := range encoders {
			buf.b, err = enc(ctx, sc, buf.b, p)
			if err != nil {
				return b, err
			}
		}

		b = appendArrayBeginBytes(b, sc.writtenField)
		b = append(b, buf.b...)

		return append(b, '}'), nil
	}, nil
}

// struct don't have `omitempty` tag, fast path
func compileStructNoOmitEmpty(rt *runtime.Type, fieldConfigs []*runtime.StructTag) (encoder, error) {
	indirect := runtime.IfaceIndir(rt)
	var encoders []encoder
	var fieldCount int64
	for i := 0; i < rt.NumField(); i++ {
		cfg := fieldConfigs[i]
		if cfg.Key == "-" || !cfg.Field.IsExported() {
			continue
		}
		field := cfg.Field

		fieldCount++

		filedNameEncoder := compileConstStringNoError(cfg.Key)

		if cfg.IsString {
			fieldValueEncoder, err := compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			offset := field.Offset
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldValueEncoder(ctx, b, p+offset)
			})
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
						return nil, err
					}
				case reflect.Slice:
					fieldValueEncoder, err = compileSlice(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}
				case reflect.Map, reflect.String:
					fieldValueEncoder, err = compile(runtime.Type2RType(field.Type.Elem()))
					if err != nil {
						return nil, err
					}

					// fieldValueEncoder = deRefEncoder(fieldValueEncoder)
				}
			}
		}

		if fieldValueEncoder == nil {
			fieldValueEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
		}

		offset := field.Offset

		if indirect && field.Type.Kind() == reflect.Map {
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldValueEncoder(ctx, b, PtrDeRef(p+offset))
			})
		} else {
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldValueEncoder(ctx, b, p+offset)
			})
		}
	}

	return structEncoderNoOmitEmpty(encoders, fieldCount), nil
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

func fieldEncoderWithEmpty(enc encoder, offset uintptr, empty emptyFunc) structFieldEncoder {
	return func(ctx *Ctx, sc *structCtx, structBuffer []byte, p uintptr) ([]byte, error) {
		shouldIgnore, err := empty(ctx, p+offset)
		if err != nil {
			return structBuffer, err
		}
		if shouldIgnore {
			return structBuffer, nil
		}

		sc.writtenField++

		return enc(ctx, structBuffer, p)
	}
}
