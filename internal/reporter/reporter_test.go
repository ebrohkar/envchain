package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/reporter"
)

func mockEnv(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestReporter_Text_AllResolved(t *testing.T) {
	env := mockEnv(map[string]string{"HOST": "localhost", "PORT": "8080"})
	c := chain.New("svc", env)
	c.Add("HOST", "")
	c.Add("PORT", "")
	c.Resolve()

	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Chain: svc") {
		t.Errorf("expected chain name in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2 resolved / 2 total") {
		t.Errorf("expected resolved count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[✓] HOST") {
		t.Errorf("expected resolved marker for HOST, got:\n%s", out)
	}
}

func TestReporter_Text_MissingVariable(t *testing.T) {
	env := mockEnv(map[string]string{"HOST": "localhost"})
	c := chain.New("svc", env)
	c.Add("HOST", "")
	c.Add("PORT", "")
	c.Resolve()

	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	r.Report(c)

	out := buf.String()
	if !strings.Contains(out, "[✗] PORT") {
		t.Errorf("expected unresolved marker for PORT, got:\n%s", out)
	}
}

func TestReporter_JSON_Output(t *testing.T) {
	env := mockEnv(map[string]string{"DB_URL": "postgres://localhost/db"})
	c := chain.New("db-chain", env)
	c.Add("DB_URL", "")
	c.Resolve()

	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Report(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"chain":"db-chain"`) {
		t.Errorf("expected chain name in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"name":"DB_URL"`) {
		t.Errorf("expected variable name in JSON output, got:\n%s", out)
	}
}

func TestReporter_NilWriter_UsesStdout(t *testing.T) {
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
