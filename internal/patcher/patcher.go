package patcher

import "fmt"

// Strategy defines how patch operations are applied.
type Strategy int

const (
	// StrategyOverwrite replaces existing values unconditionally.
	StrategyOverwrite Strategy = iota
	// StrategySkipExisting leaves existing keys untouched.
	StrategySkipExisting
	// StrategyErrorOnConflict returns an error if a key already exists with a different value.
	StrategyErrorOnConflict
)

// Patcher applies a patch map onto a base environment map.
type Patcher struct {
	strategy Strategy
}

// New returns a new Patcher with the given strategy.
func New(strategy Strategy) *Patcher {
	return &Patcher{strategy: strategy}
}

// Apply merges patch into base according to the configured strategy.
// It returns a new map and does not mutate the inputs.
func (p *Patcher) Apply(base, patch map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, v := range patch {
		existing, exists := result[k]
		switch p.strategy {
		case StrategyOverwrite:
			result[k] = v
		case StrategySkipExisting:
			if !exists {
				result[k] = v
			}
		case StrategyErrorOnConflict:
			if exists && existing != v {
				return nil, fmt.Errorf("patcher: conflict on key %q: existing=%q patch=%q", k, existing, v)
			}
			result[k] = v
		}
	}

	return result, nil
}

// Keys returns the sorted list of keys that differ between base and patch.
func Keys(base, patch map[string]string) []string {
	seen := make(map[string]struct{})
	var changed []string
	for k, v := range patch {
		if bv, ok := base[k]; !ok || bv != v {
			if _, already := seen[k]; !already {
				changed = append(changed, k)
				seen[k] = struct{}{}
			}
		}
	}
	return changed
}
