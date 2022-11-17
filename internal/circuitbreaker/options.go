//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
	WithHalfOpenErrorRate(0.5),
	WithMinSamples(1000),
	WithOpenTimeout("1s"),
	WithClosedRefreshTimeout("10s"),
}

// WithClosedErrorRate returns an option that sets error rate when breaker state is "Closed".
// The rate is expected to be between 0 and 1.0.
// When the rate is exceeded, the breaker state will be changed from "Closed" to "Open".
func WithClosedErrorRate(f float32) BreakerOption {
	return func(b *breaker) error {
		if f < 0 || f > 1.0 {
			return errors.NewErrInvalidOption("closedErrorRate", f)
		}
		b.closedErrRate = f
		return nil
	}
}

// WithClosedErrorTripper returns an option that sets whether it should trip when in "Closed" state.
func WithClosedErrorTripper(tp Tripper) BreakerOption {
	return func(b *breaker) error {
		if tp == nil {
			return errors.NewErrInvalidOption("closedErrTripper", tp)
		}
		b.closedErrShouldTrip = tp
		return nil
	}
}

// WithHalfOpenErrorRate returns an option that sets error rate when breaker state is "HalfOpen".
// The rate is expected to be between 0 and 1.0.
// When the rate is exceeded, the breaker state will be changed from "HalfOpen" to "Open".
func WithHalfOpenErrorRate(f float32) BreakerOption {
	return func(b *breaker) error {
		if f < 0 || f > 1.0 {
			return errors.NewErrInvalidOption("halfOpenErrorRate", f)
		}
		b.halfOpenErrRate = f
		return nil
	}
}

// WithHalfOpenErrorTripper returns an option that sets whether it should trip when in "Half-Open" state.
func WithHalfOpenErrorTripper(tp Tripper) BreakerOption {
	return func(b *breaker) error {
		if tp == nil {
			return errors.NewErrInvalidOption("halfOpenErrTripper", tp)
		}
		b.halfOpenErrShouldTrip = tp
		return nil
	}
}

// WithMinSamples returns an option that sets minimum sample count.
func WithMinSamples(min int64) BreakerOption {
	return func(b *breaker) error {
		if min < 1 {
			return errors.NewErrInvalidOption("minSamples", min)
		}
		b.minSamples = min
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

// WithClosedRefreshTimeout returns an option that sets the timeout of "Closed" state.
// After this period, the counter will be refreshed.
func WithClosedRefreshTimeout(timeout string) BreakerOption {
	return func(b *breaker) error {
		if len(timeout) == 0 {
			return errors.NewErrInvalidOption("closedRefreshTimeout", timeout)
		}

		d, err := timeutil.Parse(timeout)
		if err != nil {
			return errors.NewErrInvalidOption("closedRefreshTimeout", timeout, err)
		}
		b.cloedRefreshTimeout = d
		return nil
	}
}
