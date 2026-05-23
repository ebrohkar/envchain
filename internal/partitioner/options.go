package partitioner

// Option is a functional option for configuring a Partitioner.
type Option func(*Partitioner)

// WithFn sets the partition function on the Partitioner.
func WithFn(fn PartitionFn) Option {
	return func(p *Partitioner) {
		p.fn = fn
	}
}

// NewWithOptions creates a Partitioner applying the given options.
// If no partition function is provided via options, a default "all" partition is used.
func NewWithOptions(opts ...Option) *Partitioner {
	p := &Partitioner{
		fn: func(_, _ string) string { return "all" },
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
