package tagger

import (
	"sort"
	"testing"
)

func TestTag_And_Get(t *testing.T) {
	tr := New()
	tr.Tag("DB_HOST", "env", "production")
	tr.Tag("DB_HOST", "team", "platform")

	tags := tr.Get("DB_HOST")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Key != "env" || tags[0].Value != "production" {
		t.Errorf("unexpected first tag: %+v", tags[0])
	}
	if tags[1].Key != "team" || tags[1].Value != "platform" {
		t.Errorf("unexpected second tag: %+v", tags[1])
	}
}

func TestGet_ReturnsEmptySlice_WhenNoTags(t *testing.T) {
	tr := New()
	tags := tr.Get("MISSING_KEY")
	if len(tags) != 0 {
		t.Errorf("expected empty slice, got %d tags", len(tags))
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	tr := New()
	tr.Tag("API_KEY", "sensitivity", "high")

	tags := tr.Get("API_KEY")
	tags[0].Value = "mutated"

	original := tr.Get("API_KEY")
	if original[0].Value == "mutated" {
		t.Error("Get should return a copy, not a reference")
	}
}

func TestFindByTag_ReturnsMatchingKeys(t *testing.T) {
	tr := New()
	tr.Tag("DB_HOST", "env", "production")
	tr.Tag("DB_PORT", "env", "production")
	tr.Tag("API_KEY", "env", "staging")

	matches := tr.FindByTag("env", "production")
	sort.Strings(matches)

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0] != "DB_HOST" || matches[1] != "DB_PORT" {
		t.Errorf("unexpected matches: %v", matches)
	}
}

func TestFindByTag_ReturnsEmpty_WhenNoMatch(t *testing.T) {
	tr := New()
	tr.Tag("DB_HOST", "env", "production")

	matches := tr.FindByTag("env", "staging")
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %v", matches)
	}
}

func TestKeys_ReturnsTaggedKeys(t *testing.T) {
	tr := New()
	tr.Tag("FOO", "k", "v")
	tr.Tag("BAR", "k", "v")

	keys := tr.Keys()
	sort.Strings(keys)

	if len(keys) != 2 || keys[0] != "BAR" || keys[1] != "FOO" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestClear_RemovesTags(t *testing.T) {
	tr := New()
	tr.Tag("SECRET", "type", "credential")
	tr.Clear("SECRET")

	if len(tr.Get("SECRET")) != 0 {
		t.Error("expected tags to be cleared")
	}
}
