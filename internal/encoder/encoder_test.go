package encoder_test

import (
	"testing"

	"github.com/user/envchain/internal/encoder"
)

func TestEncoder_Base64_AllKeys(t *testing.T) {
	env := map[string]string{
		"SECRET": "hello",
		"TOKEN":  "world",
	}
	enc := encoder.New(encoder.FormatBase64)
	out, err := enc.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "aGVsbG8=" {
		t.Errorf("expected base64 of 'hello', got %q", out["SECRET"])
	}
	if out["TOKEN"] != "d29ybGQ=" {
		t.Errorf("expected base64 of 'world', got %q", out["TOKEN"])
	}
}

func TestEncoder_Base64_SelectiveKeys(t *testing.T) {
	env := map[string]string{
		"SECRET": "hello",
		"PLAIN":  "unchanged",
	}
	enc := encoder.New(encoder.FormatBase64, "SECRET")
	out, err := enc.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "aGVsbG8=" {
		t.Errorf("expected encoded SECRET, got %q", out["SECRET"])
	}
	if out["PLAIN"] != "unchanged" {
		t.Errorf("expected PLAIN to be unchanged, got %q", out["PLAIN"])
	}
}

func TestEncoder_Hex_Output(t *testing.T) {
	env := map[string]string{"KEY": "ab"}
	enc := encoder.New(encoder.FormatHex)
	out, err := enc.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "6162" {
		t.Errorf("expected hex '6162', got %q", out["KEY"])
	}
}

func TestEncoder_URL_SpecialChars(t *testing.T) {
	env := map[string]string{"REDIRECT": "https://example.com/path?q=1&r=2"}
	enc := encoder.New(encoder.FormatURL)
	out, err := enc.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["REDIRECT"] == env["REDIRECT"] {
		t.Error("expected URL encoding to transform special characters")
	}
}

func TestEncoder_CaseInsensitiveKeyMatch(t *testing.T) {
	env := map[string]string{"db_password": "secret"}
	enc := encoder.New(encoder.FormatBase64, "DB_PASSWORD")
	out, err := enc.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_password"] != "c2VjcmV0" {
		t.Errorf("expected encoded value for lowercase key, got %q", out["db_password"])
	}
}

func TestEncoder_EmptyEnv_ReturnsEmpty(t *testing.T) {
	enc := encoder.New(encoder.FormatBase64)
	out, err := enc.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
