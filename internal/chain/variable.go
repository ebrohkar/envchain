// Package chain provides types and logic for representing environment variable
// dependency chains used in service configurations.
package chain

import "fmt"

// VarStatus represents the resolution status of an environment variable.
type VarStatus int

const (
	StatusUnknown  VarStatus = iota
	StatusResolved           // variable is present and non-empty
	StatusMissing            // variable is absent from the environment
	StatusEmpty              // variable exists but is empty
)

func (s VarStatus) String() string {
	switch s {
	case StatusResolved:
		return "resolved"
	case StatusMissing:
		return "missing"
	case StatusEmpty:
		return "empty"
	default:
		return "unknown"
	}
}

// Variable represents a single environment variable node in a dependency chain.
type Variable struct {
	Name     string
	Value    string
	Status   VarStatus
	DependsOn []string // names of variables this one depends on
}

// Validate checks whether the variable's status is acceptable.
func (v *Variable) Validate() error {
	switch v.Status {
	case StatusMissing:
		return fmt.Errorf("variable %q is missing from the environment", v.Name)
	case StatusEmpty:
		return fmt.Errorf("variable %q is set but empty", v.Name)
	case StatusResolved:
		return nil
	default:
		return fmt.Errorf("variable %q has unknown status", v.Name)
	}
}
