package differ_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/differ"
)

func TestCompare_NoChanges(t *testing.T) {
	d := differ.New()
	before := map[string]string{"FOO": "bar", "BAZ": "qux"}
	after := map[string]string{"FOO": "bar", "BAZ": "qux"}

	diffs := d.Compare(before, after)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestCompare_Added(t *testing.T) {
	d := differ.New()
	before := map[string]string{"FOO": "bar"}
	after := map[string]string{"FOO": "bar", "NEW_VAR": "hello"}

	diffs := d.Compare(before, after)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != differ.Added || diffs[0].Key != "NEW_VAR" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	d := differ.New()
	before := map[string]string{"FOO": "bar", "OLD_VAR": "gone"}
	after := map[string]string{"FOO": "bar"}

	diffs := d.Compare(before, after)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != differ.Removed || diffs[0].Key != "OLD_VAR" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompare_Changed(t *testing.T) {
	d := differ.New()
	before := map[string]string{"FOO": "old"}
	after := map[string]string{"FOO": "new"}

	diffs := d.Compare(before, after)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != differ.Changed || diffs[0].OldValue != "old" || diffs[0].NewValue != "new" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestHasChanges(t *testing.T) {
	d := differ.New()
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "2"}

	if !d.HasChanges(before, after) {
		t.Error("expected HasChanges to return true")
	}

	if d.HasChanges(before, before) {
		t.Error("expected HasChanges to return false for identical maps")
	}
}
