package errors

import (
	"reflect"
	"strings"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	intType := reflect.TypeOf(0)
	ptrType := reflect.TypeOf((*int)(nil))
	tests := []error{
		&InvalidUnmarshalError{},
		&InvalidUnmarshalError{Type: intType},
		&InvalidUnmarshalError{Type: ptrType},
		ErrSyntax("syntax", 1),
		&UnmarshalTypeError{Value: "string", Type: intType, Offset: 2},
		&UnmarshalTypeError{Value: "string", Type: intType, Offset: 2, Struct: "S", Field: "F"},
		&UnsupportedTypeError{Type: intType},
		&UnsupportedValueError{Value: reflect.ValueOf(1), Str: "bad"},
		ErrExceededMaxDepth('a', 3),
		ErrUnexpectedStart("value", []byte("x"), 0),
		ErrUnexpectedEnd("value", 0),
		ErrUnexpected("value", 0, 'x'),
		ErrInvalidCharacter(0, "value", 0),
		ErrInvalidCharacter('x', "value", 0),
		ErrInvalidBeginningOfValue('x', 0),
		ErrInvalidBeginningOfArray('x', 0),
		ErrOverflow(256, "uint8"),
	}
	for _, err := range tests {
		if err == nil || strings.TrimSpace(err.Error()) == "" {
			t.Errorf("error %#v has an empty message", err)
		}
	}
}
