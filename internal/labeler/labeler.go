// Package labeler provides functionality for attaching and querying
// metadata labels on environment variable maps.
package labeler

import "fmt"

// Label represents a key-value metadata tag.
type Label struct {
	Key   string
	Value string
}

// Labeler manages labels associated with environment variable keys.
type Labeler struct {
	labels map[string][]Label
}

// New creates a new Labeler instance.
func New() *Labeler {
	return &Labeler{
		labels: make(map[string][]Label),
	}
}

// Attach adds a label to the given environment variable key.
func (l *Labeler) Attach(envKey, labelKey, labelValue string) {
	l.labels[envKey] = append(l.labels[envKey], Label{Key: labelKey, Value: labelValue})
}

// Get returns all labels for the given environment variable key.
func (l *Labeler) Get(envKey string) []Label {
	copy := make([]Label, len(l.labels[envKey]))
	for i, lbl := range l.labels[envKey] {
		copy[i] = lbl
	}
	return copy
}

// FindByLabel returns all environment variable keys that have a label
// matching the given key and value.
func (l *Labeler) FindByLabel(labelKey, labelValue string) []string {
	var result []string
	for envKey, labels := range l.labels {
		for _, lbl := range labels {
			if lbl.Key == labelKey && lbl.Value == labelValue {
				result = append(result, envKey)
				break
			}
		}
	}
	return result
}

// Remove deletes all labels for the given environment variable key.
func (l *Labeler) Remove(envKey string) {
	delete(l.labels, envKey)
}

// Summary returns a human-readable string of all labels for an env key.
func (l *Labeler) Summary(envKey string) string {
	labels := l.labels[envKey]
	if len(labels) == 0 {
		return fmt.Sprintf("%s: (no labels)", envKey)
	}
	summary := fmt.Sprintf("%s:", envKey)
	for _, lbl := range labels {
		summary += fmt.Sprintf(" [%s=%s]", lbl.Key, lbl.Value)
	}
	return summary
}
