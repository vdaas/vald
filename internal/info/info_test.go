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

// Package info provides build-time info
package info

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

func TestString(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if got != w.want {
			return errors.Errorf("\tgot: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return correct string with no stack trace initialized",
			beforeFunc: func() {
				infoProvider, _ = New(WithServerName(""),
					WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
						return uintptr(0), "", 0, false
					}))
			},
			afterFunc: func() {
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[]\ngit commit           ->\tmaster\ngo arch              ->\t" + runtime.GOARCH + "\ngo os                ->\t" + runtime.GOOS + "\ngo version           ->\t" + runtime.Version() + "\nvald version         ->\t\x1b[1mv0.0.1\x1b[22m",
			},
		},

		{
			name: "return correct string with no information initialized",
			beforeFunc: func() {
				infoProvider = &info{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						return uintptr(0), "", 0, false
					},
				}
			},
			afterFunc: func() {
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[]\ngit commit           ->\tmaster\ngo arch              ->\t" + runtime.GOARCH + "\ngo os                ->\t" + runtime.GOOS + "\ngo version           ->\t" + runtime.Version() + "\nvald version         ->\t\x1b[1m\x1b[22m",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type want struct {
		want Detail
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, Detail) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Detail) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return detail with initialized runtime information",
			beforeFunc: func() {
				infoProvider, _ = New(WithServerName(""), WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}))
			},
			afterFunc: func() {
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: Detail{
					ServerName:        "",
					Version:           "v0.0.1",
					BuildTime:         "",
					GitCommit:         "master",
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
					CGOEnabled:        "",
					NGTVersion:        "",
					BuildCPUInfoFlags: []string{""},
					StackTrace:        make([]StackTrace, 0, 10),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Get()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		want Info
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Info) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Info) error {
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.Comparer(func(x, y sync.Once) bool {
				return reflect.DeepEqual(x, y)
			}),
			comparator.Comparer(func(x, y func(skip int) (pc uintptr, file string, line int, ok bool)) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
			comparator.Comparer(func(x, y func(pc uintptr) *runtime.Func) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
		}
		if diff := comparator.Diff(w.want, got, opts...); len(diff) != 0 {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when the server name is not empty",
			args: args{
				name: "gateway",
			},
			want: want{
				want: &info{
					detail: Detail{
						GitCommit:  "gitcommit",
						ServerName: "gateway",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
						StackTrace: nil,
					},
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
				},
			},
			beforeFunc: func(args) {
				GitCommit = "gitcommit"
				Version = ""
				BuildTime = "1s"
				CGOEnabled = "true"
				NGTVersion = "v1.11.6"
				BuildCPUInfoFlags = "\t\tavx512f avx512dq\t"
			},
			afterFunc: func(args) {
				once = sync.Once{}
				infoProvider = nil
			},
		},
		{
			name: "set success when the server name is an empty string",
			args: args{
				name: "",
			},
			want: want{
				want: &info{
					detail: Detail{
						GitCommit:  "gitcommit",
						ServerName: "",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
						StackTrace: nil,
					},
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
				},
			},
			beforeFunc: func(args) {
				GitCommit = "gitcommit"
				Version = ""
				BuildTime = "1s"
				CGOEnabled = "true"
				NGTVersion = "v1.11.6"
				BuildCPUInfoFlags = "\t\tavx512f avx512dq\t"
			},
			afterFunc: func(args) {
				once = sync.Once{}
				infoProvider = nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			Init(test.args.name)
			if err := test.checkFunc(test.want, infoProvider); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Info
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Info, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Info, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.Comparer(func(x, y sync.Once) bool {
				return reflect.DeepEqual(x, y)
			}),
			comparator.Comparer(func(x, y func(skip int) (pc uintptr, file string, line int, ok bool)) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
			comparator.Comparer(func(x, y func(pc uintptr) *runtime.Func) bool {
				return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
			}),
		}
		if diff := comparator.Diff(got, w.want, opts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "return default info with no option set",
			args: args{
				opts: nil,
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName:        "",
						Version:           GitCommit,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						CGOEnabled:        CGOEnabled,
						NGTVersion:        NGTVersion,
						BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
						StackTrace:        nil,
					},
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
				},
			},
		},
		{
			name: "return info when 1 option set",
			args: args{
				opts: []Option{
					WithServerName("gateway"),
				},
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName:        "gateway",
						Version:           GitCommit,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						CGOEnabled:        CGOEnabled,
						NGTVersion:        NGTVersion,
						BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
					},
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
				},
			},
		},
		{
			name: "return info when multiple options set",
			args: args{
				opts: []Option{
					WithServerName("vald"),
					func(i *info) error {
						i.detail.Version = "v1.0.0"
						return nil
					},
				},
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName:        "vald",
						Version:           "v1.0.0",
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						CGOEnabled:        CGOEnabled,
						NGTVersion:        NGTVersion,
						BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
						StackTrace:        nil,
					},
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
				},
			},
		},
		{
			name: "return info and log the error when an invalid option set",
			args: args{
				opts: []Option{
					func(i *info) error {
						return errors.NewErrInvalidOption("field", "err")
					},
				},
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName:        "",
						Version:           GitCommit,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						CGOEnabled:        CGOEnabled,
						NGTVersion:        NGTVersion,
						BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
						StackTrace:        nil,
					},
					prepOnce: func() (o sync.Once) {
						o.Do(func() {})
						return
					}(),
					rtCaller:    runtime.Caller,
					rtFuncForPC: runtime.FuncForPC,
				},
			},
		},
		{
			name: "return an error when a critical error occurred",
			args: args{
				opts: []Option{
					func(i *info) error {
						return errors.NewErrCriticalOption("field", "err")
					},
				},
			},
			want: want{
				err: errors.NewErrCriticalOption("field", "err"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_info_String(t *testing.T) {
	type fields struct {
		detail      Detail
		prepOnce    sync.Once
		rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
		rtFuncForPC func(pc uintptr) *runtime.Func
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return string with stack trace initialized",
			fields: fields{
				detail: Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					CGOEnabled:        "true",
					NGTVersion:        "1.2",
					BuildCPUInfoFlags: nil,
					StackTrace: []StackTrace{
						{
							URL:      "url",
							FuncName: "func",
							File:     "file",
							Line:     10,
						},
					},
				},
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[]\nbuild time           ->\tbt\ncgo enabled          ->\ttrue\ngit commit           ->\tcommit\ngo arch              ->\tgoarch\ngo os                ->\tgoos\ngo version           ->\t1.1\nngt version          ->\t1.2\nserver name          ->\tsrv\nstack trace-0        ->\turl\tfunc\nvald version         ->\t\x1b[1m1.0\x1b[22m",
			},
		},
		{
			name: "return string with no stack trace initialized",
			fields: fields{
				detail: Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					CGOEnabled:        "true",
					NGTVersion:        "1.2",
					BuildCPUInfoFlags: nil,
					StackTrace:        []StackTrace{},
				},
				rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				},
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[avx512f avx512dq]\nbuild time           ->\tbt\ncgo enabled          ->\ttrue\ngit commit           ->\tcommit\ngo arch              ->\tgoarch\ngo os                ->\tgoos\ngo version           ->\t1.1\nngt version          ->\t1.2\nserver name          ->\tsrv\nvald version         ->\t\x1b[1m1.0\x1b[22m",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := info{
				detail:      test.fields.detail,
				prepOnce:    test.fields.prepOnce,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}

			got := i.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDetail_String(t *testing.T) {
	type fields struct {
		Version           string
		ServerName        string
		GitCommit         string
		BuildTime         string
		GoVersion         string
		GoOS              string
		GoArch            string
		CGOEnabled        string
		NGTVersion        string
		BuildCPUInfoFlags []string
		StackTrace        []StackTrace
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return string with stack trace initialized",
			fields: fields{
				Version:           "1.0",
				ServerName:        "srv",
				GitCommit:         "commit",
				BuildTime:         "bt",
				GoVersion:         "1.1",
				GoOS:              "goos",
				GoArch:            "goarch",
				CGOEnabled:        "true",
				NGTVersion:        "1.2",
				BuildCPUInfoFlags: nil,
				StackTrace: []StackTrace{
					{
						URL:      "url",
						FuncName: "func",
						File:     "file",
						Line:     10,
					},
				},
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[]\nbuild time           ->\tbt\ncgo enabled          ->\ttrue\ngit commit           ->\tcommit\ngo arch              ->\tgoarch\ngo os                ->\tgoos\ngo version           ->\t1.1\nngt version          ->\t1.2\nserver name          ->\tsrv\nstack trace-0        ->\turl\tfunc\nvald version         ->\t\x1b[1m1.0\x1b[22m",
			},
		},
		{
			name: "return string with no stack trace initialized",
			fields: fields{
				Version:           "1.0",
				ServerName:        "srv",
				GitCommit:         "commit",
				BuildTime:         "bt",
				GoVersion:         "1.1",
				GoOS:              "goos",
				GoArch:            "goarch",
				CGOEnabled:        "true",
				NGTVersion:        "1.2",
				BuildCPUInfoFlags: nil,
				StackTrace:        []StackTrace{},
			},
			want: want{
				want: "\nbuild cpu info flags ->\t[]\nbuild time           ->\tbt\ncgo enabled          ->\ttrue\ngit commit           ->\tcommit\ngo arch              ->\tgoarch\ngo os                ->\tgoos\ngo version           ->\t1.1\nngt version          ->\t1.2\nserver name          ->\tsrv\nvald version         ->\t\x1b[1m1.0\x1b[22m",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			d := Detail{
				Version:           test.fields.Version,
				ServerName:        test.fields.ServerName,
				GitCommit:         test.fields.GitCommit,
				BuildTime:         test.fields.BuildTime,
				GoVersion:         test.fields.GoVersion,
				GoOS:              test.fields.GoOS,
				GoArch:            test.fields.GoArch,
				CGOEnabled:        test.fields.CGOEnabled,
				NGTVersion:        test.fields.NGTVersion,
				BuildCPUInfoFlags: test.fields.BuildCPUInfoFlags,
				StackTrace:        test.fields.StackTrace,
			}

			got := d.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_Get(t *testing.T) {
	type fields struct {
		detail      Detail
		prepOnce    sync.Once
		rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
		rtFuncForPC func(pc uintptr) *runtime.Func
	}
	type want struct {
		want Detail
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, Detail) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Detail) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return detail object with no stack trace",
			fields: fields{
				rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				},
			},
			want: want{
				want: Detail{
					ServerName:        "",
					Version:           GitCommit,
					GitCommit:         "master",
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
					CGOEnabled:        CGOEnabled,
					StackTrace:        []StackTrace{},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with stack trace initialized",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					ServerName: "",
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/tree/master",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with the file name has goroot prefix",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), runtime.GOROOT() + "/src/info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					ServerName: "",
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/golang/go/blob/" + runtime.Version() + "/src/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     runtime.GOROOT() + "/src/info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with the go mod path set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with the go mod path with version set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/blob/v0.0.0-20171023180738-a3a6125de932/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with the go mod path contains pseudo version",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/blob/master/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
		{
			name: "return detail object with the go src path set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/src/github.com/vdaas/vald/internal/info_test.go", 100, true
						}
						return uintptr(1), "info_test.go", 100, false
					},
					rtFuncForPC: func(ptr uintptr) *runtime.Func {
						return runtime.FuncForPC(reflect.ValueOf(Test_info_Get).Pointer())
					},
				}
			}(),
			want: want{
				want: Detail{
					Version:    GitCommit,
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					CGOEnabled: CGOEnabled,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/blob/master/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/src/github.com/vdaas/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion:        NGTVersion,
					BuildTime:         BuildTime,
					BuildCPUInfoFlags: []string{"avx512f", "avx512dq"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := info{
				detail:      test.fields.detail,
				prepOnce:    test.fields.prepOnce,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}

			got := i.Get()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_prepare(t *testing.T) {
	type fields struct {
		detail      Detail
		prepOnce    sync.Once
		rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
		rtFuncForPC func(pc uintptr) *runtime.Func
	}
	type want struct {
		want info
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(info, want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(got info, w want) error {
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.IgnoreFields(info{}, "prepOnce"),
		}
		if diff := comparator.Diff(w.want, got, opts...); len(diff) != 0 {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	defaultBeforeFunc := func() {
		Version = ""
		GitCommit = "gitcommit"
		BuildTime = "1s"
		CGOEnabled = "true"
		NGTVersion = "v1.11.6"
		BuildCPUInfoFlags = "\t\tavx512f avx512dq\t"
	}
	tests := []test{
		{
			name: "set success with all fields are empty",
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with GitCommit set",
			fields: fields{
				detail: Detail{
					GitCommit: "internal",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "internal",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with Version set",
			fields: fields{
				detail: Detail{
					Version: "v1.0.0",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "v1.0.0",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with BuildTime set",
			fields: fields{
				detail: Detail{
					BuildTime: "10s",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "10s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with GoVersion set",
			fields: fields{
				detail: Detail{
					GoVersion: "1.14",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  "1.14",
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with GoOS set",
			fields: fields{
				detail: Detail{
					GoOS: "linux",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       "linux",
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with GoArch set",
			fields: fields{
				detail: Detail{
					GoArch: "amd",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     "amd",
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with CGOEnabled set",
			fields: fields{
				detail: Detail{
					CGOEnabled: "1",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "1",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with NGTVersion set",
			fields: fields{
				detail: Detail{
					NGTVersion: "v1.11.5",
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.5",
						BuildCPUInfoFlags: []string{
							"avx512f", "avx512dq",
						},
					},
				},
			},
		},
		{
			name: "set success with BuildCPUInfoFlags set",
			fields: fields{
				detail: Detail{
					BuildCPUInfoFlags: []string{"avx512f"},
				},
			},
			want: want{
				want: info{
					detail: Detail{
						GitCommit:  "master",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: "true",
						NGTVersion: "v1.11.6",
						BuildCPUInfoFlags: []string{
							"avx512f",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := &info{
				detail:      test.fields.detail,
				prepOnce:    test.fields.prepOnce,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}

			i.prepare()
			if err := test.checkFunc(*i, test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestStackTrace_String(t *testing.T) {
	type fields struct {
		URL      string
		FuncName string
		File     string
		Line     int
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return stack trace string",
			fields: fields{
				URL:      "https://github.com/golang/go/blob/v1.0.0/internal/info/info_test.go#L40",
				FuncName: "TestStackTrace_String",
				File:     "info_test.go",
				Line:     40,
			},
			want: want{
				want: "URL: https://github.com/golang/go/blob/v1.0.0/internal/info/info_test.go#L40\tFile: info_test.go\tLine: #40\tFuncName: TestStackTrace_String",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := StackTrace{
				URL:      test.fields.URL,
				FuncName: test.fields.FuncName,
				File:     test.fields.File,
				Line:     test.fields.Line,
			}

			got := s.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
