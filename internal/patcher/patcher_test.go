package patcher_test

import (
	"testing"

	"github.com/envchain/envchain/internal/patcher"
)

func TestApply_StrategyOverwrite_ReplacesExisting(t *testing.T) {
	p := patcher.New(patcher.StrategyOverwrite)
	base := map[string]string{"A": "1", "B": "2"}
	patch := map[string]string{"B": "99", "C": "3"}

	result, err := p.Apply(base, patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99, got %q", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %q", result["C"])
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %q", result["A"])
	}
}

func TestApply_StrategySkipExisting_KeepsOriginal(t *testing.T) {
	p := patcher.New(patcher.StrategySkipExisting)
	base := map[string]string{"A": "original"}
	patch := map[string]string{"A": "patched", "B": "new"}

	result, err := p.Apply(base, patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "original" {
		t.Errorf("expected A=original, got %q", result["A"])
	}
	if result["B"] != "new" {
		t.Errorf("expected B=new, got %q", result["B"])
	}
}

func TestApply_StrategyErrorOnConflict_ReturnsError(t *testing.T) {
	p := patcher.New(patcher.StrategyErrorOnConflict)
	base := map[string]string{"A": "1"}
	patch := map[string]string{"A": "2"}

	_, err := p.Apply(base, patch)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestApply_StrategyErrorOnConflict_SameValueNoError(t *testing.T) {
	p := patcher.New(patcher.StrategyErrorOnConflict)
	base := map[string]string{"A": "same"}
	patch := map[string]string{"A": "same"}

	result, err := p.Apply(base, patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "same" {
		t.Errorf("expected A=same, got %q", result["A"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	p := patcher.New(patcher.StrategyOverwrite)
	base := map[string]string{"A": "1"}
	patch := map[string]string{"A": "2"}

	_, _ = p.Apply(base, patch)
	if base["A"] != "1" {
		t.Errorf("base map was mutated: A=%q", base["A"])
	}
}

func TestKeys_ReturnsChangedKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	patch := map[string]string{"B": "99", "C": "new"}

	keys := patcher.Keys(base, patch)
	if len(keys) != 2 {
		t.Errorf("expected 2 changed keys, got %d: %v", len(keys), keys)
	}
}
