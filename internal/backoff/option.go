//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package backoff provides backoff function controller
package backoff

import (
	"time"

	"github.com/vdaas/vald/internal/timeutil"
)

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

func WithInitialDuration(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = 500 * time.Millisecond
		}
		b.initialDuration = float64(d)
	}
}

func WithMaximumDuration(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = 5 * time.Hour
		}
		b.maxDuration = float64(d)
	}
}

func WithJitterLimit(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		b.jitterLimit = float64(d)
	}
}

func WithBackOffFactor(f float64) Option {
	return func(b *backoff) {
		if f <= 0.0 {
			return
		}
		b.backoffFactor = f
	}
}

func WithRetryCount(c int) Option {
	return func(b *backoff) {
		if c <= 0 {
			return
		}
		b.maxRetryCount = c
	}
}

func WithBackOffTimeLimit(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 20
		}
		b.backoffTimeLimit = d
	}
}

func WithEnableErrorLog() Option {
	return func(b *backoff) {
		b.errLog = true
	}
}

func WithDisableErrorLog() Option {
	return func(b *backoff) {
		b.errLog = false
	}
}
