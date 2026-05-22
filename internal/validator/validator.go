package validator

import (
	"fmt"

	"github.com/envchain/envchain/internal/chain"
)

// Result holds the outcome of a validation run.
type Result struct {
	Passed  []string
	Failed  []string
	Warnings []string
}

// IsValid returns true when no variables failed validation.
func (r *Result) IsValid() bool {
	return len(r.Failed) == 0
}

// Summary returns a human-readable one-line summary.
func (r *Result) Summary() string {
	return fmt.Sprintf("%d passed, %d failed, %d warnings",
		len(r.Passed), len(r.Failed), len(r.Warnings))
}

// Validator validates a chain of environment variables.
type Validator struct {
	chain *chain.Chain
}

// New creates a new Validator for the given Chain.
func New(c *chain.Chain) *Validator {
	return &Validator{chain: c}
}

// Validate iterates over all variables in the chain and checks their status.
// Variables that are resolved are marked as passed; missing or empty ones are
// recorded as failures; unresolved dependencies produce warnings.
func (v *Validator) Validate() *Result {
	result := &Result{}

	for _, variable := range v.chain.Variables() {
		switch variable.Status() {
		case chain.StatusResolved:
			result.Passed = append(result.Passed, variable.Name())
		case chain.StatusMissing, chain.StatusEmpty:
			result.Failed = append(result.Failed, variable.Name())
		case chain.StatusUnresolved:
			result.Warnings = append(result.Warnings, variable.Name())
		}
	}

	return result
}
