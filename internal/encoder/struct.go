package encoder

import (
	"strings"

	"github.com/goccy/go-reflect"
)

func compileStruct(rt reflect.Type, rv reflect.Value) (encoder, error) {
	var encoders []encoder

	var fields int64
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := getFieldConfig(field)
		if cfg.Ignore {
			continue
		}

		enc, err := compileConstString(cfg.Name)
		if err != nil {
			return nil, err
		}
		encoders = append(encoders, enc)

		if cfg.AsString {
			enc, err = compileAsString(field.Type, rv.Field(i))
			if err != nil {
				return nil, err
			}

			fields++
			offset := field.Offset
			encoders = append(encoders, func(buf *Ctx, p uintptr) error {
				return enc(buf, p+offset)
			})

			continue
		}

		enc, err = compile(field.Type, rv.Field(i))
		if err != nil {
			return nil, err
		}

		fields++
		offset := field.Offset
		encoders = append(encoders, func(buf *Ctx, p uintptr) error {
			return enc(buf, p+offset)
		})
	}

	return func(ctx *Ctx, p uintptr) error {
		appendArrayBegin(ctx, fields)

		for _, enc := range encoders {
			if err := enc(ctx, p); err != nil {
				return err
			}
		}
		ctx.b = append(ctx.b, '}')
		return nil
	}, nil
}

func getFieldName(field reflect.StructField) string {
	tag := field.Tag.Get(DefaultStructTag)

	if tag == "" {
		return field.Name
	}

	if tag == "-" {
		return ""
	}

	i := strings.IndexByte(tag, ',')
	if i <= 0 {
		return tag
	}

	if i == 0 {
		return field.Name
	}

	return tag[:i]
}

type fieldConfig struct {
	Name     string
	Ignore   bool
	AsString bool
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
