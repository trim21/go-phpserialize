package encoder

import (
	"fmt"
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structEncoder struct {
	index int
	// a direct value handler, like `encodeInt`
	// struct encoder should de-ref pointers and pass real address to encoder.
	// address of map, slice, array may still be 0, bug theirs encoder will handle that at null.
	encode    encoder
	fieldName string // field fieldName
	omitEmpty bool
	ptr       bool
}

type compileSeenMap = map[reflect.Type]*structRecEncoder

type structRecEncoder struct {
	enc encoder
}

func (s *structRecEncoder) Encode(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
	return s.enc(ctx, b, rv)
}

func compileStruct(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	recursiveEnc, hasSeen := seen[rt]

	if hasSeen {
		return recursiveEnc.Encode, nil
	} else {
		seen[rt] = &structRecEncoder{}
	}

	enc, err := compileStructFields(rt, seen)
	if err != nil {
		return nil, err
	}

	recursiveEnc, recursiveStruct := seen[rt]
	if recursiveStruct {
		if recursiveEnc.enc == nil {
			recursiveEnc.enc = enc
			return recursiveEnc.Encode, nil
		}
	}

	return enc, nil
}

// struct don't have `omitempty` tag, fast path
func compileStructFields(rt reflect.Type, seen compileSeenMap) (encoder, error) {
	fields, compileErr := compileStructFieldsEncoders(rt, seen)
	if compileErr != nil {
		return nil, compileErr
	}

	return func(ctx *Ctx, b []byte, rv reflect.Value) ([]byte, error) {
		// shadow compiler's error
		var err error

		buf := newBuffer()
		defer freeBuffer(buf)
		structBuffer := buf.b

		var writtenField int64

	FIELD:
		for _, field := range fields {
			v := rv.Field(field.index)

			if field.omitEmpty {
				if v.IsZero() {
					continue
				}
			}

			writtenField++

			structBuffer = appendPhpStringVariable(ctx, structBuffer, field.fieldName)

			if field.ptr {
				if v.IsNil() {
					structBuffer = appendNull(structBuffer)
					continue FIELD
				}

				v = v.Elem()
			}

			structBuffer, err = field.encode(ctx, structBuffer, v)
			if err != nil {
				return b, err
			}
		}

		b = appendArrayBegin(b, writtenField)
		b = append(b, structBuffer...)
		buf.b = structBuffer

		return append(b, '}'), nil
	}, nil
}

func compileStructFieldsEncoders(rt reflect.Type, seen compileSeenMap) ([]structEncoder, error) {
	var encoders []structEncoder

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if cfg.Key == "-" || !cfg.Field.IsExported() {
			continue
		}

		var fieldEncoder encoder
		var err error

		var isPtrField = field.Type.Kind() == reflect.Ptr

		if field.Type.Kind() == reflect.Ptr {
			if field.Type.Elem().Kind() == reflect.Ptr {
				return nil, fmt.Errorf("encoding nested ptr is not supported %s", field.Type.String())
			}
		}

		if field.Anonymous {
			if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Kind() == reflect.Struct) {
				return nil, fmt.Errorf("supported for Anonymous struct field has been removed: %s", field.Type.String())
			}
		}

		if cfg.IsString {
			if field.Type.Kind() == reflect.Ptr {
				fieldEncoder, err = compileAsString(field.Type.Elem())
			} else {
				fieldEncoder, err = compileAsString(field.Type)
			}
		} else {
			if field.Type.Kind() == reflect.Ptr {
				fieldEncoder, err = compile(field.Type.Elem(), seen)
			} else {
				fieldEncoder, err = compile(field.Type, seen)
			}
		}

		if err != nil {
			return nil, err
		}

		encoders = append(encoders, structEncoder{
			index:     i,
			encode:    fieldEncoder,
			fieldName: cfg.Name(),
			omitEmpty: cfg.IsOmitEmpty,
			ptr:       isPtrField,
		})
	}

	return encoders, nil
}
