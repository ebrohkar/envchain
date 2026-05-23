package partitioner_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/partitioner"
)

func TestNewWithOptions_WithCustomFn(t *testing.T) {
	called := false
	p := partitioner.NewWithOptions(
		partitioner.WithFn(func(key, _ string) string {
			called = true
			return key
		}),
	)
	p.Partition(map[string]string{"X": "y"})
	if !called {
		t.Error("expected custom partition function to be called")
	}
}

func TestNewWithOptions_OverridesDefault(t *testing.T) {
	p := partitioner.NewWithOptions(
		partitioner.WithFn(func(_, _ string) string { return "custom" }),
	)
	result := p.Partition(map[string]string{"A": "1", "B": "2"})
	if _, ok := result["custom"]; !ok {
		t.Fatal("expected 'custom' partition")
	}
	if _, ok := result["all"]; ok {
		t.Error("expected default 'all' partition to be overridden")
	}
}
