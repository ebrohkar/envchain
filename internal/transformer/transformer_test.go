package transformer_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain/internal/transformer"
)

func TestApply_NoTransforms_ReturnsOriginal(t *testing.T) {
	tr := transformer.New()
	env := map[string]string{"KEY": "value", "OTHER": "data"}

	result, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "value" || result["OTHER"] != "data" {
		t.Errorf("expected original values, got %v", result)
	}
}

func TestApply_TrimSpace(t *testing.T) {
	tr := transformer.New().Add("trim", transformer.TrimSpace())
	env := map[string]string{"KEY": "  hello  ", "B": "\tworld\n"}

	result, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", result["KEY"])
	}
	if result["B"] != "world" {
		t.Errorf("expected 'world', got %q", result["B"])
	}
}

func TestApply_ToUpper(t *testing.T) {
	tr := transformer.New().Add("upper", transformer.ToUpper())
	env := map[string]string{"KEY": "hello"}

	result, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", result["KEY"])
	}
}

func TestApply_ChainedTransforms(t *testing.T) {
	tr := transformer.New().
		Add("trim", transformer.TrimSpace()).
		Add("upper", transformer.ToUpper())
	env := map[string]string{"KEY": "  hello  "}

	result, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", result["KEY"])
	}
}

func TestApply_TransformError_PropagatesError(t *testing.T) {
	expectedErr := errors.New("transform failed")
	tr := transformer.New().Add("fail", func(v string) (string, error) {
		return "", expectedErr
	})
	env := map[string]string{"KEY": "value"}

	_, err := tr.Apply(env)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := transformer.New().Add("upper", transformer.ToUpper())
	env := map[string]string{"KEY": "hello"}

	_, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "hello" {
		t.Errorf("input map was mutated, expected 'hello', got %q", env["KEY"])
	}
}

func TestReplace_Transform(t *testing.T) {
	tr := transformer.New().Add("replace", transformer.Replace("localhost", "prod.example.com"))
	env := map[string]string{"DB_HOST": "localhost:5432"}

	result, err := tr.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "prod.example.com:5432" {
		t.Errorf("unexpected result: %q", result["DB_HOST"])
	}
}
