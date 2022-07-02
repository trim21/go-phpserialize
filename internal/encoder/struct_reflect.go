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
	var fields int64
	for i := 0; i < rv.NumField(); i++ {
		if rt.Field(i).Tag.Get(DefaultStructTag) != "-" {
			fields++
		}
	}

	appendArrayBegin(ctx, fields)

	for i := 0; i < rv.NumField(); i++ {
		name := getFieldName(rt.Field(i))
		if name == "" {
			continue
		}

		appendString(ctx, name)

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
