package comparator_test

import (
	"testing"

	"github.com/envchain/envchain/internal/comparator"
)

func TestCompare_NoChanges(t *testing.T) {
	c := comparator.New()
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080"}
	current := map[string]string{"HOST": "localhost", "PORT": "8080"}

	result := c.Compare(baseline, current)

	if result.HasDrift() {
		t.Error("expected no drift but drift was detected")
	}
	if len(result.ByType(comparator.ChangeUnchanged)) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(result.ByType(comparator.ChangeUnchanged)))
	}
}

func TestCompare_AddedKey(t *testing.T) {
	c := comparator.New()
	baseline := map[string]string{"HOST": "localhost"}
	current := map[string]string{"HOST": "localhost", "PORT": "8080"}

	result := c.Compare(baseline, current)

	if !result.HasDrift() {
		t.Error("expected drift but none detected")
	}
	added := result.ByType(comparator.ChangeAdded)
	if len(added) != 1 || added[0].Key != "PORT" {
		t.Errorf("expected PORT to be added, got %+v", added)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	c := comparator.New()
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080"}
	current := map[string]string{"HOST": "localhost"}

	result := c.Compare(baseline, current)

	removed := result.ByType(comparator.ChangeRemoved)
	if len(removed) != 1 || removed[0].Key != "PORT" {
		t.Errorf("expected PORT to be removed, got %+v", removed)
	}
	if removed[0].OldValue != "8080" {
		t.Errorf("expected OldValue=8080, got %s", removed[0].OldValue)
	}
}

func TestCompare_UpdatedKey(t *testing.T) {
	c := comparator.New()
	baseline := map[string]string{"HOST": "localhost"}
	current := map[string]string{"HOST": "production.example.com"}

	result := c.Compare(baseline, current)

	updated := result.ByType(comparator.ChangeUpdated)
	if len(updated) != 1 {
		t.Fatalf("expected 1 updated change, got %d", len(updated))
	}
	if updated[0].OldValue != "localhost" || updated[0].NewValue != "production.example.com" {
		t.Errorf("unexpected update values: %+v", updated[0])
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	c := comparator.New()
	baseline := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	current := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}

	result := c.Compare(baseline, current)

	keys := make([]string, len(result.Changes))
	for i, ch := range result.Changes {
		keys[i] = ch.Key
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected key[%d]=%s, got %s", i, k, keys[i])
		}
	}
}
