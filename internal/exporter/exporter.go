// Package exporter provides functionality for exporting resolved environment
// variables to various output formats such as shell scripts or dotenv files.
package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the export output format.
type Format string

const (
	// FormatShell exports variables as shell export statements.
	FormatShell Format = "shell"
	// FormatDotenv exports variables in dotenv file format.
	FormatDotenv Format = "dotenv"
)

// Exporter writes resolved environment variables to an output destination.
type Exporter struct {
	format Format
	writer io.Writer
}

// New creates a new Exporter with the given format and writer.
// If writer is nil, os.Stdout is used.
func New(format Format, writer io.Writer) *Exporter {
	if writer == nil {
		writer = os.Stdout
	}
	return &Exporter{format: format, writer: writer}
}

// Export writes the provided environment variables map to the configured output.
// Keys are sorted alphabetically for deterministic output.
func (e *Exporter) Export(vars map[string]string) error {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line, err := e.formatLine(k, vars[k])
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(e.writer, line); err != nil {
			return fmt.Errorf("exporter: write error for key %q: %w", k, err)
		}
	}
	return nil
}

func (e *Exporter) formatLine(key, value string) (string, error) {
	switch e.format {
	case FormatShell:
		return fmt.Sprintf("export %s=%s", key, shellQuote(value)), nil
	case FormatDotenv:
		return fmt.Sprintf("%s=%s", key, value), nil
	default:
		return "", fmt.Errorf("exporter: unsupported format %q", e.format)
	}
}

// shellQuote wraps a value in single quotes, escaping any existing single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
