package encoder

import (
	"strings"

	"github.com/goccy/go-reflect"
)

func compileStruct(rt reflect.Type, rv reflect.Value) (encoder, error) {
	var encoders []encoder

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		enc, err := compileFieldName(field)
		if err != nil {
			return nil, err
		}
		encoders = append(encoders, enc)

		enc, err = compile(field.Type, rv.Field(i))
		if err != nil {
			return nil, err
		}

		offset := field.Offset
		encoders = append(encoders, func(buf *Ctx, p uintptr) error {
			return enc(buf, p+offset)
		})
	}

	fields := int64(rt.NumField())
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

func compileFieldName(field reflect.StructField) (encoder, error) {
	var name = field.Name

	tag := field.Tag.Get(DefaultStructTag)

	if tag != "" {
		name = strings.Split(tag, ",")[0]
	}

	return compileConstString(name)
}

func getFieldName(field reflect.StructField) string {
	var name = field.Name

	tag := field.Tag.Get(DefaultStructTag)

	if tag != "" {
		name = strings.Split(tag, ",")[0]
	}

	return name
}
