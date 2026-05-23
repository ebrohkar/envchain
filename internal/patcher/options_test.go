package patcher_test

import (
	"testing"

	"github.com/envchain/envchain/internal/patcher"
)

func TestNewWithOptions_DefaultsToOverwrite(t *testing.T) {
	p := patcher.NewWithOptions()
	base := map[string]string{"X": "old"}
	patch := map[string]string{"X": "new"}

	result, err := p.Apply(base, patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["X"] != "new" {
		t.Errorf("expected X=new (overwrite default), got %q", result["X"])
	}
}

func TestNewWithOptions_WithStrategySkipExisting(t *testing.T) {
	p := patcher.NewWithOptions(patcher.WithStrategy(patcher.StrategySkipExisting))
	base := map[string]string{"X": "keep"}
	patch := map[string]string{"X": "ignore", "Y": "add"}

	result, err := p.Apply(base, patch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["X"] != "keep" {
		t.Errorf("expected X=keep, got %q", result["X"])
	}
	if result["Y"] != "add" {
		t.Errorf("expected Y=add, got %q", result["Y"])
	}
}

func TestNewWithOptions_WithStrategyErrorOnConflict(t *testing.T) {
	p := patcher.NewWithOptions(patcher.WithStrategy(patcher.StrategyErrorOnConflict))
	base := map[string]string{"Z": "a"}
	patch := map[string]string{"Z": "b"}

	_, err := p.Apply(base, patch)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}
