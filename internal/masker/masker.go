// Package masker provides utilities for masking sensitive environment
// variable values before display, logging, or export.
package masker

import "strings"

// Strategy defines how a value should be masked.
type Strategy int

const (
	// StrategyFull replaces the entire value with asterisks.
	StrategyFull Strategy = iota
	// StrategyPartial reveals the first and last characters.
	StrategyPartial
	// StrategyPrefix reveals only the first N characters.
	StrategyPrefix
)

// Masker masks environment variable values based on a set of sensitive key patterns.
type Masker struct {
	sensitiveKeys []string
	strategy      Strategy
	prefixLen     int
}

// New returns a Masker that treats the given keys as sensitive.
func New(sensitiveKeys []string, strategy Strategy) *Masker {
	normalized := make([]string, len(sensitiveKeys))
	for i, k := range sensitiveKeys {
		normalized[i] = strings.ToUpper(k)
	}
	return &Masker{
		sensitiveKeys: normalized,
		strategy:      strategy,
		prefixLen:     4,
	}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range m.sensitiveKeys {
		if strings.Contains(upper, s) {
			return true
		}
	}
	return false
}

// Mask returns the masked form of value if the key is sensitive,
// otherwise it returns the value unchanged.
func (m *Masker) Mask(key, value string) string {
	if !m.IsSensitive(key) {
		return value
	}
	if value == "" {
		return ""
	}
	switch m.strategy {
	case StrategyPartial:
		if len(value) <= 2 {
			return "**"
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	case StrategyPrefix:
		n := m.prefixLen
		if len(value) <= n {
			return strings.Repeat("*", len(value))
		}
		return value[:n] + strings.Repeat("*", len(value)-n)
	default:
		return strings.Repeat("*", len(value))
	}
}

// MaskAll applies Mask to every entry in the provided map, returning a new map.
func (m *Masker) MaskAll(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = m.Mask(k, v)
	}
	return out
}
