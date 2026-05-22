package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/chain"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes chain validation results to an output.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a new Reporter writing to the given writer.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes the validation result of the given chain.
func (r *Reporter) Report(c *chain.Chain) error {
	switch r.format {
	case FormatJSON:
		return r.reportJSON(c)
	default:
		return r.reportText(c)
	}
}

func (r *Reporter) reportText(c *chain.Chain) error {
	vars := c.Variables()
	resolved := 0
	for _, v := range vars {
		if v.Status == chain.StatusResolved {
			resolved++
		}
	}

	fmt.Fprintf(r.out, "Chain: %s\n", c.Name())
	fmt.Fprintf(r.out, "Variables: %d resolved / %d total\n", resolved, len(vars))
	fmt.Fprintln(r.out, strings.Repeat("-", 40))

	for _, v := range vars {
		marker := "✓"
		if v.Status != chain.StatusResolved {
			marker = "✗"
		}
		fmt.Fprintf(r.out, "  [%s] %s (%s)\n", marker, v.Name, v.Status)
	}
	return nil
}

func (r *Reporter) reportJSON(c *chain.Chain) error {
	vars := c.Variables()
	fmt.Fprintf(r.out, "{\"chain\":%q,\"variables\":[\n", c.Name())
	for i, v := range vars {
		comma := ","
		if i == len(vars)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, "  {\"name\":%q,\"status\":%q}%s\n", v.Name, v.Status, comma)
	}
	fmt.Fprintln(r.out, "]}")
	return nil
}
