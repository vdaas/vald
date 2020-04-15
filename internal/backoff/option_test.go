//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package backoff

import (
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

func TestWithInitialDuration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.initialDuration != float64(10*time.Second) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "dur",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.initialDuration != float64(500*time.Millisecond) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithInitialDuration(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithMaximumDuration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxDuration != float64(10*time.Second) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "dur",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxDuration != float64(5*time.Hour) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMaximumDuration(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithJitterLimit(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.jitterLimit != float64(10*time.Second) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "dur",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.jitterLimit != float64(time.Minute) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithJitterLimit(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackOffFactor(t *testing.T) {
	type test struct {
		name      string
		f         float64
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			f:    10.0,
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffFactor != 10.0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			f:    -10.0,
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffFactor != 0.0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackOffFactor(tt.f)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRetryCount(t *testing.T) {
	type test struct {
		name      string
		c         int
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			c:    10,
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxRetryCount != 10 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			c:    -10,
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxRetryCount != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetryCount(tt.c)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackOffTimeLimit(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffTimeLimit != 10*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "dur",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffTimeLimit != 20*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackOffTimeLimit(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithEWithEnableErrorLog(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.errLog != true {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEnableErrorLog()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDisableErrorLog(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.errLog != false {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDisableErrorLog()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	type test struct {
		name      string
		checkFunc func([]Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opts []Option) error {
				got := new(backoff)

				for _, opt := range opts {
					opt(got)
				}

				if got.initialDuration != float64(10*time.Millisecond) {
					return errors.New("invalid param (initialDuration) was set")
				}

				if got.backoffTimeLimit != 5*time.Minute {
					return errors.New("invalid param (backoffTimeLimit) was set")
				}

				if got.maxDuration != float64(time.Hour) {
					return errors.New("invalid param (maxDuration) was set")
				}

				if got.jitterLimit != float64(time.Minute) {
					return errors.New("invalid param (jittedInitialDuration) was set")
				}

				if got.backoffFactor != 1.5 {
					return errors.New("invalid param (backoffFactor) was set")
				}

				if got.maxRetryCount != 50 {
					return errors.New("invalid param (maxRetryCount) was set")
				}

				if got.errLog != true {
					return errors.New("invalid param (errLog) was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.checkFunc(defaultOpts); err != nil {
				t.Error(err)
			}
		})
	}
}
