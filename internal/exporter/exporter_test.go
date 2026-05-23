package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/exporter"
)

func TestExporter_Shell_Output(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatShell, &buf)

	vars := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}

	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "export APP_ENV='production'") {
		t.Errorf("expected shell export for APP_ENV, got:\n%s", output)
	}
	if !strings.Contains(output, "export DB_HOST='localhost'") {
		t.Errorf("expected shell export for DB_HOST, got:\n%s", output)
	}
}

func TestExporter_Dotenv_Output(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatDotenv, &buf)

	vars := map[string]string{
		"PORT": "8080",
		"HOST": "0.0.0.0",
	}

	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "HOST=0.0.0.0") {
		t.Errorf("expected dotenv line for HOST, got:\n%s", output)
	}
	if !strings.Contains(output, "PORT=8080") {
		t.Errorf("expected dotenv line for PORT, got:\n%s", output)
	}
}

func TestExporter_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatDotenv, &buf)

	vars := map[string]string{"Z_VAR": "z", "A_VAR": "a", "M_VAR": "m"}
	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_VAR") {
		t.Errorf("expected first line to be A_VAR, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_VAR") {
		t.Errorf("expected last line to be Z_VAR, got %q", lines[2])
	}
}

func TestExporter_NilWriter_UsesStdout(t *testing.T) {
	e := exporter.New(exporter.FormatDotenv, nil)
	if e == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestExporter_ShellQuote_SpecialChars(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatShell, &buf)

	vars := map[string]string{"SECRET": "it's a secret"}
	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `'it'\''s a secret'`) {
		t.Errorf("expected escaped single quote in output, got:\n%s", output)
	}
}
