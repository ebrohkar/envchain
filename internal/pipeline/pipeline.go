// Package pipeline orchestrates the full envchain validation pipeline:
// loading variables, resolving references, and validating rules.
package pipeline

import (
	"fmt"
	"io"

	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/reporter"
	"github.com/user/envchain/internal/resolver"
	"github.com/user/envchain/internal/validator"
)

// Options configures a Pipeline run.
type Options struct {
	// EnvSource provides environment variable lookups.
	EnvSource func(string) (string, bool)
	// Format is the output format (text or json).
	Format reporter.Format
	// Writer is the output destination; defaults to stdout if nil.
	Writer io.Writer
}

// Result holds the aggregated outcome of a pipeline run.
type Result struct {
	ChainOK     bool
	ValidatorOK bool
	Errors      []string
}

// OK returns true only when both the chain resolution and validation succeeded.
func (r *Result) OK() bool {
	return r.ChainOK && r.ValidatorOK
}

// Run executes the full pipeline against the provided chain.
func Run(c *chain.Chain, rules map[string][]validator.Rule, opts Options) (*Result, error) {
	if c == nil {
		return nil, fmt.Errorf("pipeline: chain must not be nil")
	}

	res := &Result{}

	// Step 1: resolve variable references within the chain.
	r := resolver.New(opts.EnvSource)
	resolved, missing := r.ResolveChain(c)
	if len(missing) > 0 {
		for _, m := range missing {
			res.Errors = append(res.Errors, fmt.Sprintf("unresolved reference: %s", m))
		}
	}
	res.ChainOK = len(missing) == 0

	// Step 2: validate resolved variables against supplied rules.
	v := validator.New(opts.EnvSource)
	for name, varRules := range rules {
		v.AddRules(name, varRules...)
	}
	vResult := v.Validate(resolved)
	res.ValidatorOK = vResult.AllPassed()
	for _, f := range vResult.Failures() {
		res.Errors = append(res.Errors, f.String())
	}

	// Step 3: report results.
	rep := reporter.New(opts.Writer, opts.Format)
	rep.Report(c, vResult)

	return res, nil
}
