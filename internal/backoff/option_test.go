package backoff

import (
	"fmt"
	"testing"
	"time"
)

func TestWithInitialDuration(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.initialDuration != float64(10*time.Second) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				dur: "dur",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.initialDuration != float64(500*time.Millisecond) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithInitialDuration(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithMaximumDuration(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxDuration != float64(10*time.Second) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				dur: "dur",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxDuration != float64(5*time.Hour) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMaximumDuration(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithJitterLimit(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.jitterLimit != float64(10*time.Second) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				dur: "dur",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.jitterLimit != float64(time.Minute) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithJitterLimit(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackOffFactor(t *testing.T) {
	type args struct {
		f float64
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				f: 10.0,
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffFactor != 10.0 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				f: -10.0,
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffFactor != 0.0 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackOffFactor(tt.args.f)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRetryCount(t *testing.T) {
	type args struct {
		c int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				c: 10,
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxRetryCount != 10 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				c: -10,
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.maxRetryCount != 0 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetryCount(tt.args.c)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackOffTimeLimit(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffTimeLimit != 10*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				dur: "dur",
			},
			checkFunc: func(opt Option) error {
				got := new(backoff)
				opt(got)

				if got.backoffTimeLimit != 20*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackOffTimeLimit(tt.args.dur)
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
					return fmt.Errorf("invalid param was set")
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
					return fmt.Errorf("invalid param was set")
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
					return fmt.Errorf("invalid param (initialDuration) was set")
				}

				if got.backoffTimeLimit != 5*time.Minute {
					return fmt.Errorf("invalid param (backoffTimeLimit) was set")
				}

				if got.maxDuration != float64(time.Hour) {
					return fmt.Errorf("invalid param (maxDuration) was set")
				}

				if got.jitterLimit != float64(time.Minute) {
					return fmt.Errorf("invalid param (jittedInitialDuration) was set")
				}

				if got.backoffFactor != 1.5 {
					return fmt.Errorf("invalid param (backoffFactor) was set")
				}

				if got.maxRetryCount != 50 {
					return fmt.Errorf("invalid param (maxRetryCount) was set")
				}

				if got.errLog != true {
					return fmt.Errorf("invalid param (errLog) was set")
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
