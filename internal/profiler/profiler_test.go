package profiler_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envchain/internal/profiler"
)

func TestTrack_RecordsEntry(t *testing.T) {
	p := profiler.New()
	p.Track("load", func() {
		time.Sleep(1 * time.Millisecond)
	})
	entries := p.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "load" {
		t.Errorf("expected name 'load', got %q", entries[0].Name)
	}
	if entries[0].Duration < time.Millisecond {
		t.Errorf("expected duration >= 1ms, got %s", entries[0].Duration)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	p := profiler.New()
	p.Track("resolve", func() {})
	a := p.Entries()
	a[0].Name = "mutated"
	b := p.Entries()
	if b[0].Name == "mutated" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestSlowest_ReturnsLongestEntry(t *testing.T) {
	p := profiler.New()
	p.Track("fast", func() {})
	p.Track("slow", func() { time.Sleep(5 * time.Millisecond) })
	p.Track("medium", func() { time.Sleep(2 * time.Millisecond) })

	slowest := p.Slowest()
	if slowest == nil {
		t.Fatal("expected a slowest entry, got nil")
	}
	if slowest.Name != "slow" {
		t.Errorf("expected slowest to be 'slow', got %q", slowest.Name)
	}
}

func TestSlowest_EmptyProfiler_ReturnsNil(t *testing.T) {
	p := profiler.New()
	if p.Slowest() != nil {
		t.Error("expected nil for empty profiler")
	}
}

func TestWriteSummary_ContainsStageNames(t *testing.T) {
	p := profiler.New()
	p.Track("validate", func() {})
	p.Track("export", func() {})

	var buf bytes.Buffer
	p.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "validate") {
		t.Error("summary should contain 'validate'")
	}
	if !strings.Contains(out, "export") {
		t.Error("summary should contain 'export'")
	}
	if !strings.Contains(out, "profiler summary") {
		t.Error("summary should contain header")
	}
}
