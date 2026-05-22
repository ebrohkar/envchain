package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestFromFile_ValidEntries(t *testing.T) {
	path := writeTempEnvFile(t, "# comment\nDB_HOST=localhost\nDB_PORT=5432\n")
	env, err := loader.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", env["DB_HOST"])
	}
	if env["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", env["DB_PORT"])
	}
}

func TestFromFile_EmptyAndCommentLines(t *testing.T) {
	path := writeTempEnvFile(t, "\n# just a comment\n\nKEY=value\n")
	env, err := loader.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}

func TestFromFile_InvalidFormat(t *testing.T) {
	path := writeTempEnvFile(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := loader.FromFile(path)
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestFromFile_NotFound(t *testing.T) {
	_, err := loader.FromFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFromEnv_ReadsProcessEnv(t *testing.T) {
	t.Setenv("TEST_ENVCHAIN_KEY", "hello")
	env := loader.FromEnv([]string{"TEST_ENVCHAIN_KEY", "TEST_ENVCHAIN_MISSING"})
	if env["TEST_ENVCHAIN_KEY"] != "hello" {
		t.Errorf("expected hello, got %q", env["TEST_ENVCHAIN_KEY"])
	}
	if env["TEST_ENVCHAIN_MISSING"] != "" {
		t.Errorf("expected empty string for missing key, got %q", env["TEST_ENVCHAIN_MISSING"])
	}
}
