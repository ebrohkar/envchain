package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/snapshot"
)

func TestNew_CopiesVariables(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New(vars)

	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if s.Variables["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", s.Variables["FOO"])
	}
	// Mutating original should not affect snapshot
	vars["FOO"] = "mutated"
	if s.Variables["FOO"] != "bar" {
		t.Error("snapshot should be independent of original map")
	}
}

func TestNew_TimestampSet(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New(map[string]string{})
	after := time.Now().UTC()

	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range", s.Timestamp)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "8080"}
	s := snapshot.New(vars)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := s.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Variables["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", loaded.Variables["HOST"])
	}
	if loaded.Variables["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", loaded.Variables["PORT"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := snapshot.New(vars)
	keys := s.Keys()

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
