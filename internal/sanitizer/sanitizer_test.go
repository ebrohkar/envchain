package sanitizer_test

import (
	"testing"

	"github.com/envchain/envchain/internal/sanitizer"
)

func TestSanitizeKey_Uppercase(t *testing.T) {
	s := sanitizer.New()
	if got := s.SanitizeKey("my_var"); got != "MY_VAR" {
		t.Errorf("expected MY_VAR, got %s", got)
	}
}

func TestSanitizeKey_ReplacesInvalidChars(t *testing.T) {
	s := sanitizer.New()
	if got := s.SanitizeKey("my-var.name"); got != "MY_VAR_NAME" {
		t.Errorf("expected MY_VAR_NAME, got %s", got)
	}
}

func TestSanitizeKey_TrimsSpace(t *testing.T) {
	s := sanitizer.New()
	if got := s.SanitizeKey("  DB_HOST  "); got != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", got)
	}
}

func TestSanitizeKey_NormalizeDisabled(t *testing.T) {
	s := sanitizer.New(sanitizer.WithNormalizeKeys(false))
	if got := s.SanitizeKey("my-var"); got != "my-var" {
		t.Errorf("expected my-var unchanged, got %s", got)
	}
}

func TestSanitizeValue_TrimsSpace(t *testing.T) {
	s := sanitizer.New()
	if got := s.SanitizeValue("  hello  "); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestSanitizeValue_StripsControlChars(t *testing.T) {
	s := sanitizer.New()
	input := "hello\x01world\x7F"
	if got := s.SanitizeValue(input); got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}

func TestSanitizeValue_ControlCharsDisabled(t *testing.T) {
	s := sanitizer.New(sanitizer.WithStripControlChars(false))
	input := "hello\x01world"
	if got := s.SanitizeValue(input); got != "hello\x01world" {
		t.Errorf("expected control char preserved, got %q", got)
	}
}

func TestSanitizeMap_AllEntries(t *testing.T) {
	s := sanitizer.New()
	env := map[string]string{
		"db-host":   "  localhost  ",
		"api.token": "secret\x00value",
	}
	result := s.SanitizeMap(env)

	if v, ok := result["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", v)
	}
	if v, ok := result["API_TOKEN"]; !ok || v != "secretvalue" {
		t.Errorf("expected API_TOKEN=secretvalue, got %q", v)
	}
}

func TestSanitizeMap_DoesNotMutateInput(t *testing.T) {
	s := sanitizer.New()
	env := map[string]string{"MY_KEY": "  value  "}
	_ = s.SanitizeMap(env)
	if env["MY_KEY"] != "  value  " {
		t.Error("input map was mutated")
	}
}
