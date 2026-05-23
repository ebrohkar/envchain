// Package scoper provides namespace-scoped views of environment variable maps,
// allowing consumers to work with a subset of variables under a given prefix.
package scoper

import "strings"

// Scoper partitions environment variables by namespace prefix.
type Scoper struct {
	env    map[string]string
	sep    string
}

// New creates a Scoper over the provided environment map.
// The separator defaults to "_" if an empty string is given.
func New(env map[string]string, sep string) *Scoper {
	if sep == "" {
		sep = "_"
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Scoper{env: copy, sep: sep}
}

// Scope returns all variables whose keys begin with the given namespace prefix
// (case-insensitive). The returned map keys have the prefix and separator stripped.
func (s *Scoper) Scope(namespace string) map[string]string {
	prefix := strings.ToUpper(namespace) + s.sep
	result := make(map[string]string)
	for k, v := range s.env {
		upper := strings.ToUpper(k)
		if strings.HasPrefix(upper, prefix) {
			stripped := k[len(prefix):]
			if stripped != "" {
				result[stripped] = v
			}
		}
	}
	return result
}

// Namespaces returns the distinct top-level namespace prefixes present in the
// environment, determined by splitting each key on the separator.
func (s *Scoper) Namespaces() []string {
	seen := make(map[string]struct{})
	for k := range s.env {
		parts := strings.SplitN(k, s.sep, 2)
		if len(parts) == 2 && parts[1] != "" {
			seen[parts[0]] = struct{}{}
		}
	}
	ns := make([]string, 0, len(seen))
	for k := range seen {
		ns = append(ns, k)
	}
	return ns
}

// Has reports whether the given key exists within the specified namespace scope.
func (s *Scoper) Has(namespace, key string) bool {
	scoped := s.Scope(namespace)
	_, ok := scoped[key]
	return ok
}
