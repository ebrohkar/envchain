package partitioner_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/partitioner"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "supersecretpassword123",
		"PORT":        "8080",
		"DEBUG":       "true",
		"API_KEY":     "abc123xyz",
	}
}

func TestPartition_ByValueLength(t *testing.T) {
	p := partitioner.New(partitioner.ByValueLength(8, 16))
	result := p.Partition(baseEnv())

	if _, ok := result["short"]; !ok {
		t.Fatal("expected 'short' partition")
	}
	if _, ok := result["long"]; !ok {
		t.Fatal("expected 'long' partition")
	}
	// "true" and "8080" are short (<= 8)
	if _, ok := result["short"]["DEBUG"]; !ok {
		t.Errorf("expected DEBUG in short partition")
	}
	if _, ok := result["short"]["PORT"]; !ok {
		t.Errorf("expected PORT in short partition")
	}
	// "supersecretpassword123" is long (> 16)
	if _, ok := result["long"]["DB_PASSWORD"]; !ok {
		t.Errorf("expected DB_PASSWORD in long partition")
	}
}

func TestPartition_CustomFn(t *testing.T) {
	p := partitioner.New(func(key, _ string) string {
		if len(key) > 5 {
			return "verbose"
		}
		return "terse"
	})
	result := p.Partition(baseEnv())

	if _, ok := result["terse"]["PORT"]; !ok {
		t.Errorf("expected PORT in terse partition")
	}
	if _, ok := result["verbose"]["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST in verbose partition")
	}
}

func TestPartitionNames_ReturnsSorted(t *testing.T) {
	partitions := map[string]map[string]string{
		"zebra": {"Z": "1"},
		"alpha": {"A": "2"},
		"mango": {"M": "3"},
	}
	names := partitioner.PartitionNames(partitions)
	expected := []string{"alpha", "mango", "zebra"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, name, expected[i])
		}
	}
}

func TestPartition_EmptyEnv(t *testing.T) {
	p := partitioner.New(partitioner.ByValueLength(8, 16))
	result := p.Partition(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d partitions", len(result))
	}
}

func TestNewWithOptions_DefaultPartition(t *testing.T) {
	p := partitioner.NewWithOptions()
	result := p.Partition(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if len(result) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(result))
	}
	if _, ok := result["all"]; !ok {
		t.Errorf("expected default 'all' partition")
	}
	if len(result["all"]) != 2 {
		t.Errorf("expected 2 keys in 'all', got %d", len(result["all"]))
	}
}
