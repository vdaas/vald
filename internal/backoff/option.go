// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package backoff provides backoff function controller
package backoff

import (
	"time"

	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*backoff)

var (
	defaultOpts = []Option{
		WithInitialDuration("10ms"),
		WithBackOffTimeLimit("5m"),
		WithMaximumDuration("1h"),
		WithJitterWidthLimit("1m"),
		WithBackOffFactor(1.5),
		WithRetryCount(50),
		WithEnableErrorLog(),
	}
)

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

func WithJitterWidthLimit(dur string) Option {
	return func(b *backoff) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		b.jitterWidthLimit = float64(d)
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
		b.errLog = true
	}
}
