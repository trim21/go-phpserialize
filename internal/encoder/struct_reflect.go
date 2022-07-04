package encoder

import (
	"github.com/goccy/go-reflect"
)

func reflectStruct(ctx *Ctx, rv reflect.Value) error {
	rt := rv.Type()
	// typeID := uintptr(unsafe.Pointer(rt))
	//
	// struct {
	//
	// }{}
	//

	var fieldsCfg = make([]fieldConfig, rv.NumField())

	var fields int64
	for i := 0; i < rv.NumField(); i++ {
		cfg := getFieldConfig(rt.Field(i))
		fieldsCfg = append(fieldsCfg, cfg)
		if cfg.Ignore {
			fields++
		}
	}

	appendArrayBegin(ctx, fields)

	for i := 0; i < rv.NumField(); i++ {
		cfg := fieldsCfg[i]
		if cfg.Ignore {
			continue
		}

		appendString(ctx, cfg.Name)

		if cfg.AsString {
			err := reflectInterfaceValue(ctx, rv.Field(i))
			if err != nil {
				return err
			}

			continue
		}

		err := reflectInterfaceValue(ctx, rv.Field(i))
		if err != nil {
			return err
		}
	}

	ctx.b = append(ctx.b, '}')

	return nil
}

type structFieldConfig struct {
	fieldNames []string
}
