//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

const str = "success"

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	info.Init("backoff test")
	os.Exit(m.Run())
}

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

func Test_backoff_addJitter(t *testing.T) {
	t.Parallel()
	type args struct {
		dur float64
	}
	type fields struct {
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
		func() test {
			return test{
				name: "success when dur is 0",
				args: args{
					dur: 0,
				},
				fields: fields{
					jitterLimit: 100,
				},
				want: want{},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &backoff{
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
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_backoff_Close(t *testing.T) {
	t.Parallel()
	type fields struct{}
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
		{
			name:   "success backoff Close",
			fields: fields{},
			want:   want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &backoff{}

			b.Close()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_backoff_Do(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(ctx context.Context) (val interface{}, retryable bool, err error)
	}
	type fields struct {
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
		wantRes interface{}
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			err := errors.New("error is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				return nil, false, err
			}
			return test{
				name: "return nil response and error when function returns (nil, false, error) and not retriable",
				args: args{
					ctx: ctx,
					f:   f,
				},
				want: want{
					err: err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			f := func(context.Context) (interface{}, bool, error) {
				return nil, true, nil
			}
			return test{
				name: "return nil response and nil error when function returns (nil, true, nil) and not retriable",
				args: args{
					ctx: ctx,
					f:   f,
				},
				want: want{},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				return nil, false, err
			}
			return test{
				name: "return nil response and error when function return (nil, false, error) and maxRetryCount = 0",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         0,
					maxDuration:           0,
					maxRetryCount:         0,
					backoffTimeLimit:      0,
					errLog:                false,
				},
				want: want{
					wantRes: nil,
					err:     err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				return str, true, err
			}
			return test{
				name: "return response and nil error when function return (string, true, error) and maxRetryCount = 1",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         0,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      10 * time.Minute,
					errLog:                false,
				},
				want: want{
					wantRes: str,
					err:     err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			cnt := 0
			f := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt == 2 {
					return str, false, err
				}
				return str, true, err
			}
			return test{
				name: "return response and error when function return (string, false, error) at 2nd times and maxRetryCount = 1",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         0,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      10 * time.Minute,
					errLog:                false,
				},
				want: want{
					wantRes: str,
					err:     err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			cnt := 0
			f := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt == 2 {
					return str, true, nil
				}
				return str, true, err
			}
			return test{
				name: "return response and nil error when function return (string, true, nil) at 2nd times and maxRetryCount = 1",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         0,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      10 * time.Minute,
					errLog:                false,
				},
				want: want{
					wantRes: str,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				return str, true, err
			}
			return test{
				name: "return response and error when function return (string, true, error) and maxRetryCount = 1, errLog is true",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         10,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      10 * time.Minute,
					errLog:                true,
				},
				want: want{
					wantRes: str,
					err:     err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				return str, true, err
			}
			return test{
				name: "return nil response and error when function returns (string, true, error) and context will be closed due to timelimit",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         0,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      0,
					errLog:                true,
				},
				want: want{
					err: errors.ErrBackoffTimeout(err),
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("erros is occurred")
			f := func(context.Context) (interface{}, bool, error) {
				cancel()
				return str, true, err
			}
			return test{
				name: "return nil response and error when function returns (string, true, error) and calls cancel()",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         10,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      100 * time.Microsecond,
					errLog:                true,
				},
				want: want{
					err: err,
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("erros is occurred")
			cnt := 0
			f := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt > 1 {
					cancel()
				}
				return str, true, err
			}
			return test{
				name: "return nil response and error when function returns (string, true, error) and calls cancel() in 2nd times",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         0,
					initialDuration:       0,
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         10,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      100 * time.Microsecond,
					errLog:                true,
				},
				want: want{
					err: err,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			err := errors.New("erros is occurred")
			cnt := 0
			f := func(context.Context) (interface{}, bool, error) {
				cnt++
				if cnt > 1 {
					time.Sleep(10 * time.Millisecond)
				}
				return str, true, err
			}
			return test{
				name: "return nil response and error when function returns ends due to backoffTimeLimit",
				args: args{
					ctx: ctx,
					f:   f,
				},
				fields: fields{
					backoffFactor:         1.1,
					initialDuration:       float64(time.Millisecond * 5),
					jittedInitialDuration: 0,
					jitterLimit:           0,
					durationLimit:         10,
					maxDuration:           0,
					maxRetryCount:         1,
					backoffTimeLimit:      30 * time.Microsecond,
					errLog:                true,
				},
				want: want{
					err: errors.ErrBackoffTimeout(err),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &backoff{
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

			gotRes, err := b.Do(test.args.ctx, test.args.f)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

func TestMetrics(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	type want struct {
		want map[string]int64
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, map[string]int64) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got map[string]int64) error {
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
		           in0:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           in0:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Metrics(test.args.in0)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
