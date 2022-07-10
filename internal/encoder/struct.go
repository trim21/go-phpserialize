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
		field := rt.Field(i)
		offset := field.Offset
		cfg := fieldConfigs[i]
		if cfg.Key == "-" {
			continue
		}

		filedNameEncoder := compileConstStringNoError(cfg.Key)

		var fieldEncoder, wrappedEncoder encoder
		var err error
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
			fieldEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			if indirect && field.Type.Kind() == reflect.Map {
				wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
					b = filedNameEncoder(ctx, b)
					return fieldEncoder(ctx, b, ptrOfPtr(p+offset))
				}
			} else {
				wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
					b = filedNameEncoder(ctx, b)
					return fieldEncoder(ctx, b, p+offset)
				}
			}
		}

		if cfg.IsOmitEmpty {
			hasOmitEmptyField = true
			isEmpty, err := compileEmptyer(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
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
		field := rt.Field(i)
		cfg := fieldConfigs[i]
		if cfg.Key == "-" {
			continue
		}

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

		fieldValueEncoder, err := compile(runtime.Type2RType(field.Type))
		if err != nil {
			return nil, err
		}

		offset := field.Offset
		if indirect && field.Type.Kind() == reflect.Map {
			encoders = append(encoders, func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldValueEncoder(ctx, b, ptrOfPtr(p+offset))
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

func fieldEncoderWithEmpty(enc encoder, offset uintptr, empty isEmpty) structFieldEncoder {
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
