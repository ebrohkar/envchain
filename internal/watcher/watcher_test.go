package watcher_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/watcher"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envwatch-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0o644)
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempFile(t, "FOO=bar\n")

	w := watcher.New(20*time.Millisecond, path)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() { _ = w.Run(ctx) }()

	// Allow watcher to seed initial state.
	time.Sleep(40 * time.Millisecond)

	// Modify the file.
	_ = os.WriteFile(path, []byte("FOO=baz\n"), 0o644)

	select {
	case ev := <-w.Events():
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempFile(t, "FOO=bar\n")

	w := watcher.New(20*time.Millisecond, path)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go func() { _ = w.Run(ctx) }()

	// No writes — expect no events.
	select {
	case ev := <-w.Events():
		t.Errorf("unexpected event for unchanged file: %v", ev)
	case <-ctx.Done():
		// expected
	}
}

func TestWatcher_MissingFileIsSkipped(t *testing.T) {
	w := watcher.New(20*time.Millisecond, "/nonexistent/path.env")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should not panic or error.
	if err := w.Run(ctx); err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
