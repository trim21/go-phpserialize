package encoder

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/goccy/go-reflect"
)

const DefaultStructTag = "php"

var (
	typeToEncoderMap sync.Map
	bufpool          = sync.Pool{
		New: func() interface{} {
			return &buffer{
				b: make([]byte, 0, 1024),
			}
		},
	}
)

type buffer struct {
	b []byte
}

type encoder func(buf *buffer, p uintptr) error

func Marshal(v interface{}) ([]byte, error) {
	// Technique 1.
	// Get type information and pointer from interface{} rValue without allocation.
	typ, ptr := reflect.TypeAndPtrOf(v)
	typeID := reflect.TypeID(typ)
	p := uintptr(ptr)

	// Technique 2.
	// Reuse the buffer once allocated using sync.Pool
	buf := bufpool.Get().(*buffer)
	buf.b = buf.b[:0]
	defer bufpool.Put(buf)

	// Technique 3.
	// builds an optimized path by typeID and caches it
	if enc, ok := typeToEncoderMap.Load(typeID); ok {
		if err := enc.(encoder)(buf, p); err != nil {
			return nil, err
		}

		// allocate a new buffer required length only
		b := make([]byte, len(buf.b))
		copy(b, buf.b)
		return b, nil
	}

	// First time,
	// builds an optimized path by type and caches it with typeID.
	rv := reflect.ValueOf(v)
	enc, err := compile(typ, rv)
	if err != nil {
		return nil, err
	}
	typeToEncoderMap.Store(typeID, enc)
	if err := enc(buf, p); err != nil {
		return nil, err
	}

	runtime.KeepAlive(v) // didn't keep ref, so just hold the variable

	// allocate a new buffer required length only
	b := make([]byte, len(buf.b))

	copy(b, buf.b)
	return b, nil
}

func compile(typ reflect.Type, rv reflect.Value) (encoder, error) {
	switch typ.Kind() {
	case reflect.Bool:
		return compileBool(typ)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return compileInt(typ)
	case reflect.String:
		return encodeStringVariable, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return compileUint(typ)
	case reflect.Float32, reflect.Float64:
		return compileFloat(typ)
	case reflect.Struct:
		return compileStruct(typ, rv)
	case reflect.Slice:
		return compileSlice(typ, rv)
	case reflect.Map:
		return compileMap(typ, rv)
	case reflect.Interface:
		return compileInterface(typ)
	}

	return nil, fmt.Errorf("failed to build encoder, unsupported type %s (kind %s)", typ.String(), typ.Kind())
}
