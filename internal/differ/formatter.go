package differ

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Formatter renders a slice of Diffs to a writer.
type Formatter struct {
	w io.Writer
}

// NewFormatter returns a Formatter writing to w.
// If w is nil, os.Stdout is used.
func NewFormatter(w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{w: w}
}

// Write outputs a human-readable diff summary.
func (f *Formatter) Write(diffs []Diff) {
	if len(diffs) == 0 {
		fmt.Fprintln(f.w, "No changes detected.")
		return
	}

	sorted := make([]Diff, len(diffs))
	copy(sorted, diffs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, d := range sorted {
		switch d.Kind {
		case Added:
			fmt.Fprintf(f.w, "+ %s=%q\n", d.Key, d.NewValue)
		case Removed:
			fmt.Fprintf(f.w, "- %s=%q\n", d.Key, d.OldValue)
		case Changed:
			fmt.Fprintf(f.w, "~ %s: %q -> %q\n", d.Key, d.OldValue, d.NewValue)
		}
	}
}
