package scoper_test

import (
	"sort"
	"testing"

	"github.com/envchain/envchain/internal/scoper"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"APP_PORT":   "8080",
		"APP_SECRET": "s3cr3t",
		"NOPREFIX":   "value",
	}
}

func TestScope_ReturnsMatchingKeys(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	got := s.Scope("DB")
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", got["HOST"])
	}
	if got["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", got["PORT"])
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestScope_EmptyWhenNoMatch(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	got := s.Scope("CACHE")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestScope_CaseInsensitiveNamespace(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	got := s.Scope("app")
	if got["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", got["PORT"])
	}
}

func TestNamespaces_ReturnsDistinctPrefixes(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	ns := s.Namespaces()
	sort.Strings(ns)
	if len(ns) != 2 {
		t.Errorf("expected 2 namespaces, got %d: %v", len(ns), ns)
	}
	if ns[0] != "APP" || ns[1] != "DB" {
		t.Errorf("unexpected namespaces: %v", ns)
	}
}

func TestHas_ReturnsTrueForExistingKey(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	if !s.Has("DB", "HOST") {
		t.Error("expected Has(DB, HOST) to be true")
	}
}

func TestHas_ReturnsFalseForMissingKey(t *testing.T) {
	s := scoper.New(baseEnv(), "_")
	if s.Has("DB", "PASSWORD") {
		t.Error("expected Has(DB, PASSWORD) to be false")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	env := map[string]string{"SVC_KEY": "val"}
	s := scoper.New(env, "")
	got := s.Scope("SVC")
	if got["KEY"] != "val" {
		t.Errorf("expected KEY=val with default separator, got %q", got["KEY"])
	}
}
