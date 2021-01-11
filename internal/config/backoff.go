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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/backoff"

type Backoff struct {
	InitialDuration  string  `json:"initial_duration" yaml:"initial_duration"`
	BackoffTimeLimit string  `json:"backoff_time_limit" yaml:"backoff_time_limit"`
	MaximumDuration  string  `json:"maximum_duration" yaml:"maximum_duration"`
	JitterLimit      string  `json:"jitter_limit" yaml:"jitter_limit"`
	BackoffFactor    float64 `json:"backoff_factor" yaml:"backoff_factor"`
	RetryCount       int     `json:"retry_count" yaml:"retry_count"`
	EnableErrorLog   bool    `json:"enable_error_log" yaml:"enable_error_log"`
}

func (b *Backoff) Bind() *Backoff {
	b.InitialDuration = GetActualValue(b.InitialDuration)
	b.BackoffTimeLimit = GetActualValue(b.BackoffTimeLimit)
	b.MaximumDuration = GetActualValue(b.MaximumDuration)
	b.JitterLimit = GetActualValue(b.JitterLimit)
	return b
}

func (b *Backoff) Opts() []backoff.Option {
	opts := make([]backoff.Option, 0, 7)
	opts = append(opts,
		backoff.WithBackOffTimeLimit(b.BackoffTimeLimit),
		backoff.WithInitialDuration(b.InitialDuration),
		backoff.WithJitterLimit(b.JitterLimit),
		backoff.WithMaximumDuration(b.MaximumDuration),
		backoff.WithRetryCount(b.RetryCount),
		backoff.WithBackOffFactor(b.BackoffFactor),
	)

	if b.EnableErrorLog {
		opts = append(opts,
			backoff.WithEnableErrorLog(),
		)
	}

	return opts
}
