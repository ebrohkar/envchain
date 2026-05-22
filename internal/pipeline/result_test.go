package pipeline

import (
	"strings"
	"testing"
)

func TestResult_AllOK_BothTrue(t *testing.T) {
	r := &Result{ChainOK: true, ValidatorOK: true}
	if !r.AllOK() {
		t.Error("expected AllOK to be true")
	}
}

func TestResult_AllOK_ChainFailed(t *testing.T) {
	r := &Result{ChainOK: false, ValidatorOK: true}
	if r.AllOK() {
		t.Error("expected AllOK to be false when chain failed")
	}
}

func TestResult_AllOK_ValidatorFailed(t *testing.T) {
	r := &Result{ChainOK: true, ValidatorOK: false}
	if r.AllOK() {
		t.Error("expected AllOK to be false when validator failed")
	}
}

func TestResult_Summary_OK(t *testing.T) {
	r := &Result{ChainOK: true, ValidatorOK: true}
	if !strings.HasPrefix(r.Summary(), "OK:") {
		t.Errorf("expected summary to start with OK:, got: %s", r.Summary())
	}
}

func TestResult_Summary_Fail(t *testing.T) {
	r := &Result{
		ChainOK:     false,
		ValidatorOK: false,
		Errors:      []string{"unresolved reference: DB_PORT"},
	}
	s := r.Summary()
	if !strings.HasPrefix(s, "FAIL:") {
		t.Errorf("expected summary to start with FAIL:, got: %s", s)
	}
	if !strings.Contains(s, "DB_PORT") {
		t.Errorf("expected summary to contain error detail, got: %s", s)
	}
}

func TestResult_ExitCode(t *testing.T) {
	ok := &Result{ChainOK: true, ValidatorOK: true}
	if ok.ExitCode() != 0 {
		t.Errorf("expected exit code 0, got %d", ok.ExitCode())
	}

	fail := &Result{ChainOK: false, ValidatorOK: true}
	if fail.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", fail.ExitCode())
	}
}
