package decoder

import (
	"reflect"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
)

var (
	sliceType = reflect.TypeOf((*sliceHeader)(nil)).Elem()
	nilSlice  = unsafe.Pointer(&sliceHeader{})
)

type sliceDecoder struct {
	elemType          reflect.Type
	isElemPointerType bool
	valueDecoder      Decoder
	size              uintptr
	structName        string
	fieldName         string
}

// If use reflect.SliceHeader, data type is uintptr.
// In this case, Go compiler cannot trace reference created by newArray().
// So, define using unsafe.Pointer as data type
type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

const (
	defaultSliceCapacity = 2
)

func newSliceDecoder(dec Decoder, elemType reflect.Type, size uintptr, structName, fieldName string) *sliceDecoder {
	return &sliceDecoder{
		valueDecoder:      dec,
		elemType:          elemType,
		isElemPointerType: elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Map,
		size:              size,
		structName:        structName,
		fieldName:         fieldName,
	}
}

func (d *sliceDecoder) errNumber(offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  "number",
		Type:   reflect.SliceOf(d.elemType),
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

func (d *sliceDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, rv reflect.Value) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	switch buf[cursor] {
	case 'N':
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		cursor += 2
		rv.SetZero()
		return cursor, nil
	case 'a':
		cursor++
		if buf[cursor] != ':' {
			return cursor, errors.ErrExpected("':' before array length", cursor)
		}

		cursor++
		if buf[cursor] == '0' {
			err := validateEmptyArray(buf, cursor)
			if err != nil {
				return cursor, err
			}
			rv.SetZero()
			return cursor + 4, nil
		}

		arrLen, end, err := readLengthInt(buf, cursor-1)
		if err != nil {
			return cursor, err
		}
		cursor = end

		if buf[cursor] != '{' {
			return cursor, errors.ErrInvalidBeginningOfArray(buf[cursor], cursor)
		}
		cursor++

		// pre-alloc
		slice := &sliceHeader{
			data: newArray(d.elemType, arrLen),
			len:  0,
			cap:  arrLen,
		}
		srcLen := slice.len
		capacity := slice.cap
		data := slice.data

		idx := 0
		for {
			currentIndex, end, err := readInt(buf, cursor)
			if err != nil {
				return 0, err
			}

			idx = currentIndex
			cursor = end

			if capacity <= idx {
				src := sliceHeader{data: data, len: idx, cap: capacity}
				capacity *= 2
				data = newArray(d.elemType, capacity)
				dst := sliceHeader{data: data, len: idx, cap: capacity}
				copySlice(d.elemType, dst, src)
			}
			ep := unsafe.Pointer(uintptr(data) + uintptr(idx)*d.size)
			// if srcLen is greater than idx, keep the original reference
			if srcLen <= idx {
				if d.isElemPointerType {
					**(**unsafe.Pointer)(unsafe.Pointer(&ep)) = nil // initialize elem pointer
				} else {
					// assign new element to the slice
					typedmemmove(d.elemType, ep, unsafe_New(d.elemType))
				}
			}

			c, err := d.valueDecoder.Decode(ctx, cursor, depth, ep)
			if err != nil {
				return 0, err
			}

			cursor = c
			if buf[cursor] == '}' {
				slice.cap = capacity
				slice.len = idx + 1
				slice.data = data
				dst := (*sliceHeader)(p)
				dst.len = idx + 1
				if dst.len > dst.cap {
					dst.data = newArray(d.elemType, dst.len)
					dst.cap = dst.len
				}
				copySlice(d.elemType, *dst, *slice)
				d.releaseSlice(slice)
				cursor++
				return cursor, nil
			}
		}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return 0, d.errNumber(cursor)
	default:
		return 0, errors.ErrUnexpectedEnd("slice", cursor)
	}
}
