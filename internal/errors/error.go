package errors

import (
	"fmt"
	"reflect"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "php-serialize: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return fmt.Sprintf("php-serialize: Unmarshal(non-pointer %s)", e.Type)
	}
	return fmt.Sprintf("php-serialize: Unmarshal(nil %s)", e.Type)
}

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return e.msg }

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return fmt.Sprintf("php-serialize: cannot unmarshal %s into Go struct field %s.%s of type %s (offset %d)",
			e.Value, e.Struct, e.Field, e.Type, e.Offset,
		)
	}
	return fmt.Sprintf("php-serialize: cannot unmarshal %s into Go value of type %s (offset: %d)", e.Value, e.Type, e.Offset)
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("php-serialize: unsupported type: %s", e.Type)
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return fmt.Sprintf("php-serialize: unsupported value: %s", e.Str)
}

func ErrSyntax(msg string, offset int64) *SyntaxError {
	return &SyntaxError{msg: msg, Offset: offset}
}

func ErrExceededMaxDepth(c byte, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf(`php-serialize: invalid character "%c" exceeded max depth`, c),
		Offset: cursor,
	}
}

func ErrUnexpectedStart(typ string, buf []byte, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf("php-serialize: unexpected %c at beginning of %s", buf[cursor], typ),
		Offset: cursor,
	}
}

func ErrUnexpectedEnd(msg string, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf("php-serialize: %s unexpected end of input", msg),
		Offset: cursor,
	}
}

func ErrUnexpected(msg string, cursor int64, c byte) *SyntaxError {
	return &SyntaxError{msg: fmt.Sprintf("php-serialize: expecting %s, get char '%c' instead", msg, c), Offset: cursor}
}

func ErrInvalidCharacter(c byte, context string, cursor int64) *SyntaxError {
	if c == 0 {
		return &SyntaxError{
			msg:    fmt.Sprintf("php-serialize: invalid character as %s", context),
			Offset: cursor,
		}
	}
	return &SyntaxError{
		msg:    fmt.Sprintf("php-serialize: invalid character %c as %s", c, context),
		Offset: cursor,
	}
}

func ErrInvalidBeginningOfValue(c byte, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf("php-serialize: invalid character '%c' looking for beginning of value", c),
		Offset: cursor,
	}
}

func ErrInvalidBeginningOfArray(c byte, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf("php-serialize: invalid character '%c' looking for beginning of array", c),
		Offset: cursor,
	}
}

func ErrOverflow(v any, t string) error {
	return &overflowError{
		v: v,
		t: t,
	}
}

type overflowError struct {
	t string
	v any
}

func (o overflowError) Error() string {
	return fmt.Sprintf("php-serialize: %v overflow type %s", o.v, o.t)
}
