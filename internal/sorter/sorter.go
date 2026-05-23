// Package sorter provides utilities for ordering environment variable maps
// by key name, value length, or custom priority lists before export or display.
package sorter

import (
	"sort"
)

// Strategy defines how keys are ordered.
type Strategy int

const (
	// ByKey sorts keys alphabetically (default).
	ByKey Strategy = iota
	// ByValueLength sorts keys by the length of their value, ascending.
	ByValueLength
	// ByPriority sorts keys by a user-defined priority list; unmatched keys follow alphabetically.
	ByPriority
)

// Sorter orders the keys of an environment map.
type Sorter struct {
	strategy Strategy
	priority []string
}

// New returns a Sorter using the given strategy.
func New(s Strategy) *Sorter {
	return &Sorter{strategy: s}
}

// WithPriority sets an explicit key ordering for ByPriority strategy.
// Keys not in the list are appended alphabetically after the priority keys.
func (s *Sorter) WithPriority(keys []string) *Sorter {
	s.priority = keys
	return s
}

// Sort returns the keys of env in the order defined by the Sorter's strategy.
func (s *Sorter) Sort(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch s.strategy {
	case ByValueLength:
		sort.Slice(keys, func(i, j int) bool {
			li, lj := len(env[keys[i]]), len(env[keys[j]])
			if li != lj {
				return li < lj
			}
			return keys[i] < keys[j]
		})
	case ByPriority:
		rank := make(map[string]int, len(s.priority))
		for i, k := range s.priority {
			rank[k] = i
		}
		sort.Slice(keys, func(i, j int) bool {
			ri, iOK := rank[keys[i]]
			rj, jOK := rank[keys[j]]
			switch {
			case iOK && jOK:
				return ri < rj
			case iOK:
				return true
			case jOK:
				return false
			default:
				return keys[i] < keys[j]
			}
		})
	default: // ByKey
		sort.Strings(keys)
	}

	return keys
}
