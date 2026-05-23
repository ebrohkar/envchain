// Package prefixer provides utilities for adding, removing, and replacing
// key prefixes across environment variable maps.
package prefixer

import "strings"

// Prefixer applies prefix transformations to environment variable maps.
type Prefixer struct {
	prefix    string
	stripOnly bool
}

// New returns a Prefixer that adds or removes the given prefix.
func New(prefix string) *Prefixer {
	return &Prefixer{prefix: prefix}
}

// Add returns a new map with the prefix prepended to every key.
// Keys that already have the prefix are left unchanged.
func (p *Prefixer) Add(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if strings.HasPrefix(k, p.prefix) {
			out[k] = v
		} else {
			out[p.prefix+k] = v
		}
	}
	return out
}

// Strip returns a new map with the prefix removed from every matching key.
// Keys that do not have the prefix are included unchanged.
func (p *Prefixer) Strip(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if strings.HasPrefix(k, p.prefix) {
			out[strings.TrimPrefix(k, p.prefix)] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// Replace returns a new map where oldPrefix is replaced by newPrefix on
// every key that starts with oldPrefix. Other keys are passed through.
func Replace(env map[string]string, oldPrefix, newPrefix string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if strings.HasPrefix(k, oldPrefix) {
			newKey := newPrefix + strings.TrimPrefix(k, oldPrefix)
			out[newKey] = v
		} else {
			out[k] = v
		}
	}
	return out
}
