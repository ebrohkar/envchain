package renamer_test

import (
	"testing"

	"github.com/user/envchain/internal/renamer"
)

func TestApply_RenamesKeys(t *testing.T) {
	r := renamer.New()
	env := map[string]string{"OLD_KEY": "value", "KEEP": "kept"}
	mapping := map[string]string{"OLD_KEY": "NEW_KEY"}

	out, err := r.Apply(env, mapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["KEEP"] != "kept" {
		t.Errorf("expected KEEP=kept, got %q", out["KEEP"])
	}
}

func TestApply_MissingSourceKey_Skipped(t *testing.T) {
	r := renamer.New()
	env := map[string]string{"A": "1"}
	mapping := map[string]string{"MISSING": "NEW"}

	out, err := r.Apply(env, mapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["NEW"]; ok {
		t.Error("expected NEW to not be present")
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestApply_StrategySkip_KeepsDestination(t *testing.T) {
	r := renamer.NewWithStrategy(renamer.StrategySkip)
	env := map[string]string{"SRC": "new-val", "DST": "original"}
	mapping := map[string]string{"SRC": "DST"}

	out, err := r.Apply(env, mapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DST"] != "original" {
		t.Errorf("expected DST=original, got %q", out["DST"])
	}
}

func TestApply_StrategyError_ConflictReturnsError(t *testing.T) {
	r := renamer.NewWithStrategy(renamer.StrategyError)
	env := map[string]string{"SRC": "new-val", "DST": "original"}
	mapping := map[string]string{"SRC": "DST"}

	_, err := r.Apply(env, mapping)
	if err == nil {
		t.Fatal("expected error for conflicting destination key")
	}
}

func TestApply_StrategyOverwrite_ReplacesDestination(t *testing.T) {
	r := renamer.NewWithStrategy(renamer.StrategyOverwrite)
	env := map[string]string{"SRC": "new-val", "DST": "original"}
	mapping := map[string]string{"SRC": "DST"}

	out, err := r.Apply(env, mapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DST"] != "new-val" {
		t.Errorf("expected DST=new-val, got %q", out["DST"])
	}
}

func TestKeys_ReturnsSources(t *testing.T) {
	mapping := map[string]string{"A": "X", "B": "Y"}
	keys := renamer.Keys(mapping)
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
