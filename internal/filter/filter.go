// Package filter provides functionality for selecting and excluding
// environment variables by name patterns or prefixes.
package filter

import (
	"strings"
)

// Filter holds include and exclude rules for environment variable names.
type Filter struct {
	prefixes []string
	excludes []string
}

// New creates a new Filter with no rules applied.
func New() *Filter {
	return &Filter{}
}

// WithPrefix adds a required prefix. Variables must match at least one
// prefix to be included. If no prefixes are set, all variables pass.
func (f *Filter) WithPrefix(prefixes ...string) *Filter {
	f.prefixes = append(f.prefixes, prefixes...)
	return f
}

// WithExclude adds name patterns to exclude. Variables whose names contain
// any excluded pattern are dropped from results.
func (f *Filter) WithExclude(patterns ...string) *Filter {
	f.excludes = append(f.excludes, patterns...)
	return f
}

// Apply returns only the key-value pairs from env that pass the filter rules.
func (f *Filter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if f.excluded(k) {
			continue
		}
		if !f.included(k) {
			continue
		}
		out[k] = v
	}
	return out
}

// Keys returns only the variable names from env that pass the filter rules.
func (f *Filter) Keys(env map[string]string) []string {
	filtered := f.Apply(env)
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	return keys
}

func (f *Filter) included(name string) bool {
	if len(f.prefixes) == 0 {
		return true
	}
	for _, p := range f.prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func (f *Filter) excluded(name string) bool {
	for _, pattern := range f.excludes {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}
