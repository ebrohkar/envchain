package registry

import "fmt"

// MissingResult describes a required key absent from the provided environment.
type MissingResult struct {
	Key         string
	Description string
}

// CheckResult holds the outcome of a registry check against an environment.
type CheckResult struct {
	Missing  []MissingResult
	Defaults map[string]string // keys resolved via their registered default
}

// AllPresent returns true when no required keys are missing.
func (c CheckResult) AllPresent() bool {
	return len(c.Missing) == 0
}

// Summary returns a human-readable summary string.
func (c CheckResult) Summary() string {
	if c.AllPresent() {
		return "all required variables present"
	}
	return fmt.Sprintf("%d required variable(s) missing", len(c.Missing))
}

// Check validates that all required entries in r are present in env.
// If a required key is absent but has a non-empty Default, it is recorded
// in CheckResult.Defaults rather than Missing.
func (r *Registry) Check(env func(string) (string, bool)) CheckResult {
	result := CheckResult{
		Defaults: make(map[string]string),
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for key, entry := range r.entries {
		if !entry.Required {
			continue
		}
		val, ok := env(key)
		if ok && val != "" {
			continue
		}
		if entry.Default != "" {
			result.Defaults[key] = entry.Default
			continue
		}
		result.Missing = append(result.Missing, MissingResult{
			Key:         key,
			Description: entry.Description,
		})
	}
	return result
}
