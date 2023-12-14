//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package backoff provides backoff function controller
package backoff

import (
	"time"

	"github.com/vdaas/vald/internal/timeutil"
)

// Option represents the functional option for backoff.
type Option func(*backoff)

var defaultOptions = []Option{
	WithInitialDuration("10ms"),
	WithBackOffTimeLimit("5m"),
	WithMaximumDuration("1h"),
	WithJitterLimit("1m"),
	WithBackOffFactor(1.5),
	WithRetryCount(50),
	WithEnableErrorLog(),
}

// WithInitialDuration returns the option to set the initial duration of backoff.
func WithInitialDuration(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = 500 * time.Millisecond
		}
		b.initialDuration = float64(d)
	}
}

// WithMaximumDuration returns the option to set the maximum duration of backoff.
func WithMaximumDuration(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = 5 * time.Hour
		}
		b.maxDuration = float64(d)
	}
}

// WithJitterLimit returns the option to set the jitter limit duration of backoff.
func WithJitterLimit(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		b.jitterLimit = float64(d)
	}
}

// WithBackOffFactor returns the option to set the factor of backoff.
func WithBackOffFactor(f float64) Option {
	return func(b *backoff) {
		if f <= 0.0 {
			return
		}
		b.backoffFactor = f
	}
}

// WithRetryCount returns the option to set the retry count of backoff.
func WithRetryCount(c int) Option {
	return func(b *backoff) {
		if c <= 0 {
			return
		}
		b.maxRetryCount = c
	}
}

// WithBackOffTimeLimit returns the option to set the limit of backoff.
func WithBackOffTimeLimit(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 20
		}
		b.backoffTimeLimit = d
	}
}

// WithEnableErrorLog returns the option to set the enable for error log.
func WithEnableErrorLog() Option {
	return func(b *backoff) {
		b.errLog = true
	}
}

// WithDisableErrorLog returns the option to set the disable for error log.
func WithDisableErrorLog() Option {
	return func(b *backoff) {
		b.errLog = false
	}
}
