package chain

import (
	"fmt"
	"os"
)

// Chain holds a collection of Variables and their declared dependencies.
type Chain struct {
	vars map[string]*Variable
}

// New creates an empty Chain.
func New() *Chain {
	return &Chain{vars: make(map[string]*Variable)}
}

// Add registers a variable name with optional dependencies into the chain.
func (c *Chain) Add(name string, dependsOn ...string) {
	c.vars[name] = &Variable{
		Name:      name,
		DependsOn: dependsOn,
	}
}

// Resolve reads actual values from the environment and sets statuses.
func (c *Chain) Resolve(environ func(string) (string, bool)) {
	if environ == nil {
		enviran := func(k string) (string, bool) { return os.LookupEnv(k) }
		enviran = enviran // satisfy compiler; reassign
		enviran = func(k string) (string, bool) { return os.LookupEnv(k) }
		for _, v := range c.vars {
			c.resolveVar(v, enviran)
		}
		return
	}
	for _, v := range c.vars {
		c.resolveVar(v, environ)
	}
}

func (c *Chain) resolveVar(v *Variable, lookup func(string) (string, bool)) {
	val, ok := lookup(v.Name)
	if !ok {
		v.Status = StatusMissing
		return
	}
	if val == "" {
		v.Status = StatusEmpty
		return
	}
	v.Value = val
	v.Status = StatusResolved
}

// Validate checks all variables and their dependency ordering.
// Returns a slice of errors — one per problematic variable.
func (c *Chain) Validate() []error {
	var errs []error
	for _, v := range c.vars {
		if err := v.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}
		for _, dep := range v.DependsOn {
			depVar, exists := c.vars[dep]
			if !exists {
				errs = append(errs, fmt.Errorf("variable %q depends on undeclared variable %q", v.Name, dep))
				continue
			}
			if depVar.Status != StatusResolved {
				errs = append(errs, fmt.Errorf("variable %q depends on unresolved variable %q (%s)", v.Name, dep, depVar.Status))
			}
		}
	}
	return errs
}

// Variables returns a copy of the registered variable map.
func (c *Chain) Variables() map[string]*Variable {
	copy := make(map[string]*Variable, len(c.vars))
	for k, v := range c.vars {
		copy[k] = v
	}
	return copy
}
