package sorter_test

import (
	"testing"

	"github.com/envchain/envchain/internal/sorter"
)

var sampleEnv = map[string]string{
	"ZEBRA":    "z",
	"ALPHA":    "long_value_here",
	"MONGO_URI": "mongodb://localhost",
	"PORT":     "8080",
}

func TestSort_ByKey(t *testing.T) {
	s := sorter.New(sorter.ByKey)
	keys := s.Sort(sampleEnv)

	expected := []string{"ALPHA", "MONGO_URI", "PORT", "ZEBRA"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestSort_ByValueLength(t *testing.T) {
	s := sorter.New(sorter.ByValueLength)
	keys := s.Sort(sampleEnv)

	// values by length: "z"(1), "8080"(4), "long_value_here"(15), "mongodb://localhost"(18)
	expected := []string{"ZEBRA", "PORT", "ALPHA", "MONGO_URI"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestSort_ByPriority_MatchedFirst(t *testing.T) {
	s := sorter.New(sorter.ByPriority).WithPriority([]string{"PORT", "MONGO_URI"})
	keys := s.Sort(sampleEnv)

	if keys[0] != "PORT" {
		t.Errorf("expected PORT first, got %q", keys[0])
	}
	if keys[1] != "MONGO_URI" {
		t.Errorf("expected MONGO_URI second, got %q", keys[1])
	}
	// remaining keys should be alphabetical
	if keys[2] != "ALPHA" {
		t.Errorf("expected ALPHA third, got %q", keys[2])
	}
	if keys[3] != "ZEBRA" {
		t.Errorf("expected ZEBRA fourth, got %q", keys[3])
	}
}

func TestSort_ByPriority_NoPrioritySet(t *testing.T) {
	s := sorter.New(sorter.ByPriority)
	keys := s.Sort(sampleEnv)

	// without priority list behaves like ByKey
	expected := []string{"ALPHA", "MONGO_URI", "PORT", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestSort_EmptyMap(t *testing.T) {
	s := sorter.New(sorter.ByKey)
	keys := s.Sort(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
