package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/differ"
)

func TestFormatter_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	f := differ.NewFormatter(&buf)
	f.Write(nil)

	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-change message, got: %q", buf.String())
	}
}

func TestFormatter_Added(t *testing.T) {
	var buf bytes.Buffer
	f := differ.NewFormatter(&buf)
	f.Write([]differ.Diff{
		{Key: "NEW_KEY", Kind: differ.Added, NewValue: "hello"},
	})

	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added line, got: %q", out)
	}
}

func TestFormatter_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := differ.NewFormatter(&buf)
	f.Write([]differ.Diff{
		{Key: "OLD_KEY", Kind: differ.Removed, OldValue: "bye"},
	})

	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removed line, got: %q", out)
	}
}

func TestFormatter_Changed(t *testing.T) {
	var buf bytes.Buffer
	f := differ.NewFormatter(&buf)
	f.Write([]differ.Diff{
		{Key: "FOO", Kind: differ.Changed, OldValue: "old", NewValue: "new"},
	})

	out := buf.String()
	if !strings.Contains(out, "~ FOO") {
		t.Errorf("expected changed line, got: %q", out)
	}
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("expected old and new values in output, got: %q", out)
	}
}

func TestFormatter_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	f := differ.NewFormatter(&buf)
	f.Write([]differ.Diff{
		{Key: "ZZZ", Kind: differ.Added, NewValue: "z"},
		{Key: "AAA", Kind: differ.Added, NewValue: "a"},
	})

	out := buf.String()
	idxA := strings.Index(out, "AAA")
	idxZ := strings.Index(out, "ZZZ")
	if idxA > idxZ {
		t.Errorf("expected AAA before ZZZ in sorted output")
	}
}
