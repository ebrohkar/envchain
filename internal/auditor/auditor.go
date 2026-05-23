// Package auditor records and retrieves audit events for environment
// variable resolution and validation runs.
package auditor

import (
	"fmt"
	"io"
	"os"
	"time"
)

// EventKind classifies what happened during an audit event.
type EventKind string

const (
	EventResolved   EventKind = "resolved"
	EventMissing    EventKind = "missing"
	EventValidated  EventKind = "validated"
	EventFailed     EventKind = "failed"
)

// Event represents a single auditable action.
type Event struct {
	Timestamp time.Time
	Kind      EventKind
	Key       string
	Message   string
}

// Auditor collects audit events and can write a summary.
type Auditor struct {
	events []Event
	w      io.Writer
}

// New returns an Auditor that writes output to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Auditor {
	if w == nil {
		w = os.Stdout
	}
	return &Auditor{w: w}
}

// Record appends a new event to the audit log.
func (a *Auditor) Record(kind EventKind, key, message string) {
	a.events = append(a.events, Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Key:       key,
		Message:   message,
	})
}

// Events returns a copy of all recorded events.
func (a *Auditor) Events() []Event {
	out := make([]Event, len(a.events))
	copy(out, a.events)
	return out
}

// WriteSummary writes all recorded events to the auditor's writer.
func (a *Auditor) WriteSummary() {
	for _, e := range a.events {
		fmt.Fprintf(a.w, "[%s] %-10s %-30s %s\n",
			e.Timestamp.Format(time.RFC3339),
			string(e.Kind),
			e.Key,
			e.Message,
		)
	}
}

// HasFailures returns true if any event of kind EventFailed or EventMissing exists.
func (a *Auditor) HasFailures() bool {
	for _, e := range a.events {
		if e.Kind == EventFailed || e.Kind == EventMissing {
			return true
		}
	}
	return false
}
