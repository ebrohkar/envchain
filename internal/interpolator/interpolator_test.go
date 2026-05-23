package interpolator_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/interpolator"
)

func mockEnv(vars map[string]string) interpolator.Env {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	interp := interpolator.New(mockEnv(nil))
	got, err := interp.Interpolate("plain value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "plain value" {
		t.Errorf("expected %q, got %q", "plain value", got)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	interp := interpolator.New(mockEnv(map[string]string{"HOST": "example.com"}))
	got, err := interp.Interpolate("https://${HOST}/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.com/api" {
		t.Errorf("expected %q, got %q", "https://example.com/api", got)
	}
}

func TestInterpolate_ShorthandStyle(t *testing.T) {
	interp := interpolator.New(mockEnv(map[string]string{"PORT": "8080"}))
	got, err := interp.Interpolate("port=$PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "port=8080" {
		t.Errorf("expected %q, got %q", "port=8080", got)
	}
}

func TestInterpolate_DefaultValue_UsedWhenMissing(t *testing.T) {
	interp := interpolator.New(mockEnv(nil))
	got, err := interp.Interpolate("${TIMEOUT:-30s}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "30s" {
		t.Errorf("expected %q, got %q", "30s", got)
	}
}

func TestInterpolate_DefaultValue_IgnoredWhenSet(t *testing.T) {
	interp := interpolator.New(mockEnv(map[string]string{"TIMEOUT": "60s"}))
	got, err := interp.Interpolate("${TIMEOUT:-30s}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "60s" {
		t.Errorf("expected %q, got %q", "60s", got)
	}
}

func TestInterpolate_MissingVariable_ReturnsError(t *testing.T) {
	interp := interpolator.New(mockEnv(nil))
	_, err := interp.Interpolate("${MISSING_VAR}")
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolateMap_AllResolved(t *testing.T) {
	env := mockEnv(map[string]string{"BASE_URL": "https://api.example.com"})
	interp := interpolator.New(env)
	input := map[string]string{
		"ENDPOINT": "${BASE_URL}/v1/users",
		"HEALTH":   "${BASE_URL}/health",
	}
	out, err := interp.InterpolateMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENDPOINT"] != "https://api.example.com/v1/users" {
		t.Errorf("unexpected ENDPOINT: %q", out["ENDPOINT"])
	}
	if out["HEALTH"] != "https://api.example.com/health" {
		t.Errorf("unexpected HEALTH: %q", out["HEALTH"])
	}
}

func TestInterpolateMap_MissingVariable_ReturnsError(t *testing.T) {
	interp := interpolator.New(mockEnv(nil))
	_, err := interp.InterpolateMap(map[string]string{"KEY": "${UNDEFINED}"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
