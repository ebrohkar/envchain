package resolver

import (
	"fmt"
	"strings"
)

// EnvGetter is a function type for retrieving environment variable values.
type EnvGetter func(key string) (string, bool)

// Result holds the outcome of resolving a single variable reference.
type Result struct {
	Key      string
	Resolved string
	Missing  []string
}

// Resolver expands variable references within values using a provided env source.
type Resolver struct {
	getenv EnvGetter
}

// New creates a new Resolver with the given environment getter.
func New(getenv EnvGetter) *Resolver {
	if getenv == nil {
		panic("resolver: getenv must not be nil")
	}
	return &Resolver{getenv: getenv}
}

// Resolve expands all ${VAR} references in value using the resolver's env source.
// Returns the expanded string and a slice of any variable names that were missing.
func (r *Resolver) Resolve(value string) (string, []string) {
	var missing []string
	result := expandVars(value, func(key string) string {
		val, ok := r.getenv(key)
		if !ok || val == "" {
			missing = append(missing, key)
			return ""
		}
		return val
	})
	return result, missing
}

// ResolveAll resolves a map of key→value pairs, returning results for each entry.
func (r *Resolver) ResolveAll(vars map[string]string) []Result {
	results := make([]Result, 0, len(vars))
	for k, v := range vars {
		resolved, missing := r.Resolve(v)
		results = append(results, Result{
			Key:      k,
			Resolved: resolved,
			Missing:  missing,
		})
	}
	return results
}

// expandVars replaces ${VAR} patterns in s using the provided mapping function.
func expandVars(s string, mapping func(string) string) string {
	var sb strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '$' && i+1 < len(s) && s[i+1] == '{' {
			end := strings.Index(s[i+2:], "}")
			if end < 0 {
				sb.WriteString(s[i:])
				break
			}
			key := s[i+2 : i+2+end]
			sb.WriteString(mapping(key))
			i += 2 + end + 1
			continue
		}
		sb.WriteByte(s[i])
		i++
	}
	return sb.String()
}

// HasMissing returns true if any result contains unresolved references.
func HasMissing(results []Result) bool {
	for _, r := range results {
		if len(r.Missing) > 0 {
			return true
		}
	}
	return false
}

// MissingSummary returns a human-readable summary of all unresolved references.
func MissingSummary(results []Result) string {
	var parts []string
	for _, r := range results {
		if len(r.Missing) > 0 {
			parts = append(parts, fmt.Sprintf("%s references missing: %s", r.Key, strings.Join(r.Missing, ", ")))
		}
	}
	return strings.Join(parts, "; ")
}
