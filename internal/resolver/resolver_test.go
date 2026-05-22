package resolver_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/resolver"
)

func mockEnv(m map[string]string) resolver.EnvGetter {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func TestResolve_NoReferences(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{}))
	got, missing := r.Resolve("plainvalue")
	if got != "plainvalue" {
		t.Errorf("expected 'plainvalue', got %q", got)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestResolve_SingleReference(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{"HOST": "localhost"}))
	got, missing := r.Resolve("${HOST}:5432")
	if got != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %q", got)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestResolve_MultipleReferences(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{"USER": "admin", "PASS": "secret"}))
	got, missing := r.Resolve("${USER}:${PASS}@db")
	if got != "admin:secret@db" {
		t.Errorf("expected 'admin:secret@db', got %q", got)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing vars, got %v", missing)
	}
}

func TestResolve_MissingReference(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{}))
	got, missing := r.Resolve("${MISSING_VAR}")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
	if len(missing) != 1 || missing[0] != "MISSING_VAR" {
		t.Errorf("expected [MISSING_VAR], got %v", missing)
	}
}

func TestResolveAll_HasMissing(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{"DB_HOST": "localhost"}))
	vars := map[string]string{
		"DSN":     "${DB_HOST}:${DB_PORT}/mydb",
		"API_URL": "https://api.example.com",
	}
	results := r.ResolveAll(vars)
	if !resolver.HasMissing(results) {
		t.Error("expected HasMissing to return true")
	}
	summary := resolver.MissingSummary(results)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestResolveAll_AllResolved(t *testing.T) {
	r := resolver.New(mockEnv(map[string]string{"A": "foo", "B": "bar"}))
	vars := map[string]string{
		"X": "${A}-${B}",
	}
	results := r.ResolveAll(vars)
	if resolver.HasMissing(results) {
		t.Error("expected no missing variables")
	}
	if len(results) != 1 || results[0].Resolved != "foo-bar" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestNew_NilGetenvPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil getenv")
		}
	}()
	resolver.New(nil)
}
