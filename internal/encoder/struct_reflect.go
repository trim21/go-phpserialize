package encoder

import "github.com/goccy/go-reflect"

func reflectStruct(ctx *Ctx, rv reflect.Value) error {
	rt := rv.Type()
	appendArrayBegin(ctx, int64(rv.NumField()))

	for i := 0; i < rv.NumField(); i++ {
		appendString(ctx, getFieldName(rt.Field(i)))

		err := reflectInterfaceValue(ctx, rv.Field(i))
		if err != nil {
			return err
		}
		continue
	}

	ctx.b = append(ctx.b, '}')

	return nil
}
