package decoder

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/errors"
	"github.com/trim21/go-phpserialize/internal/runtime"
)

var (
	sliceType = runtime.Type2RType(
		reflect.TypeOf((*sliceHeader)(nil)).Elem(),
	)
	nilSlice = unsafe.Pointer(&sliceHeader{})
)

type sliceDecoder struct {
	elemType          *runtime.Type
	isElemPointerType bool
	valueDecoder      Decoder
	size              uintptr
	arrayPool         sync.Pool
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

func newSliceDecoder(dec Decoder, elemType *runtime.Type, size uintptr, structName, fieldName string) *sliceDecoder {
	return &sliceDecoder{
		valueDecoder:      dec,
		elemType:          elemType,
		isElemPointerType: elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Map,
		size:              size,
		arrayPool: sync.Pool{
			New: func() interface{} {
				fmt.Println("slice decoder new slice")
				return &sliceHeader{
					data: newArray(elemType, defaultSliceCapacity),
					len:  0,
					cap:  defaultSliceCapacity,
				}
			},
		},
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *sliceDecoder) newSlice(src *sliceHeader, length int) *sliceHeader {
	slice := d.arrayPool.Get().(*sliceHeader)
	if src.len > 0 {
		// copy original elem
		if slice.cap < src.cap || src.cap < length {
			data := newArray(d.elemType, length)
			slice = &sliceHeader{data: data, len: src.len, cap: length}
		} else {
			slice.len = src.len
		}
		copySlice(d.elemType, *slice, *src)
	} else {
		slice.len = 0
	}
	return slice
}

func (d *sliceDecoder) releaseSlice(p *sliceHeader) {
	d.arrayPool.Put(p)
}

//go:linkname copySlice reflect.typedslicecopy
func copySlice(elemType *runtime.Type, dst, src sliceHeader) int

//go:linkname newArray reflect.unsafe_NewArray
func newArray(*runtime.Type, int) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t *runtime.Type, dst, src unsafe.Pointer)

func (d *sliceDecoder) errNumber(offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  "number",
		Type:   reflect.SliceOf(runtime.RType2Type(d.elemType)),
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

func (d *sliceDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
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
		typedmemmove(sliceType, p, nilSlice)
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

			dst := (*sliceHeader)(p)
			if dst.data == nil {
				dst.data = newArray(d.elemType, 0)
			} else {
				dst.len = 0
			}

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
