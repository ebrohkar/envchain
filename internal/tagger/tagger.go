package tagger

import "sync"

// Tag represents a key-value metadata pair attached to an env variable.
type Tag struct {
	Key   string
	Value string
}

// Tagger manages tags associated with environment variable keys.
type Tagger struct {
	mu   sync.RWMutex
	tags map[string][]Tag
}

// New returns a new Tagger instance.
func New() *Tagger {
	return &Tagger{
		tags: make(map[string][]Tag),
	}
}

// Tag attaches a key-value tag to the given environment variable name.
func (t *Tagger) Tag(envKey, tagKey, tagValue string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tags[envKey] = append(t.tags[envKey], Tag{Key: tagKey, Value: tagValue})
}

// Get returns all tags associated with the given environment variable name.
// Returns an empty slice if none are found.
func (t *Tagger) Get(envKey string) []Tag {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]Tag, len(t.tags[envKey]))
	copy(result, t.tags[envKey])
	return result
}

// FindByTag returns all env variable names that have a tag matching the given key and value.
func (t *Tagger) FindByTag(tagKey, tagValue string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var matches []string
	for envKey, tags := range t.tags {
		for _, tag := range tags {
			if tag.Key == tagKey && tag.Value == tagValue {
				matches = append(matches, envKey)
				break
			}
		}
	}
	return matches
}

// Keys returns all env variable names that have at least one tag.
func (t *Tagger) Keys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	keys := make([]string, 0, len(t.tags))
	for k := range t.tags {
		keys = append(keys, k)
	}
	return keys
}

// Clear removes all tags for the given environment variable name.
func (t *Tagger) Clear(envKey string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tags, envKey)
}
