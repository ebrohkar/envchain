// Package redactor provides utilities for redacting sensitive environment
// variable values from maps before logging, reporting, or exporting.
package redactor

import "strings"

// Strategy defines how a value is redacted.
type Strategy int

const (
	// StrategyFull replaces the entire value with a fixed placeholder.
	StrategyFull Strategy = iota
	// StrategyPartial keeps the first N characters and masks the rest.
	StrategyPartial
)

const defaultPlaceholder = "[REDACTED]"

// Redactor filters sensitive keys from an environment map.
type Redactor struct {
	keys        []string
	strategy    Strategy
	visibleRune int
	placeholder string
}

// New returns a Redactor that redacts the given key names (case-insensitive).
func New(sensitiveKeys []string) *Redactor {
	normalized := make([]string, len(sensitiveKeys))
	for i, k := range sensitiveKeys {
		normalized[i] = strings.ToUpper(k)
	}
	return &Redactor{
		keys:        normalized,
		strategy:    StrategyFull,
		visibleRune: 4,
		placeholder: defaultPlaceholder,
	}
}

// WithStrategy sets the redaction strategy.
func (r *Redactor) WithStrategy(s Strategy) *Redactor {
	r.strategy = s
	return r
}

// WithPlaceholder sets a custom placeholder string for StrategyFull.
func (r *Redactor) WithPlaceholder(p string) *Redactor {
	r.placeholder = p
	return r
}

// Apply returns a copy of env with sensitive values redacted.
func (r *Redactor) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.isSensitive(k) {
			out[k] = r.redact(v)
		} else {
			out[k] = v
		}
	}
	return out
}

func (r *Redactor) isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, k := range r.keys {
		if k == upper {
			return true
		}
	}
	return false
}

func (r *Redactor) redact(value string) string {
	if r.strategy == StrategyPartial && len(value) > r.visibleRune {
		return value[:r.visibleRune] + strings.Repeat("*", len(value)-r.visibleRune)
	}
	return r.placeholder
}
