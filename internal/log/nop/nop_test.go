// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package nop

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		want logger.Logger
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, logger.Logger) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got logger.Logger) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return logger.Logger",
			want: want{
				want: new(nopLogger),
			},
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

			got := New()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Debug(t *testing.T) {
	t.Parallel()
	type args struct {
		vals []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Debug(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Debugf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Debugf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Debugd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Debugd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Info(t *testing.T) {
	t.Parallel()
	type args struct {
		vals []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Info(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Infof(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Infof(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Infod(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Infod(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Warn(t *testing.T) {
	t.Parallel()
	type args struct {
		vals []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Warn(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Warnf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Warnf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Warnd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Warnd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Error(t *testing.T) {
	t.Parallel()
	type args struct {
		vals []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Error(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Errorf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Errorf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Errord(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Errord(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Fatal(t *testing.T) {
	t.Parallel()
	type args struct {
		vals []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Fatal(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Fatalf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Fatalf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Fatald(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		n          *nopLogger
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "successful call",
		},
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
			n := &nopLogger{}

			n.Fatald(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_nopLogger_Close(t *testing.T) {
	t.Parallel()
	type want struct {
		err error
	}
	type test struct {
		name       string
		n          *nopLogger
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "successful call",
			want: want{
				err: nil,
			},
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
			n := &nopLogger{}

			err := n.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
