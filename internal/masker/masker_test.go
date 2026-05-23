package masker_test

import (
	"testing"

	"github.com/envchain/envchain/internal/masker"
)

func TestIsSensitive_MatchesSubstring(t *testing.T) {
	m := masker.New([]string{"SECRET", "PASSWORD"}, masker.StrategyFull)
	if !m.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
	if !m.IsSensitive("API_SECRET_KEY") {
		t.Error("expected API_SECRET_KEY to be sensitive")
	}
	if m.IsSensitive("DB_HOST") {
		t.Error("expected DB_HOST to not be sensitive")
	}
}

func TestMask_FullStrategy(t *testing.T) {
	m := masker.New([]string{"SECRET"}, masker.StrategyFull)
	got := m.Mask("API_SECRET", "mysecretvalue")
	want := "*************"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_PartialStrategy(t *testing.T) {
	m := masker.New([]string{"TOKEN"}, masker.StrategyPartial)
	got := m.Mask("AUTH_TOKEN", "abcdef")
	want := "a****f"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_PartialStrategy_ShortValue(t *testing.T) {
	m := masker.New([]string{"TOKEN"}, masker.StrategyPartial)
	got := m.Mask("AUTH_TOKEN", "ab")
	if got != "**" {
		t.Errorf("got %q, want **", got)
	}
}

func TestMask_PrefixStrategy(t *testing.T) {
	m := masker.New([]string{"KEY"}, masker.StrategyPrefix)
	got := m.Mask("API_KEY", "abcdefghij")
	want := "abcd******"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_NonSensitiveKey_Unchanged(t *testing.T) {
	m := masker.New([]string{"SECRET"}, masker.StrategyFull)
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	m := masker.New([]string{"SECRET"}, masker.StrategyFull)
	got := m.Mask("API_SECRET", "")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMaskAll_MasksOnlySensitive(t *testing.T) {
	m := masker.New([]string{"SECRET", "PASSWORD"}, masker.StrategyFull)
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "hunter2",
		"API_SECRET":  "tok123",
	}
	result := m.MaskAll(env)
	if result["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should be unchanged, got %q", result["DB_HOST"])
	}
	if result["DB_PASSWORD"] != "*******" {
		t.Errorf("DB_PASSWORD should be masked, got %q", result["DB_PASSWORD"])
	}
	if result["API_SECRET"] != "******" {
		t.Errorf("API_SECRET should be masked, got %q", result["API_SECRET"])
	}
}
