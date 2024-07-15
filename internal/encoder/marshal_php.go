package encoder

import (
	"reflect"
)

var marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()

type Marshaler interface {
	MarshalPHP() ([]byte, error)
}

func compileMarshaler(rt reflect.Type) (encoder, error) {
	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		raw, err := rv.Interface().(Marshaler).MarshalPHP()
		if err != nil {
			return nil, err
		}

		return append(b, raw...), nil
	}, nil
}
