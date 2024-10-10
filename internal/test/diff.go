package test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func init() {
	color.NoColor = false // force color
}

func diff(a, b string) []diffmatchpatch.Diff {
	return diffmatchpatch.New().DiffMain(a, b, false)
}

// characterDiff returns an inline diff between the two strings, using (++added++) and (~~deleted~~) markup.
func characterDiff(a, b string) string {
	return diffsToString(diff(a, b))
}

func diffsToString(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	buff.WriteString("expected : ")
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			buff.WriteString(color.RedString(text))
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}

	buff.WriteString("\nactual   : ")

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

type stringLike interface {
	~string | ~[]byte
}

func StringEqual[A stringLike, B stringLike](t *testing.T, expected A, actual B) {
	t.Helper()

	if string(expected) == string(actual) {
		return
	}

	t.Errorf("Result not as expected:\n%v", characterDiff(strconv.QuoteToASCII(string(expected)), strconv.QuoteToASCII(string(actual))))
	t.FailNow()
}
