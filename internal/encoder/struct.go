package encoder

import (
	"fmt"
	"reflect"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

type structEncoder struct {
	offset uintptr
	// a direct value handler, like `encodeInt`
	// struct encoder should de-ref pointers and pass real address to encoder.
	// address of map, slice, array may still be 0, bug theirs encoder will handle that at null.
	encode    encoder
	fieldName string // field fieldName
	zero      emptyFunc
	indirect  bool
	ptr       bool
	ptrDepth  int
}

type seenMap = map[reflect.Type]*structRecEncoder

type structRecEncoder struct {
	enc encoder
}

func (s *structRecEncoder) Encode(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	return s.enc(ctx, b, p)
}

func compileStruct(rt reflect.Type, seen seenMap) (encoder, error) {
	recursiveEnc, hasSeen := seen[rt]

	if hasSeen {
		return recursiveEnc.Encode, nil
	} else {
		seen[rt] = &structRecEncoder{}
	}

	hasOmitEmpty, err := hasOmitEmptyField(rt)
	if err != nil {
		return nil, err
	}

	var enc encoder
	if !hasOmitEmpty {
		enc, err = compileStructNoOmitEmptyFastPath(rt, seen)
	} else {
		enc, err = compileStructBufferSlowPath(rt, seen)
	}
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

func hasOmitEmptyField(rt reflect.Type) (bool, error) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if field.Type.Kind() == reflect.Struct {
			if cfg.IsOmitEmpty {
				return false, fmt.Errorf("can't use 'omitempty' config with struct field: %s{}.%s", rt.String(), field.Name)
			}

			v, err := hasOmitEmptyField(field.Type)
			if err != nil {
				return false, err
			}
			if v {
				return true, nil
			}
		}
		if cfg.IsOmitEmpty {
			return true, nil
		}
	}

	return false, nil
}

// struct don't have `omitempty` tag, fast path
func compileStructNoOmitEmptyFastPath(rt reflect.Type, seen seenMap) (encoder, error) {
	fields, err := compileStructFieldsEncoders(rt, 0, seen)
	if err != nil {
		return nil, err
	}

	var fieldCount int64 = int64(len(fields))
	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		b = appendArrayBegin(b, fieldCount)

		var err error

	FIELD:
		for _, field := range fields {
			b = appendPhpStringVariable(ctx, b, field.fieldName)

			fp := field.offset + p

			if field.ptr {
				if field.indirect {
					fp = PtrDeRef(fp)
				}

				if fp == 0 {
					b = appendNull(b)
					continue
				}

				for i := 0; i < field.ptrDepth; i++ {
					fp = PtrDeRef(fp)
					if fp == 0 {
						b = appendNull(b)
						continue FIELD
					}
				}
			}

			b, err = field.encode(ctx, b, fp)
			if err != nil {
				return b, err
			}
		}

		return append(b, '}'), nil
	}, nil
}

func compileStructBufferSlowPath(rt reflect.Type, seen seenMap) (encoder, error) {
	encoders, err := compileStructFieldsEncoders(rt, 0, seen)
	if err != nil {
		return nil, err
	}

	return func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
		buf := newBuffer()
		defer freeBuffer(buf)
		structBuffer := buf.b

		var err error
		var writtenField int64

	FIELD:
		for _, field := range encoders {
			fp := field.offset + p

			if field.ptr {
				if field.indirect {
					fp = PtrDeRef(fp)
				}

				if fp == 0 {
					if field.zero != nil {
						continue FIELD
					}

					structBuffer = appendPhpStringVariable(ctx, structBuffer, field.fieldName)
					writtenField++
					structBuffer = appendNull(structBuffer)
					continue
				}

				for i := 0; i < field.ptrDepth; i++ {
					fp = PtrDeRef(fp)
					if fp == 0 {
						if field.zero != nil {
							continue FIELD
						}

						structBuffer = appendPhpStringVariable(ctx, structBuffer, field.fieldName)
						structBuffer = appendNull(structBuffer)
						writtenField++
						continue FIELD
					}
				}
			}

			if field.zero != nil {
				empty, err := field.zero(ctx, fp)
				if err != nil {
					return nil, err
				}
				if empty {
					continue
				}
			}

			structBuffer = appendPhpStringVariable(ctx, structBuffer, field.fieldName)
			structBuffer, err = field.encode(ctx, structBuffer, fp)
			if err != nil {
				return b, err
			}

			writtenField++
		}

		b = appendArrayBegin(b, writtenField)
		b = append(b, structBuffer...)
		buf.b = structBuffer

		return append(b, '}'), nil
	}, nil
}

func compileStructFieldsEncoders(rt reflect.Type, baseOffset uintptr, seen seenMap) (encoders []structEncoder, err error) {
	indirect := runtime.IfaceIndir(rt)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		cfg := runtime.StructTagFromField(field)
		if cfg.Key == "-" || !cfg.Field.IsExported() {
			continue
		}
		offset := field.Offset + baseOffset

		var ptrDepth = 0

		var isEmpty emptyFunc
		var fieldEncoder encoder
		var err error

		var isPtrField = field.Type.Kind() == reflect.Ptr

		if field.Type.Kind() == reflect.Ptr {
			isEmpty = EmptyPtr

			switch field.Type.Elem().Kind() {
			case reflect.Ptr:
				return nil, fmt.Errorf("encoding nested ptr is not supported %s", field.Type.String())

			case reflect.Map:
				ptrDepth++
				fallthrough
			default:
				fieldEncoder, err = compile(field.Type.Elem(), seen)
				if err != nil {
					return nil, err
				}
			}
		}

		if fieldEncoder == nil {
			if field.Type.Kind() == reflect.Struct && field.Anonymous {
				enc, err := compileStructFieldsEncoders(field.Type, offset, seen)
				if err != nil {
					return nil, err
				}

				encoders = append(encoders, enc...)
				continue
			}

			fieldEncoder, err = compile(field.Type, seen)
			if err != nil {
				return nil, err
			}
		}

		var enc encoder
		if cfg.IsString {
			if field.Type.Kind() == reflect.Ptr {
				enc, err = compileAsString(field.Type.Elem())
				fieldEncoder = func(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
					// fmt.Println(p)
					// fmt.Println(**(**bool)(unsafe.Pointer(&p)))
					// fmt.Println(PtrDeRef(p))
					// fmt.Println(PtrDeRef(PtrDeRef(p)))
					return enc(ctx, b, p)
				}
				// if !indirect && field.abiType.Elem().Kind() != reflect.Bool {
				// 	ptrDepth++
				// }
			} else {
				fieldEncoder, err = compileAsString(field.Type)
			}
			if err != nil {
				return nil, err
			}
		} else {
			if indirect && (field.Type.Kind() == reflect.Map) {
				isPtrField = true
			}
		}

		if cfg.IsOmitEmpty && isEmpty == nil {
			isEmpty, err = compileEmptyFunc(field.Type)
			if err != nil {
				return nil, err
			}
		}

		encoders = append(encoders, structEncoder{
			offset:    offset,
			encode:    fieldEncoder,
			fieldName: cfg.Name(),
			zero:      isEmpty,
			indirect:  indirect,
			ptrDepth:  ptrDepth,
			ptr:       isPtrField,
		})
	}

	return encoders, nil
}
