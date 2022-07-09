package encoder

import (
	"reflect"
	"strings"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structFieldEncoder func(ctx *Ctx, sc *structCtx, p uintptr) error

func compileStruct(rt *runtime.Type) (encoder, error) {
	indirect := runtime.IfaceIndir(rt)
	var encoders []structFieldEncoder

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := getFieldConfig(field)
		if cfg.Ignore {
			continue
		}

		filedNameEncoder := compileConstStringNoError(cfg.Name)

		var enc encoder
		var err error
		var wrappedEncoder encoder
		if cfg.AsString {
			enc, err = compileAsString(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			offset := field.Offset
			wrappedEncoder = func(ctx *Ctx, p uintptr) error {
				filedNameEncoder(ctx)
				return enc(ctx, p+offset)
			}
		} else {
			enc, err = compile(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}

			offset := field.Offset
			if indirect && field.Type.Kind() == reflect.Map {
				wrappedEncoder = func(ctx *Ctx, p uintptr) error {
					filedNameEncoder(ctx)
					return enc(ctx, ptrOfPtr(p+offset))
				}
			} else {
				wrappedEncoder = func(ctx *Ctx, p uintptr) error {
					filedNameEncoder(ctx)
					return enc(ctx, p+offset)
				}
			}
		}

		var isEmpty = notIgnore

		if cfg.OmitEmpty {
			isEmpty, err = compileEmptyer(runtime.Type2RType(field.Type))
			if err != nil {
				return nil, err
			}
		}

		encoders = append(encoders, func(ctx *Ctx, sc *structCtx, p uintptr) error {
			empty, err := isEmpty(ctx, p)
			if err != nil {
				return err
			}
			if empty {
				return nil
			}
			sc.writtenField++
			return wrappedEncoder(ctx, p)
		})
	}

	return func(ctx *Ctx, p uintptr) error {
		sc := newStructCtx()
		defer freeStructCtx(sc)

		newCtx := newCtx()
		defer freeCtx(newCtx)

		for _, enc := range encoders {
			if err := enc(newCtx, sc, p); err != nil {
				return err
			}
		}

		appendArrayBegin(ctx, sc.writtenField)
		ctx.b = append(ctx.b, newCtx.b...)
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
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
