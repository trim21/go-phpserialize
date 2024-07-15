package runtime

import (
	"reflect"
	"strings"
	"unicode"
)

func getTag(field reflect.StructField) string {
	return field.Tag.Get("php")
}

func IsIgnoredStructField(field reflect.StructField) bool {
	if !field.IsExported() {
		return true
	}

	tag := getTag(field)

	return tag == "-"
}

type StructTag struct {
	Key         string
	IsOmitEmpty bool
	IsString    bool
	Field       reflect.StructField
}

func (s StructTag) Name() string {
	if s.Key != "" {
		return s.Key
	}
	return s.Field.Name
}

type StructTags []*StructTag

func (t StructTags) ExistsKey(key string) bool {
	for _, tt := range t {
		if tt.Key == key {
			return true
		}
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

func StructTagFromField(field reflect.StructField) *StructTag {
	keyName := field.Name
	tag := getTag(field)
	st := &StructTag{Field: field}
	opts := strings.Split(tag, ",")
	if len(opts) > 0 {
		if opts[0] != "" && isValidTag(opts[0]) {
			keyName = opts[0]
		}
	}
	st.Key = keyName
	if len(opts) > 1 {
		for _, opt := range opts[1:] {
			switch opt {
			case "omitempty":
				st.IsOmitEmpty = true
			case "string":
				st.IsString = true
			}
		}
	}
	return st
}
