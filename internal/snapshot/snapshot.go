// Package snapshot provides functionality for capturing and persisting
// the current state of resolved environment variables to disk, enabling
// comparison between deployments or environment states over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"\n	"time"
)

// Snapshot represents a point-in-time capture of resolved environment variables.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Variables map[string]string `json:"variables"`
}

// New creates a new Snapshot from the given map of resolved variables.
func New(vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Variables: copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &s, nil
}

// Keys returns a sorted list of variable names in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.Variables))
	for k := range s.Variables {
		keys = append(keys, k)
	}
	return keys
}
