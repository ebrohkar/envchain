package filter_test

import (
	"sort"
	"testing"

	"github.com/user/envchain/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"INTERNAL_KEY": "hidden",
	"LOG_LEVEL":    "info",
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestFilter_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New()
	result := f.Apply(sampleEnv)
	if len(result) != len(sampleEnv) {
		t.Errorf("expected %d entries, got %d", len(sampleEnv), len(result))
	}
}

func TestFilter_WithPrefix_FiltersCorrectly(t *testing.T) {
	f := filter.New().WithPrefix("APP_")
	result := f.Apply(sampleEnv)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestFilter_WithMultiplePrefixes(t *testing.T) {
	f := filter.New().WithPrefix("APP_", "DB_")
	result := f.Apply(sampleEnv)
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
}

func TestFilter_WithExclude_DropsMatches(t *testing.T) {
	f := filter.New().WithExclude("PASSWORD", "KEY")
	result := f.Apply(sampleEnv)
	keys := sortedKeys(result)
	for _, k := range keys {
		if k == "DB_PASSWORD" || k == "INTERNAL_KEY" {
			t.Errorf("key %q should have been excluded", k)
		}
	}
	if len(result) != 4 {
		t.Errorf("expected 4 entries after exclusion, got %d", len(result))
	}
}

func TestFilter_PrefixAndExclude_Combined(t *testing.T) {
	f := filter.New().WithPrefix("DB_").WithExclude("PASSWORD")
	result := f.Apply(sampleEnv)
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilter_Keys_ReturnsNames(t *testing.T) {
	f := filter.New().WithPrefix("LOG_")
	keys := f.Keys(sampleEnv)
	if len(keys) != 1 || keys[0] != "LOG_LEVEL" {
		t.Errorf("expected [LOG_LEVEL], got %v", keys)
	}
}

func TestFilter_EmptyEnv_ReturnsEmpty(t *testing.T) {
	f := filter.New().WithPrefix("APP_")
	result := f.Apply(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
