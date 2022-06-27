package circuitbreaker

type Option func(*breakerManager) error

var defaultOpts = []Option{}

func WithBreakerOpts(opts ...BreakerOption) Option {
	return func(bm *breakerManager) error {
		if bm.opts != nil && len(bm.opts) > 0 {
			bm.opts = append(bm.opts, opts...)
		} else {
			bm.opts = opts
		}
		return nil
	}
}

type BreakerOption func(*breaker) error

var defaultBreakerOpts = []BreakerOption{
	WithClosedErrorThreshold(10),
	WithHarfOpenSuccessThreshold(10),
	WithOpenTimeout("1s"),
}

// WithErrorThreshold returns an option that sets error threshold.
// When this number is exceeded, the state will be changed from Closed to Open.
func WithClosedErrorThreshold(n int) BreakerOption {
	return func(b *breaker) error {
		return nil
	}
}

// WithHarfOpenSuccessThreshold returns an option that sets success threshold.
// When this number is exceeded, the state will be changed from HalfOpen to Closed.
func WithHarfOpenSuccessThreshold(n int) BreakerOption {
	return func(bg *breaker) error {
		return nil
	}
}

// WithOpenTimeout returns an option that sets the timeout of Open state.
// After this period, the state will be changed from Open to HalfOpen.
func WithOpenTimeout(str string) BreakerOption {
	return func(b *breaker) error {
		return nil
	}
}
