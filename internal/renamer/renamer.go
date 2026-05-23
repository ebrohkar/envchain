package renamer

import "fmt"

// Strategy defines how key renaming conflicts are handled.
type Strategy int

const (
	// StrategyOverwrite replaces existing keys in the output.
	StrategyOverwrite Strategy = iota
	// StrategySkip leaves the destination key unchanged if it already exists.
	StrategySkip
	// StrategyError returns an error on conflicting destination keys.
	StrategyError
)

// Renamer renames environment variable keys according to a mapping.
type Renamer struct {
	strategy Strategy
}

// New returns a Renamer with the default StrategyOverwrite.
func New() *Renamer {
	return &Renamer{strategy: StrategyOverwrite}
}

// NewWithStrategy returns a Renamer using the given conflict strategy.
func NewWithStrategy(s Strategy) *Renamer {
	return &Renamer{strategy: s}
}

// Apply renames keys in env according to the provided mapping.
// Keys not present in the mapping are passed through unchanged.
// If a destination key already exists in env the configured strategy applies.
func (r *Renamer) Apply(env map[string]string, mapping map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))

	// Copy all keys that are NOT being renamed.
	renamed := make(map[string]bool, len(mapping))
	for src := range mapping {
		renamed[src] = true
	}
	for k, v := range env {
		if !renamed[k] {
			out[k] = v
		}
	}

	// Apply renames.
	for src, dst := range mapping {
		val, exists := env[src]
		if !exists {
			continue
		}
		if existing, conflict := out[dst]; conflict {
			switch r.strategy {
			case StrategySkip:
				_ = existing
				continue
			case StrategyError:
				return nil, fmt.Errorf("renamer: destination key %q already exists", dst)
			default: // StrategyOverwrite
				out[dst] = val
			}
		} else {
			out[dst] = val
		}
	}

	return out, nil
}

// Keys returns the list of source keys that would be renamed.
func Keys(mapping map[string]string) []string {
	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	return keys
}
