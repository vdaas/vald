// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrTimeoutParseFailed(t *testing.T) {
	type args struct {
		timeout string
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("invalid timeout value: 10hours\t:timeout parse error out put failed")
			return test{
				name: "return an ErrTimeoutParseFailed error when timeout is not empty.",
				args: args{
					timeout: "10hours",
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("invalid timeout value: \t:timeout parse error out put failed")
			return test{
				name: "return an ErrTimeoutParseFailed error when timeout is empty.",
				args: args{
					timeout: "",
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrTimeoutParseFailed(test.args.timeout)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrServerNotFound(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("server gateway.vald.svc.cluster.local not found")
			return test{
				name: "return an ErrServerNotFound error when the name is not empty.",
				args: args{
					name: "gateway.vald.svc.cluster.local",
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("server  not found")
			return test{
				name: "return an ErrServerNotFound error when the name is empty.",
				args: args{
					name: "",
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrServerNotFound(test.args.name)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrOptionFailed(t *testing.T) {
	type args struct {
		err error
		ref reflect.Value
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("failed to setup option :\tfmt.Println: option failed error")
			return test{
				name: "return an ErrOptionFailed error when err and ref are not empty.",
				args: args{
					err: New("option failed error"),
					ref: func() reflect.Value {
						var i any = fmt.Println
						return reflect.ValueOf(i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to setup option :\t")
			return test{
				name: "return an ErrOptionFailed error when err is empty and ref is zero value.",
				args: args{
					ref: func() reflect.Value {
						var i int
						return reflect.ValueOf(&i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to setup option :\t: option failed error")
			return test{
				name: "return an ErrOptionFailed error when err is not empty and ref is nil.",
				args: args{
					err: New("option failed error"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to setup option :\t")
			return test{
				name: "return an ErrOptionFailed error when err is empty and ref is <invalid reflect.Value>.",
				args: args{
					ref: reflect.Value{},
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrOptionFailed(test.args.err, test.args.ref)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrArgumentPraseFailed(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("argument parse failed: argument parse error")
			return test{
				name: "return an ErrArgumentParseFailed error when err is not empty.",
				args: args{
					err: New("argument parse error"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("argument parse failed")
			return test{
				name: "return an ErrArgumentParseFailed error when err is empty.",
				args: args{},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrArgumentParseFailed(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrBackoffTimeout(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("backoff timeout by limitation: backoff is timeout")
			return test{
				name: "return an ErrBackoffTimeout error when err is not empty.",
				args: args{
					err: New("backoff is timeout"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("backoff timeout by limitation")
			return test{
				name: "return an ErrBackoffTimeout error when err is empty.",
				args: args{},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrBackoffTimeout(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrLoggingRetry(t *testing.T) {
	type args struct {
		err error
		ref reflect.Value
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("failed to output fmt.Println logs, retrying...: logging retry")
			return test{
				name: "return an ErrLoggingRetry error when err and ref are not empty.",
				args: args{
					err: New("logging retry"),
					ref: func() reflect.Value {
						var i any = fmt.Println
						return reflect.ValueOf(i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs, retrying...: logging retry")
			return test{
				name: "return an ErrLoggingRetry error when err is not empty and ref is nil.",
				args: args{
					err: New("logging retry"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs, retrying...")
			return test{
				name: "return an ErrLoggingRetry error when err is empty and ref is zero value.",
				args: args{
					ref: func() reflect.Value {
						var i int
						return reflect.ValueOf(&i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs, retrying...")
			return test{
				name: "return an ErrLoggingRetry error when err is empty and ref is <invalid reflect.Value>.",
				args: args{
					ref: reflect.Value{},
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrLoggingRetry(test.args.err, test.args.ref)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrLoggingFailed(t *testing.T) {
	type args struct {
		err error
		ref reflect.Value
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("failed to output fmt.Println logs: logging retry")
			return test{
				name: "return an ErrLoggingFailed error when err and ref are not empty.",
				args: args{
					err: New("logging retry"),
					ref: func() reflect.Value {
						var i any = fmt.Println
						return reflect.ValueOf(i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs: logging retry")
			return test{
				name: "return an ErrLoggingFailed error when err is not empty and ref is nil.",
				args: args{
					err: New("logging retry"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs")
			return test{
				name: "return an ErrLoggingFailed error when err is empty and ref is zero value.",
				args: args{
					ref: func() reflect.Value {
						var i int
						return reflect.ValueOf(&i)
					}(),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("failed to output  logs")
			return test{
				name: "return an ErrLoggingFailed error when err is empty and ref is <invalid reflect,Value>.",
				args: args{
					ref: reflect.Value{},
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := ErrLoggingFailed(test.args.err, test.args.ref)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		msg string
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("error is occurred")
			return test{
				name: "return a New error when msg is not empty.",
				args: args{
					msg: "error is occurred",
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil when msg is empty.",
				args: args{},
				want: want{},
			}
		}(),
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

			got := New(test.args.msg)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := fmt.Errorf("error is occurred: err")
			return test{
				name: "return an error when err and msg are not empty.",
				args: args{
					err: New("err"),
					msg: "error is occurred",
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("err")
			return test{
				name: "return an error when err is not empty and msg is empty.",
				args: args{
					err: New("err"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			wantErr := New("error is occurred")
			return test{
				name: "return an error when err is empty and msg is not empty.",
				args: args{
					msg: "error is occurred",
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil when err and msg are empty.",
				args: args{},
				want: want{},
			}
		}(),
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

			got := Wrap(test.args.err, test.args.msg)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWrapf(t *testing.T) {
	type args struct {
		err    error
		format string
		args   []any
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			err := New("err: ")
			format := "error is occurred: %v"
			val := []any{
				"timeout error",
			}
			wantErr := fmt.Errorf("%s: %w", fmt.Sprintf(format, val...), err)
			return test{
				name: "return an error when err and format are not empty and args has a single value.",
				args: args{
					err:    err,
					format: format,
					args:   val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			err := New("err: ")
			format := "error is occurred: %v : %v"
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := fmt.Errorf("%s: %w", fmt.Sprintf(format, val...), err)
			return test{
				name: "return an error when err and format are not empty and args has multiple values.",
				args: args{
					err:    err,
					format: format,
					args:   val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			err := New("err: ")
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := err
			return test{
				name: "return an error when err is not empty and format is empty and args has multiple values.",
				args: args{
					err:  err,
					args: val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			err := New("err: ")
			format := "error is occurred: %v : %v"
			wantErr := fmt.Errorf("%s: %w", format, err)
			return test{
				name: "return an error when err and format are not empty and args is empty.",
				args: args{
					err:    err,
					format: format,
				},
				want: want{
					wantErr,
				},
				checkFunc: func(w want, got error) error {
					if got == nil || got.Error() != w.want.Error() || !Is(got, err) {
						return fmt.Errorf("got: \"%w\",\n\t\t\t\twant: \"%w\"", got, w.want)
					}
					return nil
				},
			}
		}(),
		func() test {
			err := New("err: ")
			wantErr := err
			return test{
				name: "return an error when err is not empty and format and args are empty.",
				args: args{
					err: err,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			format := "error is occurred: %v : %v"
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := fmt.Errorf(format, val...)
			return test{
				name: "return an error when err is empty and format and args are not empty.",
				args: args{
					format: format,
					args:   val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			format := "error is occurred: %v : %v"
			wantErr := New(format)
			return test{
				name: "return an error when err and args are empty and format is not empty.",
				args: args{
					format: format,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when all of the input is empty.",
				args: args{},
				want: want{},
			}
		}(),
		func() test {
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := fmt.Errorf("%v %v", val[0], val[1])
			return test{
				name: "return nil when a format is empty and args has multiple values.",
				args: args{
					args: val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			val := []any{
				map[string]int{"invalid time_duration": 10},
			}
			wantErr := fmt.Errorf("%v", val[0])
			return test{
				name: "return an error when a format is empty and args has a single value",
				args: args{
					args: val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
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

			got := Wrapf(test.args.err, test.args.format, test.args.args...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCause(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := Unwrap(Wrap(New("err"), "invalid parameter"))
			return test{
				name: "return an unwrapd error when err is not empty.",
				args: args{
					err: Wrap(New("err"), "invalid parameter"),
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil when err is empty.",
				args: args{},
				want: want{},
			}
		}(),
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

			got := Cause(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUnwarp(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			wantErr := New("err")
			err := fmt.Errorf("%s: %w", "error occurs", wantErr)
			return test{
				name: "return an unwrapped error when err is not empty.",
				args: args{
					err: err,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil when err is empty.",
				args: args{},
				want: want{},
			}
		}(),
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

			got := Unwrap(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		args   []any
	}
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return fmt.Errorf("got: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			format := "error is occurred: %v"
			val := []any{
				"timeout error",
			}
			wantErr := fmt.Errorf(format, val...)
			return test{
				name: "return an error when a format is not empty and args has a single value.",
				args: args{
					format: format,
					args:   val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			format := "error is occurred: %v : %v"
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := fmt.Errorf(format, val...)
			return test{
				name: "return an error when a format is not empty and args has multiple values.",
				args: args{
					format: format,
					args:   val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			val := []any{
				"invalid time_duration",
				10,
			}
			wantErr := fmt.Errorf("%v %v", val[0], val[1])
			return test{
				name: "return an error when a format is empty and args has multiple values.",
				args: args{
					args: val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			val := []any{
				map[string]int{"invalid time_duration": 10},
			}
			wantErr := fmt.Errorf("%v", val[0])
			return test{
				name: "return nil when a format is empty and args has a single value.",
				args: args{
					args: val,
				},
				want: want{
					wantErr,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil when all of the input is empty.",
				args: args{},
				want: want{},
			}
		}(),
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

			got := Errorf(test.args.format, test.args.args...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

type interErr struct {
	msg string
}

type uncomparableErr struct {
	err []interErr
}

func (err uncomparableErr) Error() string {
	var str strings.Builder
	for _, e := range err.err {
		str.WriteString(e.msg)
	}
	return str.String()
}

type wrapErr struct {
	err error
}

func (err wrapErr) Error() string {
	return err.err.Error()
}

func (err wrapErr) Unwrap() error {
	return err.err
}

type isErr struct {
	err error
}

func (err isErr) Error() string {
	return err.err.Error()
}

func (err isErr) Is(e error) bool {
	return err.err.Error() == e.Error()
}

func TestIs(t *testing.T) {
	type args struct {
		err    error
		target error
	}
	type want struct {
		want bool
	}
	type test struct {
		args       args
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got bool) error {
		if got != w.want {
			return fmt.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return false when target is nil.",
			args: args{
				err: New("invalid parameter"),
			},
			want: want{},
		},
		{
			name: "return false when err is nil.",
			args: args{
				target: New("invalid parameter"),
			},
			want: want{},
		},
		{
			name: "return true when err is same comparable errors type and same error as target.",
			args: args{
				err:    New("invalid parameter"),
				target: New("invalid parameter"),
			},
			want: want{
				true,
			},
		},
		{
			name: "return false when err is same comparable errors type and differ error as target.",
			args: args{
				err:    New("invalid parameter"),
				target: New("err is occurred"),
			},
			want: want{},
		},
		{
			name: "return true when err is comparable error and target is uncomparable error and both err msg is same.",
			args: args{
				err: New("err is occurred"),
				target: uncomparableErr{
					[]interErr{
						{
							msg: "err is occurred",
						},
					},
				},
			},
			want: want{
				true,
			},
		},
		{
			name: "return false when err is comparable error and target is uncomparable error and both err msg is not same.",
			args: args{
				err: New("err is occurred"),
				target: uncomparableErr{
					[]interErr{
						{
							msg: "invalid parameter",
						},
					},
				},
			},
			want: want{},
		},
		{
			name: "return true when err is wrapped comparable error and target is uncomparable error and err.err.Error() and target msg are same.",
			args: args{
				err: wrapErr{
					err: New("invalid parameter"),
				},
				target: uncomparableErr{
					[]interErr{
						{
							msg: "invalid parameter",
						},
					},
				},
			},
			want: want{
				true,
			},
		},
		{
			name: "return false when err is wrapped comparable error and target is uncomparable error and err.err.Error() and target msg are not same.",
			args: args{
				err: wrapErr{
					err: New("err is occurred"),
				},
				target: uncomparableErr{
					[]interErr{
						{
							msg: "invalid parameter",
						},
					},
				},
			},
			want: want{},
		},
		{
			name: "return false when err is comparable error with Is() implemented and target is uncomparable error and target msg is empty.",
			args: args{
				err: isErr{
					err: New("err is occurred"),
				},
				target: uncomparableErr{},
			},
			want: want{},
		},
		{
			name: "return true when err is comparable error with Is() implemented and target is uncomparable error and target msg is not empty.",
			args: args{
				err: isErr{
					err: New("err is occurred"),
				},
				target: uncomparableErr{
					[]interErr{
						{
							msg: "err is occurred",
						},
					},
				},
			},
			want: want{
				true,
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
			got := Is(test.args.err, test.args.target)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsNot(t *testing.T) {
	type args struct {
		err     error
		targets []error
	}
	type want struct {
		want bool
	}
	type test struct {
		name string
		args args
		want want
	}
	tests := []test{
		{
			name: "return false when err is nil and targets are given, unlike the negation of Is.",
			args: args{
				targets: []error{New("invalid parameter")},
			},
			want: want{},
		},
		{
			name: "return false when err is nil and no targets are given.",
			args: args{},
			want: want{},
		},
		{
			name: "return true when err is not nil and no targets are given.",
			args: args{
				err: New("err is occurred"),
			},
			want: want{
				true,
			},
		},
		{
			name: "return false when err matches one of the targets.",
			args: args{
				err:     New("invalid parameter"),
				targets: []error{New("err is occurred"), New("invalid parameter")},
			},
			want: want{},
		},
		{
			name: "return true when err matches none of the targets.",
			args: args{
				err:     New("invalid parameter"),
				targets: []error{New("err is occurred"), New("connection refused")},
			},
			want: want{
				true,
			},
		},
		{
			name: "return true when err is not nil and all targets are nil.",
			args: args{
				err:     New("invalid parameter"),
				targets: []error{nil, nil},
			},
			want: want{
				true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			got := IsNot(test.args.err, test.args.targets...)
			if got != test.want.want {
				tt.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, test.want.want)
			}
		})
	}
}

func TestAs(t *testing.T) {
	type args struct {
		err    error
		target any
	}
	type want struct {
		want bool
	}
	type test struct {
		args       args
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got bool) error {
		if got != w.want {
			return fmt.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true when err and target is not empty.",
			args: args{
				err:    New("err"),
				target: New("err is occurred"),
			},
			want: want{
				true,
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

			got := As(test.args.err, &test.args.target)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "succeeds to remove duplicated errors",
			args: args{
				errs: []error{
					New("same error1"),
					New("same error1"),
					New("same error2"),
					New("same error2"),
					New("same error2"),
					New("same error3"),
				},
			},
			want: []error{
				New("same error1"),
				New("same error2"),
				New("same error3"),
			},
		},
		{
			name: "single error remains the same",
			args: args{
				errs: []error{
					New("same error"),
				},
			},
			want: []error{
				New("same error"),
			},
		},
		{
			name: "empty errs remains the same",
			args: args{
				errs: []error{},
			},
			want: []error{},
		},
	}

	equalErrs := func(errs1, errs2 []error) bool {
		if len(errs1) != len(errs2) {
			return false
		}
		for i := range errs1 {
			if !Is(errs1[i], errs2[i]) {
				return false
			}
		}
		return true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.errs); !equalErrs(got, tt.want) {
				t.Errorf("removeDuplicatedErrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		errs []error
	}
	type want struct {
		err error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, err error) error {
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when all errors are nil",
			args: args{
				errs: []error{
					nil, nil, nil,
				},
			},
		},
		{
			name: "returns an aggregated error when all errors are non-nil and different",
			args: args{
				errs: []error{
					New("error1"), New("error2"), New("error3"),
				},
			},
			want: want{
				err: &joinError{
					errs: []error{
						New("error1"), New("error2"), New("error3"),
					},
				},
			},
		},
		{
			name: "returns an error when errors are mixed nil and non-nil",
			args: args{
				errs: []error{
					nil, New("error1"), nil,
				},
			},
			want: want{
				err: New("error1"),
			},
		},
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

			err := Join(test.args.errs...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestIsAny(t *testing.T) {
// 	type args struct {
// 		err     error
// 		targets []error
// 	}
// 	type want struct {
// 		wantSame bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotSame bool) error {
// 		if !reflect.DeepEqual(gotSame, w.wantSame) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSame, w.wantSame)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		           targets:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           err:nil,
// 		           targets:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotSame := IsAny(test.args.err, test.args.targets...)
// 			if err := checkFunc(test.want, gotSame); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestIsNot(t *testing.T) {
// 	type args struct {
// 		err     error
// 		targets []error
// 	}
// 	type want struct {
// 		wantSame bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotSame bool) error {
// 		if !reflect.DeepEqual(gotSame, w.wantSame) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSame, w.wantSame)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		           targets:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           err:nil,
// 		           targets:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotSame := IsNot(test.args.err, test.args.targets...)
// 			if err := checkFunc(test.want, gotSame); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_is(t *testing.T) {
// 	type args struct {
// 		err    error
// 		target error
// 	}
// 	type want struct {
// 		wantSame bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotSame bool) error {
// 		if !reflect.DeepEqual(gotSame, w.wantSame) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSame, w.wantSame)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		           target:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           err:nil,
// 		           target:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotSame := is(test.args.err, test.args.target)
// 			if err := checkFunc(test.want, gotSame); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestUnwrap(t *testing.T) {
// 	type args struct {
// 		err error
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !Is(err, w.err) {
// 			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           err:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := Unwrap(test.args.err)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_joinError_Error(t *testing.T) {
// 	type fields struct {
// 		errs []error
// 	}
// 	type want struct {
// 		wantStr string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotStr string) error {
// 		if !reflect.DeepEqual(gotStr, w.wantStr) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			e := &joinError{
// 				errs: test.fields.errs,
// 			}
//
// 			gotStr := e.Error()
// 			if err := checkFunc(test.want, gotStr); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_joinError_String(t *testing.T) {
// 	type fields struct {
// 		errs []error
// 	}
// 	type want struct {
// 		wantStr string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotStr string) error {
// 		if !reflect.DeepEqual(gotStr, w.wantStr) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			e := &joinError{
// 				errs: test.fields.errs,
// 			}
//
// 			gotStr := e.String()
// 			if err := checkFunc(test.want, gotStr); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_joinError_Unwrap(t *testing.T) {
// 	type fields struct {
// 		errs []error
// 	}
// 	type want struct {
// 		want []error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []error) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			e := &joinError{
// 				errs: test.fields.errs,
// 			}
//
// 			got := e.Unwrap()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// multiUnwrapError is a minimal error exposing Unwrap() []error, which lets the
// tests place nil elements in the multi-unwrap chain.
type multiUnwrapError struct {
	msg  string
	errs []error
}

func (e *multiUnwrapError) Error() string { return e.msg }

func (e *multiUnwrapError) Unwrap() []error { return e.errs }

func TestIs_multiUnwrapNilElement(t *testing.T) {
	t.Parallel()
	const multiMsg = "multi"
	hit := New("hit")
	other := New("other")
	missing := New("missing")
	tests := []struct {
		err    error
		target error
		name   string
		want   bool
	}{
		{
			name:   "nil first element is skipped and later element matches",
			err:    &multiUnwrapError{msg: multiMsg, errs: []error{nil, hit}},
			target: hit,
			want:   true,
		},
		{
			name:   "nil first element is skipped and nothing matches",
			err:    &multiUnwrapError{msg: multiMsg, errs: []error{nil, other}},
			target: missing,
			want:   false,
		},
		{
			name:   "nil last element does not panic when nothing matches",
			err:    &multiUnwrapError{msg: multiMsg, errs: []error{other, nil}},
			target: missing,
			want:   false,
		},
		{
			name:   "nil last element does not hide an earlier match",
			err:    &multiUnwrapError{msg: multiMsg, errs: []error{hit, nil}},
			target: hit,
			want:   true,
		},
		{
			name:   "only nil elements",
			err:    &multiUnwrapError{msg: multiMsg, errs: []error{nil, nil}},
			target: missing,
			want:   false,
		},
		{
			name:   "empty element list",
			err:    &multiUnwrapError{msg: multiMsg},
			target: missing,
			want:   false,
		},
		{
			name:   "the multi error itself matches by message",
			err:    &multiUnwrapError{msg: "missing", errs: []error{nil}},
			target: missing,
			want:   true,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if got := Is(test.err, test.target); got != test.want {
				tt.Errorf("Is() = %v, want %v", got, test.want)
			}
		})
	}
}

type callCountError struct {
	next  error
	calls *int
	msg   string
}

func (e *callCountError) Error() string {
	*e.calls++
	return e.msg
}

func (e *callCountError) Unwrap() error { return e.next }

func TestIs_errorCallCount(t *testing.T) {
	t.Parallel()
	type want struct {
		layerCalls  []int
		targetCalls int
		same        bool
	}
	tests := []struct {
		name   string
		want   want
		layers int
		match  bool
		ident  bool
	}{
		{
			name:   "mismatch calls Error once per layer and once for target",
			layers: 3,
			want:   want{same: false, targetCalls: 1, layerCalls: []int{1, 1, 1}},
		},
		{
			name:   "message match at innermost layer",
			layers: 3,
			match:  true,
			want:   want{same: true, targetCalls: 1, layerCalls: []int{1, 1, 1}},
		},
		{
			name:   "identity match does not call Error",
			layers: 1,
			ident:  true,
			want:   want{same: true, targetCalls: 0, layerCalls: []int{0}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var targetCalls int
			target := &callCountError{msg: "target", calls: &targetCalls}
			counts := make([]int, tc.layers)
			var err error
			for i := tc.layers - 1; i >= 0; i-- {
				msg := fmt.Sprintf("layer-%d", i)
				if i == tc.layers-1 && tc.match {
					msg = "TARGET"
				}
				err = &callCountError{msg: msg, calls: &counts[i], next: err}
			}
			if tc.ident {
				err = target
				counts[0] = 0
			}
			if got := Is(err, target); got != tc.want.same {
				t.Fatalf("Is() = %v, want %v", got, tc.want.same)
			}
			if targetCalls != tc.want.targetCalls {
				t.Errorf("target.Error() calls = %d, want %d", targetCalls, tc.want.targetCalls)
			}
			if !tc.ident && !reflect.DeepEqual(counts, tc.want.layerCalls) {
				t.Errorf("layer Error() calls = %v, want %v", counts, tc.want.layerCalls)
			}
		})
	}
}

// TestJoin_dedupToSingle verifies that the variadic Join returns the remaining
// error itself, not a one-element join, when duplicates and nils collapse to one.
func TestJoin_dedupToSingle(t *testing.T) {
	t.Parallel()
	e := New("same")
	tests := []struct {
		name string
		errs []error
	}{
		{name: "identical errors", errs: []error{e, e, e}},
		{name: "equal messages", errs: []error{New("same"), New("same"), New("same")}},
		{name: "duplicates and nils", errs: []error{nil, e, nil, e}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Join(tc.errs...)
			if _, ok := got.(*joinError); ok {
				t.Fatalf("Join() = %#v, want a single non-join error", got)
			}
			if got == nil || got.Error() != "same" {
				t.Errorf("Join() = %v, want error \"same\"", got)
			}
		})
	}
}

// TestJoin_uncomparable verifies that the two-argument Join does not panic
// when both errors share an uncomparable dynamic type.
func TestJoin_uncomparable(t *testing.T) {
	t.Parallel()
	const msgDup = "dup"
	tests := []struct {
		a, b     error
		name     string
		wantMsg  string
		wantJoin bool
	}{
		{
			name:    "same message collapses to the first error",
			a:       uncomparableErr{err: []interErr{{msg: msgDup}}},
			b:       uncomparableErr{err: []interErr{{msg: msgDup}}},
			wantMsg: msgDup,
		},
		{
			name:     "different messages are joined",
			a:        uncomparableErr{err: []interErr{{msg: "a"}}},
			b:        uncomparableErr{err: []interErr{{msg: "b"}}},
			wantMsg:  "a\nb",
			wantJoin: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Join() panicked: %v", r)
				}
			}()
			got := Join(tc.a, tc.b)
			if got == nil {
				t.Fatal("Join() = nil, want non-nil")
			}
			if _, ok := got.(*joinError); ok != tc.wantJoin {
				t.Errorf("Join() is *joinError = %v, want %v", ok, tc.wantJoin)
			}
			if got.Error() != tc.wantMsg {
				t.Errorf("Join().Error() = %q, want %q", got.Error(), tc.wantMsg)
			}
		})
	}
}

// TestWrap_errorCallCount verifies that Wrap does not evaluate err.Error()
// redundantly before Is, which already compares the messages.
func TestWrap_errorCallCount(t *testing.T) {
	t.Parallel()
	const msgEqual = "equal"
	tests := []struct {
		name     string
		errMsg   string
		msg      string
		want     string
		maxCalls int
	}{
		{
			// one call inside Is plus one when fmt formats the %w operand.
			name:     "different message is wrapped",
			errMsg:   "inner",
			msg:      "outer",
			want:     "outer: inner",
			maxCalls: 2,
		},
		{
			name:     "equal message returns err",
			errMsg:   msgEqual,
			msg:      msgEqual,
			want:     msgEqual,
			maxCalls: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var calls int
			err := &callCountError{msg: tc.errMsg, calls: &calls}
			got := Wrap(err, tc.msg)
			if calls > tc.maxCalls {
				t.Errorf("err.Error() calls = %d, want <= %d", calls, tc.maxCalls)
			}
			if got == nil || got.Error() != tc.want {
				t.Fatalf("Wrap() = %v, want %q", got, tc.want)
			}
			// err has no inner error, so only a wrapper returned by Wrap unwraps to non-nil.
			if tc.errMsg == tc.msg && Unwrap(got) != nil {
				t.Errorf("Wrap() = %#v, want the input error itself", got)
			}
		})
	}
}
