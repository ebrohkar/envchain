// Package watcher monitors environment variable files for changes
// and triggers re-evaluation of the pipeline when modifications are detected.
package watcher

import (
	"context"
	"os"
	"time"
)

// Event represents a file change notification.
type Event struct {
	Path    string
	ModTime time.Time
}

// Watcher polls one or more env files for modifications.
type Watcher struct {
	paths    []string
	interval time.Duration
	events   chan Event
	lastMod  map[string]time.Time
}

// New creates a Watcher that polls the given paths at the specified interval.
func New(interval time.Duration, paths ...string) *Watcher {
	return &Watcher{
		paths:    paths,
		interval: interval,
		events:   make(chan Event, len(paths)),
		lastMod:  make(map[string]time.Time),
	}
}

// Events returns the read-only channel of change events.
func (w *Watcher) Events() <-chan Event {
	return w.events
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	// Seed initial modification times.
	for _, p := range w.paths {
		if info, err := os.Stat(p); err == nil {
			w.lastMod[p] = info.ModTime()
		}
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.poll()
		}
	}
}

// poll checks each watched path for modification time changes.
func (w *Watcher) poll() {
	for _, p := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if prev, ok := w.lastMod[p]; !ok || info.ModTime().After(prev) {
			w.lastMod[p] = info.ModTime()
			select {
			case w.events <- Event{Path: p, ModTime: info.ModTime()}:
			default:
			}
		}
	}
}
