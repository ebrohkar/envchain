package watcher

import "time"

// Option is a functional option for configuring a Watcher.
type Option func(*Watcher)

// WithInterval overrides the default polling interval.
func WithInterval(d time.Duration) Option {
	return func(w *Watcher) {
		if d > 0 {
			w.interval = d
		}
	}
}

// WithPaths appends additional paths to watch.
func WithPaths(paths ...string) Option {
	return func(w *Watcher) {
		w.paths = append(w.paths, paths...)
	}
}

// NewWithOptions creates a Watcher using functional options.
// The interval defaults to 1 second if not overridden.
func NewWithOptions(opts ...Option) *Watcher {
	w := &Watcher{
		interval: time.Second,
		events:   make(chan Event, 8),
		lastMod:  make(map[string]time.Time),
	}
	for _, o := range opts {
		o(w)
	}
	return w
}
