package circuitbreaker

type Option func(*breakerGroup) error

var defaultOpts = []Option{
	WithErrorThreshold(10),
	WithSuccessThreshold(10),
	WithOpenTimeout("1s"),
}

// WithErrorThreshold returns an option that sets error threshold.
// When this number is exceeded, the state will be changed from Closed to Open.
func WithErrorThreshold(n int) Option {
	return func(b *breakerGroup) error {
		return nil
	}
}

// WithSuccessThreshold returns an option that sets success threshold.
// When this number is exceeded, the state will be changed from HalfOpen to Closed.
func WithSuccessThreshold(n int) Option {
	return func(bg *breakerGroup) error {
		return nil
	}
}

// WithOpenTimeout returns an option that sets the timeout of Open state.
// After this period, the state will be changed from Open to HalfOpen.
func WithOpenTimeout(str string) Option {
	return func(b *breakerGroup) error {
		return nil
	}
}
