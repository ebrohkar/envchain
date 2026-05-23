package transformer

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms an environment variable value.
type TransformFunc func(value string) (string, error)

// Transformer applies a chain of transformations to environment variable values.
type Transformer struct {
	transforms []namedTransform
}

type namedTransform struct {
	name string
	fn   TransformFunc
}

// New creates a new Transformer with no transforms applied.
func New() *Transformer {
	return &Transformer{}
}

// Add registers a named transform function to the transformer pipeline.
func (t *Transformer) Add(name string, fn TransformFunc) *Transformer {
	t.transforms = append(t.transforms, namedTransform{name: name, fn: fn})
	return t
}

// Apply runs all registered transforms on the given value map and returns
// a new map with transformed values. Returns an error if any transform fails.
func (t *Transformer) Apply(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	for _, tr := range t.transforms {
		for k, v := range result {
			transformed, err := tr.fn(v)
			if err != nil {
				return nil, fmt.Errorf("transform %q failed on key %q: %w", tr.name, k, err)
			}
			result[k] = transformed
		}
	}
	return result, nil
}

// TrimSpace returns a TransformFunc that trims leading and trailing whitespace.
func TrimSpace() TransformFunc {
	return func(value string) (string, error) {
		return strings.TrimSpace(value), nil
	}
}

// ToUpper returns a TransformFunc that converts values to uppercase.
func ToUpper() TransformFunc {
	return func(value string) (string, error) {
		return strings.ToUpper(value), nil
	}
}

// ToLower returns a TransformFunc that converts values to lowercase.
func ToLower() TransformFunc {
	return func(value string) (string, error) {
		return strings.ToLower(value), nil
	}
}

// Replace returns a TransformFunc that replaces all occurrences of old with new.
func Replace(old, new string) TransformFunc {
	return func(value string) (string, error) {
		return strings.ReplaceAll(value, old, new), nil
	}
}
