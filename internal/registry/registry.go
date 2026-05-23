package registry

import (
	"fmt"
	"sync"
)

// Entry holds metadata about a registered environment variable.
type Entry struct {
	Key         string
	Description string
	Required    bool
	Default     string
}

// Registry stores known environment variable definitions for a service.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		entries: make(map[string]Entry),
	}
}

// Register adds or replaces an entry in the registry.
func (r *Registry) Register(e Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[e.Key] = e
}

// Lookup returns the Entry for the given key, or an error if not found.
func (r *Registry) Lookup(key string) (Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[key]
	if !ok {
		return Entry{}, fmt.Errorf("registry: key %q not registered", key)
	}
	return e, nil
}

// Keys returns all registered keys in the registry.
func (r *Registry) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.entries))
	for k := range r.entries {
		keys = append(keys, k)
	}
	return keys
}

// RequiredKeys returns only the keys marked as required.
func (r *Registry) RequiredKeys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var keys []string
	for k, e := range r.entries {
		if e.Required {
			keys = append(keys, k)
		}
	}
	return keys
}

// Size returns the number of registered entries.
func (r *Registry) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
