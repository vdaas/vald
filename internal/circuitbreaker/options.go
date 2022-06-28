package circuitbreaker

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/timeutil"
)

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
	WithClosedErrorRate(0.7),
	WithOpenTimeout("1s"),
}

// WithClosedErrorRate returns an option that sets error rate.
// The rate is expected to be between 0 and 1.0.
// When the rate is exceeded, the breaker state will be changed from "Closed" to "Open".
func WithClosedErrorRate(f float64) BreakerOption {
	return func(b *breaker) error {
		if f < 0 || f > 1.0 {
			return errors.NewErrInvalidOption("closedErrorRate", f)
		}
		b.closedErrRate = f
		return nil
	}
}

// WithOpenTimeout returns an option that sets the timeout of "Open" state.
// After this period, the state will be changed from "Open" to "HalfOpen".
func WithOpenTimeout(timeout string) BreakerOption {
	return func(b *breaker) error {
		if len(timeout) == 0 {
			return errors.NewErrInvalidOption("openTimeout", timeout)
		}

		d, err := timeutil.Parse(timeout)
		if err != nil {
			return errors.NewErrInvalidOption("openTimeout", timeout, err)
		}
		b.openTimeout = d
		return nil
	}
}
