// Package sanitizer provides utilities for sanitizing environment variable
// keys and values before they are used in deployment pipelines.
package sanitizer

import (
	"regexp"
	"strings"
)

// InvalidKeyChars matches characters not allowed in env var keys.
var invalidKeyChars = regexp.MustCompile(`[^A-Z0-9_]`)

// Option configures the Sanitizer.
type Option func(*Sanitizer)

// Sanitizer cleans environment variable keys and values.
type Sanitizer struct {
	stripControlChars bool
	normalizeKeys     bool
}

// New returns a new Sanitizer with the given options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{
		stripControlChars: true,
		normalizeKeys:     true,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithStripControlChars controls whether control characters are removed from values.
func WithStripControlChars(enabled bool) Option {
	return func(s *Sanitizer) {
		s.stripControlChars = enabled
	}
}

// WithNormalizeKeys controls whether keys are uppercased and invalid chars replaced.
func WithNormalizeKeys(enabled bool) Option {
	return func(s *Sanitizer) {
		s.normalizeKeys = enabled
	}
}

// SanitizeKey cleans an environment variable key.
// It uppercases the key and replaces invalid characters with underscores.
func (s *Sanitizer) SanitizeKey(key string) string {
	if !s.normalizeKeys {
		return key
	}
	upper := strings.ToUpper(strings.TrimSpace(key))
	return invalidKeyChars.ReplaceAllString(upper, "_")
}

// SanitizeValue cleans an environment variable value.
// It trims surrounding whitespace and optionally strips control characters.
func (s *Sanitizer) SanitizeValue(value string) string {
	v := strings.TrimSpace(value)
	if s.stripControlChars {
		v = stripControl(v)
	}
	return v
}

// SanitizeMap applies SanitizeKey and SanitizeValue to all entries in a map.
// If two keys normalize to the same value, the last one wins.
func (s *Sanitizer) SanitizeMap(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[s.SanitizeKey(k)] = s.SanitizeValue(v)
	}
	return result
}

// stripControl removes ASCII control characters (0x00–0x1F, 0x7F) from s.
func stripControl(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x20 && r != 0x7F {
			b.WriteRune(r)
		}
	}
	return b.String()
}
