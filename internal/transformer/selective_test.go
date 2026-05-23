package transformer_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/transformer"
)

func TestSelective_OnlyTransformsMatchingKeys(t *testing.T) {
	st := transformer.NewSelective(func(key string) bool {
		return strings.HasPrefix(key, "APP_")
	}).Add("upper", transformer.ToUpper())

	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
	}

	result, err := st.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_NAME"] != "MYAPP" {
		t.Errorf("expected 'MYAPP', got %q", result["APP_NAME"])
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected 'localhost' unchanged, got %q", result["DB_HOST"])
	}
}

func TestSelective_NoMatchingKeys_PassesThrough(t *testing.T) {
	st := transformer.NewSelective(func(key string) bool {
		return strings.HasPrefix(key, "MISSING_")
	}).Add("upper", transformer.ToUpper())

	env := map[string]string{"KEY": "value", "OTHER": "data"}

	result, err := st.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "value" || result["OTHER"] != "data" {
		t.Errorf("expected unchanged values, got %v", result)
	}
}

func TestSelective_AllKeysMatch(t *testing.T) {
	st := transformer.NewSelective(func(key string) bool {
		return true
	}).Add("trim", transformer.TrimSpace())

	env := map[string]string{"A": "  foo  ", "B": " bar "}

	result, err := st.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "foo" || result["B"] != "bar" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestSelective_DoesNotMutateInput(t *testing.T) {
	st := transformer.NewSelective(func(key string) bool {
		return true
	}).Add("upper", transformer.ToUpper())

	env := map[string]string{"KEY": "hello"}

	_, err := st.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "hello" {
		t.Errorf("input map was mutated, expected 'hello', got %q", env["KEY"])
	}
}
