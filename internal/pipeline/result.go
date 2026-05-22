package pipeline

import "strings"

// AllOK returns true when both the chain and validator checks passed.
func (r *Result) AllOK() bool {
	return r.ChainOK && r.ValidatorOK
}

// Summary returns a human-readable single-line summary of the result.
func (r *Result) Summary() string {
	if r.AllOK() {
		return "OK: all environment variables resolved and validated"
	}

	var parts []string
	if !r.ChainOK {
		parts = append(parts, "chain resolution failed")
	}
	if !r.ValidatorOK {
		parts = append(parts, "validation failed")
	}

	detail := ""
	if len(r.Errors) > 0 {
		detail = ": " + strings.Join(r.Errors, "; ")
	}

	return "FAIL: " + strings.Join(parts, ", ") + detail
}

// ExitCode returns 0 when all checks passed, 1 otherwise.
// Intended for use as a process exit code.
func (r *Result) ExitCode() int {
	if r.AllOK() {
		return 0
	}
	return 1
}
