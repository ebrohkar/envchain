package validator_test

import (
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/validator"
)

func mockEnv(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestValidator_AllPassed(t *testing.T) {
	env := mockEnv(map[string]string{"HOST": "localhost", "PORT": "8080"})
	c := chain.New([]string{"HOST", "PORT"}, env)
	v := validator.New(c)
	result := v.Validate()

	if !result.IsValid() {
		t.Fatalf("expected valid result, got failures: %v", result.Failed)
	}
	if len(result.Passed) != 2 {
		t.Errorf("expected 2 passed, got %d", len(result.Passed))
	}
}

func TestValidator_SomeFailed(t *testing.T) {
	env := mockEnv(map[string]string{"HOST": "localhost"})
	c := chain.New([]string{"HOST", "PORT"}, env)
	v := validator.New(c)
	result := v.Validate()

	if result.IsValid() {
		t.Fatal("expected invalid result")
	}
	if len(result.Failed) != 1 || result.Failed[0] != "PORT" {
		t.Errorf("expected PORT in failed, got %v", result.Failed)
	}
}

func TestValidator_EmptyVariable(t *testing.T) {
	env := mockEnv(map[string]string{"HOST": "", "PORT": "8080"})
	c := chain.New([]string{"HOST", "PORT"}, env)
	v := validator.New(c)
	result := v.Validate()

	if result.IsValid() {
		t.Fatal("expected invalid result due to empty HOST")
	}
}

func TestResult_Summary(t *testing.T) {
	r := &validator.Result{
		Passed:   []string{"A", "B"},
		Failed:   []string{"C"},
		Warnings: []string{},
	}
	got := r.Summary()
	want := "2 passed, 1 failed, 0 warnings"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
