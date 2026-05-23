package redactor

import (
	"testing"
)

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "supersecret",
		"APP_NAME": "envchain",
	}
	r := New([]string{"API_KEY"})
	out := r.Apply(env)

	if out["API_KEY"] != defaultPlaceholder {
		t.Errorf("expected %q, got %q", defaultPlaceholder, out["API_KEY"])
	}
	if out["APP_NAME"] != "envchain" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out["APP_NAME"])
	}
}

func TestApply_CaseInsensitiveKeyMatch(t *testing.T) {
	env := map[string]string{"db_password": "s3cr3t"}
	r := New([]string{"DB_PASSWORD"})
	out := r.Apply(env)

	if out["db_password"] != defaultPlaceholder {
		t.Errorf("expected redacted value, got %q", out["db_password"])
	}
}

func TestApply_PartialStrategy(t *testing.T) {
	env := map[string]string{"SECRET_TOKEN": "abcdefghij"}
	r := New([]string{"SECRET_TOKEN"}).WithStrategy(StrategyPartial)
	out := r.Apply(env)

	expected := "abcd******"
	if out["SECRET_TOKEN"] != expected {
		t.Errorf("expected %q, got %q", expected, out["SECRET_TOKEN"])
	}
}

func TestApply_PartialStrategy_ShortValue(t *testing.T) {
	env := map[string]string{"SECRET_TOKEN": "abc"}
	r := New([]string{"SECRET_TOKEN"}).WithStrategy(StrategyPartial)
	out := r.Apply(env)

	if out["SECRET_TOKEN"] != defaultPlaceholder {
		t.Errorf("short value should fall back to full redaction, got %q", out["SECRET_TOKEN"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"TOKEN": "mytoken"}
	r := New([]string{"TOKEN"}).WithPlaceholder("***")
	out := r.Apply(env)

	if out["TOKEN"] != "***" {
		t.Errorf("expected custom placeholder, got %q", out["TOKEN"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"PASSWORD": "original"}
	r := New([]string{"PASSWORD"})
	_ = r.Apply(env)

	if env["PASSWORD"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_NoSensitiveKeys_ReturnsAll(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := New([]string{})
	out := r.Apply(env)

	if len(out) != 2 || out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Error("expected all keys returned unchanged")
	}
}
