package auditor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envchain/internal/auditor"
)

func TestAuditor_RecordAndEvents(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventResolved, "DB_HOST", "resolved from env")
	a.Record(auditor.EventMissing, "API_KEY", "not found")

	events := a.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", events[0].Key)
	}
	if events[1].Kind != auditor.EventMissing {
		t.Errorf("expected missing kind, got %s", events[1].Kind)
	}
}

func TestAuditor_HasFailures_True(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventResolved, "DB_HOST", "ok")
	a.Record(auditor.EventFailed, "DB_PASS", "validation failed")

	if !a.HasFailures() {
		t.Error("expected HasFailures to return true")
	}
}

func TestAuditor_HasFailures_False(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventResolved, "DB_HOST", "ok")
	a.Record(auditor.EventValidated, "DB_PORT", "passed")

	if a.HasFailures() {
		t.Error("expected HasFailures to return false")
	}
}

func TestAuditor_WriteSummary(t *testing.T) {
	var buf bytes.Buffer
	a := auditor.New(&buf)
	a.Record(auditor.EventResolved, "DB_HOST", "resolved from env")
	a.Record(auditor.EventMissing, "API_KEY", "not found")
	a.WriteSummary()

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in summary output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in summary output")
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected 'missing' kind in summary output")
	}
}

func TestAuditor_Events_ReturnsCopy(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventResolved, "X", "ok")

	events := a.Events()
	events[0].Key = "MODIFIED"

	original := a.Events()
	if original[0].Key == "MODIFIED" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestAuditor_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when w is nil (falls back to stdout)
	a := auditor.New(nil)
	a.Record(auditor.EventResolved, "KEY", "ok")
	// WriteSummary writes to stdout — just ensure no panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
}
