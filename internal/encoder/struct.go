package encoder

import (
	"reflect"
	"strings"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structFieldEncoder func(ctx *Ctx, sc *structCtx, b []byte, p uintptr) ([]byte, error)

func compileStruct(rt *runtime.Type) (encoder, error) {
	var fieldConfigs = make([]fieldConfig, 0, rt.NumField())
	var hasOmitEmptyField = false

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := getFieldConfig(field)
		fieldConfigs = append(fieldConfigs, cfg)
		if cfg.OmitEmpty {
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
		cfg := fieldConfigs[i]
		if cfg.Ignore {
			continue
		}

		filedNameEncoder := compileConstStringNoError(cfg.Name)

		var fieldEncoder, wrappedEncoder encoder
		var err error
		if cfg.AsString {
			fieldEncoder, err = compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			offset := field.Offset
			wrappedEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
				b = filedNameEncoder(ctx, b)
				return fieldEncoder(ctx, b, p+offset)
			}
		} else {
			fieldEncoder, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			offset := field.Offset
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

		if cfg.OmitEmpty {
			hasOmitEmptyField = true
			isEmpty, err := compileEmptyer(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
			encoders = append(encoders, fieldEncoderWithEmpty(wrappedEncoder, isEmpty))
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
func compileStructNoOmitEmpty(rt *runtime.Type, fieldConfigs []fieldConfig) (encoder, error) {
	indirect := runtime.IfaceIndir(rt)
	var encoders []encoder
	var fieldCount int64
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := fieldConfigs[i]
		if cfg.Ignore {
			continue
		}

		fieldCount++

		filedNameEncoder := compileConstStringNoError(cfg.Name)

		if cfg.AsString {
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
		sc := newStructCtx()
		defer freeStructCtx(sc)

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

func fieldEncoderWithEmpty(enc encoder, empty isEmpty) structFieldEncoder {
	return func(ctx *Ctx, sc *structCtx, structBuffer []byte, p uintptr) ([]byte, error) {
		shouldIgnore, err := empty(ctx, p)
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

type fieldConfig struct {
	Name      string
	Ignore    bool
	AsString  bool
	OmitEmpty bool
}

func getFieldConfig(field reflect.StructField) fieldConfig {
	tag := field.Tag.Get(DefaultStructTag)

	if tag == "" {
		return fieldConfig{Name: field.Name}
	}

	if tag == "-" {
		return fieldConfig{Ignore: true}
	}

	cfg := fieldConfig{Name: field.Name}
	s := strings.Split(tag, ",")

	if s[0] != "" {
		cfg.Name = s[0]
	}

	if len(s) == 1 {
		return cfg
	}

	if contains(s[1:], "string") {
		cfg.AsString = true
	}

	if contains(s[1:], "omitempty") {
		cfg.OmitEmpty = true
	}

	return cfg
}

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
