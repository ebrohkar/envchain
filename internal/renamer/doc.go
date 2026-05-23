// Package renamer provides utilities for renaming environment variable keys
// according to a user-supplied mapping.
//
// A Renamer accepts a source-to-destination mapping and applies it to an
// environment map, passing through any keys not present in the mapping
// unchanged. Three conflict strategies control what happens when the
// destination key already exists in the environment:
//
//   - StrategyOverwrite (default): the destination is replaced with the
//     renamed value.
//   - StrategySkip: the existing destination value is preserved and the
//     rename is silently ignored.
//   - StrategyError: an error is returned immediately, leaving the caller
//     responsible for resolving the conflict.
//
// Example usage:
//
//	r := renamer.New()
//	out, err := r.Apply(env, map[string]string{"OLD_HOST": "APP_HOST"})
package renamer
