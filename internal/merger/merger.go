// Package merger provides utilities for merging multiple environment variable
// maps with configurable precedence and conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// StrategyFirst keeps the value from the first map that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last map that defines the key.
	StrategyLast
	// StrategyError returns an error on any key conflict.
	StrategyError
)

// Merger merges multiple environment variable maps.
type Merger struct {
	strategy Strategy
}

// New returns a new Merger with the given conflict resolution strategy.
func New(strategy Strategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines the provided env maps according to the configured strategy.
// Maps are processed in order; later maps have lower priority under StrategyFirst
// and higher priority under StrategyLast.
func (m *Merger) Merge(maps ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, env := range maps {
		for k, v := range env {
			existing, exists := result[k]
			if !exists {
				result[k] = v
				continue
			}

			switch m.strategy {
			case StrategyFirst:
				// keep existing, do nothing
			case StrategyLast:
				result[k] = v
			case StrategyError:
				if existing != v {
					return nil, fmt.Errorf("merger: conflict on key %q: %q vs %q", k, existing, v)
				}
			}
		}
	}

	return result, nil
}

// Keys returns a sorted list of all unique keys across the provided maps.
func Keys(maps ...map[string]string) []string {
	seen := make(map[string]struct{})
	for _, env := range maps {
		for k := range env {
			seen[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
