package test

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func init() {
	color.NoColor = false // force color
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

func StringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("Result not as expected:\n%v", CharacterDiff(expected, actual))
		t.FailNow()
	}
}
