// Package differ provides utilities for comparing two snapshots of environment
// variable maps and surfacing the differences between them.
//
// It is typically used after the watcher detects a file change, allowing
// envchain to report exactly which variables were added, removed, or modified
// since the last known-good state.
//
// Basic usage:
//
//	d := differ.New()
//	diffs := d.Compare(before, after)
//	if d.HasChanges(before, after) {
//		f := differ.NewFormatter(os.Stdout)
//		f.Write(diffs)
//	}
package differ
