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
package logger

import (
	"os"
	"os/exec"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"google.golang.org/grpc/grpclog"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestInit(t *testing.T) {
	type test struct {
		name       string
		checkFunc  func() error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "set logger success with verbosity level is not set",
			checkFunc: func() error {
				if grpclog.V(1) {
					return errors.New("verbosity level is set")
				}
				return nil
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				once = sync.Once{}
			},
		},
		{
			name: "set logger success with verbosity level is set",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				t.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "2")
			},
			checkFunc: func() error {
				if !grpclog.V(1) {
					return errors.New("verbosity level 1 is not set")
				}
				if !grpclog.V(2) {
					return errors.New("verbosity level is not set")
				}
				if grpclog.V(3) {
					return errors.New("verbosity level is not correct")
				}
				return nil
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				once = sync.Once{}
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			Init()
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Info(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Info success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Info(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Infoln(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Infoln success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Infoln(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Infof(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		args   []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Infof success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Infof(test.args.format, test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warning(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Warning success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Warning(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warningln(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Warningln success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Warningln(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warningf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		args   []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Warningf success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Warningf(test.args.format, test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Error(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Error success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Error(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Errorln(t *testing.T) {
	t.Parallel()
	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Errorln success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Errorln(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Errorf(t *testing.T) {
	t.Parallel()
	type args struct {
		format string
		args   []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Errorf success to log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Errorf(test.args.format, test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=Test_logger_Fatal")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Fatal log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				_ = recover()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Fatal(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatalln(t *testing.T) {
	if os.Getenv("BE_CRASHER") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=Test_logger_Fatalln")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	type args struct {
		args []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Fatalln log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Fatalln(test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatalf(t *testing.T) {
	if os.Getenv("BE_CRASHER") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=Test_logger_Fatalf")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	type args struct {
		format string
		args   []interface{}
	}
	type fields struct {
		v int
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		{
			name: "Fatalf log the message",
			args: args{
				args: []interface{}{"log message"},
			},
			fields: fields{
				v: 0,
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				_ = recover()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			l.Fatalf(test.args.format, test.args.args...)
			if err := test.checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_V(t *testing.T) {
	t.Parallel()
	type args struct {
		v int
	}
	type fields struct {
		v int
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true if v is less than verbosity level",
			args: args{
				v: 3,
			},
			fields: fields{
				v: 5,
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true if v is equal than verbosity level",
			args: args{
				v: 5,
			},
			fields: fields{
				v: 5,
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false if v is larger than verbosity level",
			args: args{
				v: 5,
			},
			fields: fields{
				v: 3,
			},
			want: want{
				want: false,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &logger{
				v: test.fields.v,
			}

			got := l.V(test.args.v)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
