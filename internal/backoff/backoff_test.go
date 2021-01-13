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
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		want      *backoff
		checkFunc func(got, want *backoff) error
	}

	tests := []test{
		{
			name: "returns backoff instance",
			opts: []Option{
				WithBackOffFactor(0.5),
			},
			want: &backoff{
				initialDuration:  float64(10 * time.Millisecond),
				backoffTimeLimit: 5 * time.Minute,
				maxDuration:      float64(time.Hour),
				jitterLimit:      float64(time.Minute),
				backoffFactor:    1.1,
				maxRetryCount:    50,
				errLog:           true,
				durationLimit:    float64(time.Hour) / 1.1,
			},
			checkFunc: func(got *backoff, want *backoff) error {
				got.jittedInitialDuration, want.jittedInitialDuration = 1, 1
				if !reflect.DeepEqual(got, want) {
					return errors.Errorf("not equals. want: %v, got: %v", got, want)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			if err := tt.checkFunc(got.(*backoff), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_backoff_Do(t *testing.T) {
	t.Parallel()
	type args struct {
		fn   func(context.Context) (interface{}, bool, error)
		opts []Option
	}

	type test struct {
		name      string
		args      args
		ctxFn     func() (context.Context, context.CancelFunc)
		checkFunc func(got, want error) error
		want      error
	}

	tests := []test{
		func() test {
			cnt := 0
			fn := func(context.Context) (interface{}, bool, error) {
				cnt++
				return nil, false, nil
			}

			return test{
				name: "returns response and nil when function return not nil",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 1 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if !errors.Is(want, got) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: nil,
			}
		}(),

		func() test {
			cnt := 0
			fn := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt == 2 {
					return nil, false, nil
				}
				return nil, true, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "returns response and nil when retried twice and did not return an error",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 2 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if !errors.Is(want, got) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: nil,
			}
		}(),

		func() test {
			cnt := 0
			err := errors.New("not retryable error")
			fn := func(context.Context) (interface{}, bool, error) {
				cnt++
				return nil, false, err
			}

			return test{
				name: "returns error when retryable is false",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 1 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 1, cnt)
					}

					if !errors.Is(want, got) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: err,
			}
		}(),

		func() test {
			cnt := 0
			fn := func(context.Context) (interface{}, bool, error) {
				cnt++
				return nil, true, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "returns error when retrying the maximum number of times",
				args: args{
					fn: fn,
					opts: []Option{
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 7 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 7, cnt)
					}

					if want.Error() != got.Error() {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: errors.New("error (7)"),
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			cnt := 0
			fn := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt == 2 {
					cancel()
				}
				return nil, true, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "return response and error when context canceled",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return ctx, cancel
				},
				checkFunc: func(got, want error) error {
					if cnt != 2 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if got.Error() != want.Error() {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: errors.Wrap(context.Canceled, errors.New("error (2)").Error()),
			}
		}(),

		func() test {
			err := errors.New("error")
			fn := func(context.Context) (interface{}, bool, error) {
				return nil, true, err
			}

			return test{
				name: "return response and error when backoff timeout",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
						WithBackOffTimeLimit("0s"),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if !errors.Is(got, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
				want: err,
			}
		}(),
	}

	log.Init()
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := test.ctxFn()
			defer cancel()
			_, err := New(test.args.opts...).Do(ctx, test.args.fn)
			if test.want == nil && err != nil {
				t.Errorf("Do return err: %v", err)
			}

			if err := test.checkFunc(err, test.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	type test struct {
		name string
		bo   *backoff
	}

	tests := []test{
		{
			name: "processing is successes",
			bo: &backoff{
				wg: sync.WaitGroup{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bo.Close()
		})
	}
}

func Test_backoff_addJitter(t *testing.T) {
	t.Parallel()
	type args struct {
		dur float64
	}
	type fields struct {
		wg                    sync.WaitGroup
		backoffFactor         float64
		initialDuration       float64
		jittedInitialDuration float64
		jitterLimit           float64
		durationLimit         float64
		maxDuration           float64
		maxRetryCount         int
		backoffTimeLimit      time.Duration
		errLog                bool
	}
	type want struct {
		want float64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, float64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got float64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dur: 0,
		       },
		       fields: fields {
		           wg: sync.WaitGroup{},
		           backoffFactor: 0,
		           initialDuration: 0,
		           jittedInitialDuration: 0,
		           jitterLimit: 0,
		           durationLimit: 0,
		           maxDuration: 0,
		           maxRetryCount: 0,
		           backoffTimeLimit: nil,
		           errLog: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           dur: 0,
		           },
		           fields: fields {
		           wg: sync.WaitGroup{},
		           backoffFactor: 0,
		           initialDuration: 0,
		           jittedInitialDuration: 0,
		           jitterLimit: 0,
		           durationLimit: 0,
		           maxDuration: 0,
		           maxRetryCount: 0,
		           backoffTimeLimit: nil,
		           errLog: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &backoff{
				wg:                    test.fields.wg,
				backoffFactor:         test.fields.backoffFactor,
				initialDuration:       test.fields.initialDuration,
				jittedInitialDuration: test.fields.jittedInitialDuration,
				jitterLimit:           test.fields.jitterLimit,
				durationLimit:         test.fields.durationLimit,
				maxDuration:           test.fields.maxDuration,
				maxRetryCount:         test.fields.maxRetryCount,
				backoffTimeLimit:      test.fields.backoffTimeLimit,
				errLog:                test.fields.errLog,
			}

			got := b.addJitter(test.args.dur)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_backoff_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		wg                    sync.WaitGroup
		backoffFactor         float64
		initialDuration       float64
		jittedInitialDuration float64
		jitterLimit           float64
		durationLimit         float64
		maxDuration           float64
		maxRetryCount         int
		backoffTimeLimit      time.Duration
		errLog                bool
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           wg: sync.WaitGroup{},
		           backoffFactor: 0,
		           initialDuration: 0,
		           jittedInitialDuration: 0,
		           jitterLimit: 0,
		           durationLimit: 0,
		           maxDuration: 0,
		           maxRetryCount: 0,
		           backoffTimeLimit: nil,
		           errLog: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           wg: sync.WaitGroup{},
		           backoffFactor: 0,
		           initialDuration: 0,
		           jittedInitialDuration: 0,
		           jitterLimit: 0,
		           durationLimit: 0,
		           maxDuration: 0,
		           maxRetryCount: 0,
		           backoffTimeLimit: nil,
		           errLog: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &backoff{
				wg:                    test.fields.wg,
				backoffFactor:         test.fields.backoffFactor,
				initialDuration:       test.fields.initialDuration,
				jittedInitialDuration: test.fields.jittedInitialDuration,
				jitterLimit:           test.fields.jitterLimit,
				durationLimit:         test.fields.durationLimit,
				maxDuration:           test.fields.maxDuration,
				maxRetryCount:         test.fields.maxRetryCount,
				backoffTimeLimit:      test.fields.backoffTimeLimit,
				errLog:                test.fields.errLog,
			}

			b.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
