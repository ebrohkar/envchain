package merger_test

import (
	"testing"

	"github.com/yourusername/envchain/internal/merger"
)

func TestMerge_StrategyLast_OverwritesConflicts(t *testing.T) {
	m := merger.New(merger.StrategyLast)
	a := map[string]string{"FOO": "a", "BAR": "shared_a"}
	b := map[string]string{"BAZ": "b", "BAR": "shared_b"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["BAR"] != "shared_b" {
		t.Errorf("expected BAR=shared_b, got %q", result["BAR"])
	}
	if result["FOO"] != "a" {
		t.Errorf("expected FOO=a, got %q", result["FOO"])
	}
}

func TestMerge_StrategyFirst_KeepsOriginal(t *testing.T) {
	m := merger.New(merger.StrategyFirst)
	a := map[string]string{"FOO": "original"}
	b := map[string]string{"FOO": "override"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %q", result["FOO"])
	}
}

func TestMerge_StrategyError_ConflictReturnsError(t *testing.T) {
	m := merger.New(merger.StrategyError)
	a := map[string]string{"FOO": "one"}
	b := map[string]string{"FOO": "two"}

	_, err := m.Merge(a, b)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_StrategyError_SameValueNoError(t *testing.T) {
	m := merger.New(merger.StrategyError)
	a := map[string]string{"FOO": "same"}
	b := map[string]string{"FOO": "same"}

	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
	if result["FOO"] != "same" {
		t.Errorf("expected FOO=same, got %q", result["FOO"])
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	m := merger.New(merger.StrategyLast)
	result, err := m.Merge(map[string]string{}, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestKeys_ReturnsAllUniqueKeys(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "3", "C": "4"}

	keys := merger.Keys(a, b)
	if len(keys) != 3 {
		t.Errorf("expected 3 unique keys, got %d: %v", len(keys), keys)
	}
}
