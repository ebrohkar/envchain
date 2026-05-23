// Package comparator provides utilities for comparing two sets of environment
// variables and producing structured diff results with change classification.
package comparator

import "sort"

// ChangeType describes the kind of change detected for a variable.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
	ChangeUnchanged ChangeType = "unchanged"
)

// Change represents a single variable comparison result.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// Result holds the full comparison output between two env maps.
type Result struct {
	Changes []Change
}

// HasDrift returns true if any change is not of type unchanged.
func (r *Result) HasDrift() bool {
	for _, c := range r.Changes {
		if c.Type != ChangeUnchanged {
			return true
		}
	}
	return false
}

// ByType filters and returns changes matching the given ChangeType.
func (r *Result) ByType(ct ChangeType) []Change {
	var out []Change
	for _, c := range r.Changes {
		if c.Type == ct {
			out = append(out, c)
		}
	}
	return out
}

// Comparator compares two environment variable maps.
type Comparator struct{}

// New returns a new Comparator instance.
func New() *Comparator {
	return &Comparator{}
}

// Compare performs a full comparison between baseline and current env maps.
func (c *Comparator) Compare(baseline, current map[string]string) *Result {
	seen := make(map[string]bool)
	var changes []Change

	for key, oldVal := range baseline {
		seen[key] = true
		newVal, exists := current[key]
		if !exists {
			changes = append(changes, Change{Key: key, OldValue: oldVal, NewValue: "", Type: ChangeRemoved})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: key, OldValue: oldVal, NewValue: newVal, Type: ChangeUpdated})
		} else {
			changes = append(changes, Change{Key: key, OldValue: oldVal, NewValue: newVal, Type: ChangeUnchanged})
		}
	}

	for key, newVal := range current {
		if !seen[key] {
			changes = append(changes, Change{Key: key, OldValue: "", NewValue: newVal, Type: ChangeAdded})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return &Result{Changes: changes}
}
