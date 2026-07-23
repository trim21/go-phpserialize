package runtime

import (
	"reflect"
	"testing"
)

func TestStructFieldTags(t *testing.T) {
	type sample struct {
		Default string
		Named   int `php:"name,omitempty,string"`
		Invalid int `php:"bad\\name,omitempty"`
		Ignored int `php:"-"`
		private int
	}
	typ := reflect.TypeOf(sample{})

	defaultField, _ := typ.FieldByName("Default")
	defaultTag := StructTagFromField(defaultField)
	if defaultTag.Name() != "Default" || defaultTag.IsOmitEmpty || defaultTag.IsString {
		t.Fatalf("default tag = %#v", defaultTag)
	}

	namedField, _ := typ.FieldByName("Named")
	namedTag := StructTagFromField(namedField)
	if namedTag.Name() != "name" || !namedTag.IsOmitEmpty || !namedTag.IsString {
		t.Fatalf("named tag = %#v", namedTag)
	}

	invalidField, _ := typ.FieldByName("Invalid")
	if got := StructTagFromField(invalidField).Name(); got != "Invalid" {
		t.Fatalf("invalid tag name = %q", got)
	}

	ignoredField, _ := typ.FieldByName("Ignored")
	privateField, _ := typ.FieldByName("private")
	if !IsIgnoredStructField(ignoredField) || !IsIgnoredStructField(privateField) || IsIgnoredStructField(defaultField) {
		t.Fatal("IsIgnoredStructField returned an unexpected result")
	}

	tags := StructTags{defaultTag, namedTag}
	if !tags.ExistsKey("name") || tags.ExistsKey("missing") {
		t.Fatal("StructTags.ExistsKey returned an unexpected result")
	}

	for _, tag := range []string{"name", "with space", "a-b", "日本語"} {
		if !isValidTag(tag) {
			t.Errorf("isValidTag(%q) = false", tag)
		}
	}
	for _, tag := range []string{"", "bad\\name", `bad"name`} {
		if isValidTag(tag) {
			t.Errorf("isValidTag(%q) = true", tag)
		}
	}
}
