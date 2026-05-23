package transformer

// SelectiveTransformer applies transforms only to keys matching a predicate.
type SelectiveTransformer struct {
	base      *Transformer
	predicate func(key string) bool
}

// NewSelective creates a SelectiveTransformer that applies transforms only
// to keys for which the predicate returns true.
func NewSelective(predicate func(key string) bool) *SelectiveTransformer {
	return &SelectiveTransformer{
		base:      New(),
		predicate: predicate,
	}
}

// Add registers a named transform function on the selective transformer.
func (s *SelectiveTransformer) Add(name string, fn TransformFunc) *SelectiveTransformer {
	s.base.Add(name, fn)
	return s
}

// Apply runs all registered transforms only on keys matching the predicate.
// Non-matching keys are passed through unchanged.
func (s *SelectiveTransformer) Apply(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))

	matched := make(map[string]string)
	for k, v := range env {
		if s.predicate(k) {
			matched[k] = v
		} else {
			result[k] = v
		}
	}

	transformed, err := s.base.Apply(matched)
	if err != nil {
		return nil, err
	}

	for k, v := range transformed {
		result[k] = v
	}
	return result, nil
}
