package patcher

// Option is a functional option for configuring a Patcher.
type Option func(*Patcher)

// WithStrategy sets the conflict resolution strategy.
func WithStrategy(s Strategy) Option {
	return func(p *Patcher) {
		p.strategy = s
	}
}

// NewWithOptions constructs a Patcher using the provided options.
// Defaults to StrategyOverwrite if no strategy option is given.
func NewWithOptions(opts ...Option) *Patcher {
	p := &Patcher{strategy: StrategyOverwrite}
	for _, o := range opts {
		o(p)
	}
	return p
}
