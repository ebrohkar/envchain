package registry

import (
	"sort"
	"testing"
)

func TestRegister_And_Lookup(t *testing.T) {
	r := New()
	e := Entry{Key: "DB_HOST", Description: "Database host", Required: true}
	r.Register(e)

	got, err := r.Lookup("DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != "Database host" {
		t.Errorf("expected description %q, got %q", "Database host", got.Description)
	}
}

func TestLookup_NotFound(t *testing.T) {
	r := New()
	_, err := r.Lookup("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRegister_Overwrites_Existing(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "API_URL", Description: "old"})
	r.Register(Entry{Key: "API_URL", Description: "new"})

	got, _ := r.Lookup("API_URL")
	if got.Description != "new" {
		t.Errorf("expected %q, got %q", "new", got.Description)
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "A"})
	r.Register(Entry{Key: "B"})
	r.Register(Entry{Key: "C"})

	keys := r.Keys()
	sort.Strings(keys)
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestRequiredKeys_FiltersCorrectly(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "REQ_KEY", Required: true})
	r.Register(Entry{Key: "OPT_KEY", Required: false})

	req := r.RequiredKeys()
	if len(req) != 1 || req[0] != "REQ_KEY" {
		t.Errorf("expected [REQ_KEY], got %v", req)
	}
}

func TestSize_ReturnsCount(t *testing.T) {
	r := New()
	if r.Size() != 0 {
		t.Errorf("expected 0, got %d", r.Size())
	}
	r.Register(Entry{Key: "X"})
	r.Register(Entry{Key: "Y"})
	if r.Size() != 2 {
		t.Errorf("expected 2, got %d", r.Size())
	}
}
