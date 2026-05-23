// Package profiler provides timing and performance profiling utilities
// for measuring how long environment resolution and validation stages take.
package profiler

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// Entry records the duration of a named stage.
type Entry struct {
	Name     string
	Duration time.Duration
	Started  time.Time
}

// Profiler tracks timing entries for named stages.
type Profiler struct {
	entries []Entry
}

// New returns a new Profiler instance.
func New() *Profiler {
	return &Profiler{}
}

// Track records the duration of fn under the given stage name.
func (p *Profiler) Track(name string, fn func()) {
	start := time.Now()
	fn()
	p.entries = append(p.entries, Entry{
		Name:     name,
		Duration: time.Since(start),
		Started:  start,
	})
}

// Entries returns a copy of all recorded timing entries.
func (p *Profiler) Entries() []Entry {
	out := make([]Entry, len(p.entries))
	copy(out, p.entries)
	return out
}

// Slowest returns the entry with the longest duration, or nil if empty.
func (p *Profiler) Slowest() *Entry {
	if len(p.entries) == 0 {
		return nil
	}
	slowest := p.entries[0]
	for _, e := range p.entries[1:] {
		if e.Duration > slowest.Duration {
			slowest = e
		}
	}
	return &slowest
}

// WriteSummary writes a human-readable timing summary to w.
// If w is nil, os.Stdout is used.
func (p *Profiler) WriteSummary(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	sorted := make([]Entry, len(p.entries))
	copy(sorted, p.entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration > sorted[j].Duration
	})
	fmt.Fprintln(w, "--- profiler summary ---")
	for _, e := range sorted {
		fmt.Fprintf(w, "  %-30s %s\n", e.Name, e.Duration.Round(time.Microsecond))
	}
}
