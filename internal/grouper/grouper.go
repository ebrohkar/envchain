// Package grouper provides functionality for grouping environment variables
// by a shared key pattern, prefix, or custom classifier function.
package grouper

import "sort"

// GroupFn is a function that maps an environment variable key to a group name.
// If the key does not belong to any group, it should return an empty string.
type GroupFn func(key string) string

// Grouper holds the classification function and operates on env maps.
type Grouper struct {
	fn GroupFn
}

// New creates a new Grouper using the provided GroupFn.
func New(fn GroupFn) *Grouper {
	return &Grouper{fn: fn}
}

// Group partitions the given env map into named groups.
// Keys that map to an empty group name are placed under the "" (ungrouped) key.
func (g *Grouper) Group(env map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for k, v := range env {
		groupName := g.fn(k)
		if result[groupName] == nil {
			result[groupName] = make(map[string]string)
		}
		result[groupName][k] = v
	}
	return result
}

// GroupNames returns a sorted list of group names present in the grouped result.
func GroupNames(grouped map[string]map[string]string) []string {
	names := make([]string, 0, len(grouped))
	for name := range grouped {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ByPrefix returns a GroupFn that groups keys by their prefix up to the
// first occurrence of sep. If sep is not found, the key is ungrouped ("").
func ByPrefix(sep string) GroupFn {
	return func(key string) string {
		for i := 0; i < len(key); i++ {
			if string(key[i]) == sep {
				return key[:i]
			}
		}
		return ""
	}
}
