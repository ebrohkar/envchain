package registry

import "testing"

func makeEnv(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func TestCheck_AllPresent(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "DB_HOST", Required: true})
	r.Register(Entry{Key: "DB_PORT", Required: true})

	env := makeEnv(map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	res := r.Check(env)

	if !res.AllPresent() {
		t.Errorf("expected all present, got missing: %v", res.Missing)
	}
	if res.Summary() != "all required variables present" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestCheck_MissingRequired(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "SECRET_KEY", Required: true, Description: "App secret"})

	env := makeEnv(map[string]string{})
	res := r.Check(env)

	if res.AllPresent() {
		t.Fatal("expected missing keys")
	}
	if len(res.Missing) != 1 || res.Missing[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected missing: %v", res.Missing)
	}
}

func TestCheck_FallsBackToDefault(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "LOG_LEVEL", Required: true, Default: "info"})

	env := makeEnv(map[string]string{})
	res := r.Check(env)

	if !res.AllPresent() {
		t.Errorf("expected all present via default, got missing: %v", res.Missing)
	}
	if res.Defaults["LOG_LEVEL"] != "info" {
		t.Errorf("expected default 'info', got %q", res.Defaults["LOG_LEVEL"])
	}
}

func TestCheck_OptionalKeyIgnored(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "OPTIONAL_FEATURE", Required: false})

	env := makeEnv(map[string]string{})
	res := r.Check(env)

	if !res.AllPresent() {
		t.Errorf("optional key should not appear in missing: %v", res.Missing)
	}
}

func TestCheck_Summary_WithMissing(t *testing.T) {
	r := New()
	r.Register(Entry{Key: "A", Required: true})
	r.Register(Entry{Key: "B", Required: true})

	env := makeEnv(map[string]string{})
	res := r.Check(env)

	expected := "2 required variable(s) missing"
	if res.Summary() != expected {
		t.Errorf("expected %q, got %q", expected, res.Summary())
	}
}
