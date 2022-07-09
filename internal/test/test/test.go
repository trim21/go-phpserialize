package test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Case struct {
	Name     string
	Data     interface{}
	Expected string `php:"-" json:"-"`
}

func (tc Case) WrappedExpected() string {
	return fmt.Sprintf(`a:2:{s:4:"Name";s:%d:"%s";s:4:"Data";`, len(tc.Name), tc.Name) + tc.Expected + "}"
}

type User struct {
	ID   uint64 `php:"id" json:"id"`
	Name string `php:"name" json:"name"`
}

func StringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("Result not as expected:\n%v", CharacterDiff(expected, actual))
		t.FailNow()
	}
}

func diff(a, b string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, true)
	if len(diffs) > 2 {
		diffs = dmp.DiffCleanupSemantic(diffs)
		diffs = dmp.DiffCleanupEfficiency(diffs)
	}
	return diffs
}

// CharacterDiff returns an inline diff between the two strings, using (++added++) and (~~deleted~~) markup.
func CharacterDiff(a, b string) string {
	return diffsToString(diff(a, b))
}

func diffsToString(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			buff.WriteString(color.RedString(text))
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}

	buff.WriteString("\n")

	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString(color.GreenString(text))
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}

	return buff.String()
}
