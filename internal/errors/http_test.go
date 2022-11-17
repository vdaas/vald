// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package errors

import (
	"math"
	"testing"
	"time"
)

func TestErrInvalidAPIConfig(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrInvalidAPIConfig error",
			want: want{
				want: New("invalid api config"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidAPIConfig
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidRequest(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrInvalidRequest error",
			want: want{
				want: New("invalid request"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidRequest
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrHandler(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrHandler error when the err is database error",
			args: args{
				err: New("database error"),
			},
			want: want{
				want: New("handler returned error: database error"),
			},
		},
		{
			name: "return ErrHandler error when the err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("handler returned error"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrHandler(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrHandlerTimeout(t *testing.T) {
	type args struct {
		err error
		dur time.Duration
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrHandlerTimeout error when the err is database error and the dur is the minimum number of int64",
			args: args{
				err: New("database error"),
				dur: math.MinInt64,
			},
			want: want{
				want: Errorf("handler timeout %s: database error", time.Duration(math.MinInt64).String()),
			},
		},
		{
			name: "return wrapped ErrHandlerTimeout error when the err is database error and the dur is the 'MinInt64+1'",
			args: args{
				err: New("database error"),
				dur: math.MinInt64 + 1,
			},
			want: want{
				want: Errorf("handler timeout %s: database error", time.Duration(math.MinInt64+1).String()),
			},
		},
		{
			name: "return wrapped ErrHandlerTimeout error when the err is database error and the dur is '0'",
			args: args{
				err: New("database error"),
				dur: 0,
			},
			want: want{
				want: New("handler timeout 0s: database error"),
			},
		},
		{
			name: "return wrapped ErrHandlerTimeout error when the err is database error and the dur is '100'",
			args: args{
				err: New("database error"),
				dur: 100,
			},
			want: want{
				want: New("handler timeout 100ns: database error"),
			},
		},
		{
			name: "return wrapped ErrHandlerTimeout error when the err is database error and the dur is the maximum number of int64",
			args: args{
				err: New("database error"),
				dur: math.MaxInt64,
			},
			want: want{
				want: Errorf("handler timeout %s: database error", time.Duration(math.MaxInt64)),
			},
		},
		{
			name: "return ErrHandlerTimeout error when the err is nil and the dur is '100'",
			args: args{
				err: nil,
				dur: 100,
			},
			want: want{
				want: New("handler timeout 100ns"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrHandlerTimeout(test.args.err, test.args.dur)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRequestBodyCloseAndFlush(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrRequestBodyCloseAndFlush error when the err is close error",
			args: args{
				err: New("close error"),
			},
			want: want{
				want: New("request body flush & close failed: close error"),
			},
		},
		{
			name: "return ErrRequestBodyCloseAndFlush error when the err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("request body flush & close failed"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRequestBodyCloseAndFlush(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRequestBodyClose(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrRequestBodyClose error when the err is close error",
			args: args{
				err: New("close error"),
			},
			want: want{
				want: New("request body close failed: close error"),
			},
		},
		{
			name: "return ErrRequestBodyClose error when the err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("request body close failed"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRequestBodyClose(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRequestBodyFlush(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return wrapped ErrRequestBodyFlush error when the err is close error",
			args: args{
				err: New("close error"),
			},
			want: want{
				want: New("request body flush failed: close error"),
			},
		},
		{
			name: "return ErrRequestBodyFlush error when the err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("request body flush failed"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRequestBodyFlush(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrTransportRetryable(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrTransportRetryable error",
			want: want{
				want: New("transport is retryable"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrTransportRetryable
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
