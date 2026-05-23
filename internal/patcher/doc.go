// Package patcher provides utilities for merging environment variable maps
// using configurable conflict resolution strategies.
//
// Three strategies are supported:
//
//   - StrategyOverwrite: patch values always win, replacing existing keys.
//   - StrategySkipExisting: existing keys are preserved; only new keys are added.
//   - StrategyErrorOnConflict: returns an error if a key exists with a different value.
//
// Example usage:
//
//	p := patcher.New(patcher.StrategyOverwrite)
//	result, err := p.Apply(base, patch)
//
// The Apply method never mutates its input maps.
package patcher
