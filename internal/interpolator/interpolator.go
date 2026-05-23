// Package interpolator provides template-style interpolation of environment
// variable values using ${VAR} or $VAR syntax with support for default values.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

// Env is a function that retrieves an environment variable by name.
type Env func(key string) (string, bool)

// Interpolator replaces variable references in strings with their resolved values.
type Interpolator struct {
	env      Env
	pattern  *regexp.Regexp
}

// New creates a new Interpolator using the provided Env lookup function.
func New(env Env) *Interpolator {
	return &Interpolator{
		env:     env,
		pattern: regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`),
	}
}

// Interpolate replaces all variable references in the input string.
// Supports ${VAR}, ${VAR:-default}, and $VAR syntax.
// Returns an error if a referenced variable is missing and no default is provided.
func (i *Interpolator) Interpolate(input string) (string, error) {
	var firstErr error
	result := i.pattern.ReplaceAllStringFunc(input, func(match string) string {
		if firstErr != nil {
			return match
		}

		key, defaultVal, hasDefault := parseRef(match)

		if val, ok := i.env(key); ok {
			return val
		}
		if hasDefault {
			return defaultVal
		}

		firstErr = fmt.Errorf("interpolator: variable %q is not set", key)
		return match
	})

	if firstErr != nil {
		return "", firstErr
	}
	return result, nil
}

// InterpolateMap applies Interpolate to every value in the provided map.
// Returns a new map and the first error encountered, if any.
func (i *Interpolator) InterpolateMap(vars map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		resolved, err := i.Interpolate(v)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

// parseRef extracts the variable name and optional default value from a match.
func parseRef(match string) (key, defaultVal string, hasDefault bool) {
	// Strip leading $ and optional braces
	inner := strings.TrimPrefix(match, "$")
	inner = strings.TrimPrefix(inner, "{")
	inner = strings.TrimSuffix(inner, "}")

	if idx := strings.Index(inner, ":-"); idx != -1 {
		return inner[:idx], inner[idx+2:], true
	}
	return inner, "", false
}
