// Package differ compares two sets of environment variable snapshots
// and reports additions, removals, and value changes between them.
package differ

// ChangeKind describes the type of change detected for a variable.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Diff represents a single detected change between two env snapshots.
type Diff struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Differ compares environment variable maps.
type Differ struct{}

// New returns a new Differ instance.
func New() *Differ {
	return &Differ{}
}

// Compare takes two env maps (before, after) and returns a slice of Diffs.
func (d *Differ) Compare(before, after map[string]string) []Diff {
	var diffs []Diff

	for k, oldVal := range before {
		newVal, exists := after[k]
		if !exists {
			diffs = append(diffs, Diff{Key: k, Kind: Removed, OldValue: oldVal})
		} else if oldVal != newVal {
			diffs = append(diffs, Diff{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, newVal := range after {
		if _, exists := before[k]; !exists {
			diffs = append(diffs, Diff{Key: k, Kind: Added, NewValue: newVal})
		}
	}

	return diffs
}

// HasChanges returns true if any diffs exist between before and after.
func (d *Differ) HasChanges(before, after map[string]string) bool {
	return len(d.Compare(before, after)) > 0
}
