package labeler

import (
	"sort"
	"testing"
)

func TestAttach_And_Get(t *testing.T) {
	l := New()
	l.Attach("DB_HOST", "tier", "critical")
	l.Attach("DB_HOST", "owner", "platform")

	labels := l.Get("DB_HOST")
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	if labels[0].Key != "tier" || labels[0].Value != "critical" {
		t.Errorf("unexpected first label: %+v", labels[0])
	}
}

func TestGet_ReturnsEmptySlice_WhenNoLabels(t *testing.T) {
	l := New()
	labels := l.Get("MISSING_KEY")
	if len(labels) != 0 {
		t.Errorf("expected empty slice, got %d labels", len(labels))
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	l := New()
	l.Attach("API_KEY", "env", "prod")

	copy1 := l.Get("API_KEY")
	copy1[0].Value = "mutated"

	copy2 := l.Get("API_KEY")
	if copy2[0].Value != "prod" {
		t.Errorf("Get should return a copy, but original was mutated")
	}
}

func TestFindByLabel_ReturnsMatchingKeys(t *testing.T) {
	l := New()
	l.Attach("DB_HOST", "tier", "critical")
	l.Attach("DB_PORT", "tier", "critical")
	l.Attach("LOG_LEVEL", "tier", "optional")

	keys := l.FindByLabel("tier", "critical")
	sort.Strings(keys)

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PORT" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFindByLabel_ReturnsEmpty_WhenNoMatch(t *testing.T) {
	l := New()
	l.Attach("API_KEY", "env", "staging")

	keys := l.FindByLabel("env", "prod")
	if len(keys) != 0 {
		t.Errorf("expected no matches, got %v", keys)
	}
}

func TestRemove_DeletesLabels(t *testing.T) {
	l := New()
	l.Attach("DB_HOST", "tier", "critical")
	l.Remove("DB_HOST")

	if len(l.Get("DB_HOST")) != 0 {
		t.Error("expected labels to be removed")
	}
}

func TestSummary_WithLabels(t *testing.T) {
	l := New()
	l.Attach("DB_HOST", "tier", "critical")

	s := l.Summary("DB_HOST")
	if s == "" {
		t.Error("expected non-empty summary")
	}
	if s == "DB_HOST: (no labels)" {
		t.Error("expected labels in summary, got none")
	}
}

func TestSummary_NoLabels(t *testing.T) {
	l := New()
	s := l.Summary("UNKNOWN_KEY")
	expected := "UNKNOWN_KEY: (no labels)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
