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
			name: "return valid string with no stacktrace initialized",
			beforeFunc: func() {
				infoProvider, _ = New(WithServerName(""), WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}))
			},
			want: want{
				want: "\nbuild cpu info flags -> []\nbuild time           -> \ncgo enabled          -> \ngit commit           -> master\ngo arch              -> " + runtime.GOARCH + "\ngo os                -> " + runtime.GOOS + "\ngo version           -> " + runtime.Version() + "\nngt version          -> \nserver name          -> \nvald version         -> \x1b[1mv0.0.1\x1b[22m",
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
			name: "return detail object",
			beforeFunc: func() {
				infoProvider, _ = New(WithServerName(""), WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}))
			},
			want: want{
				want: Detail{
					Version:           "v0.0.1",
					GitCommit:         "master",
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
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
			name: "set success when all fields are empty",
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
			name: "new returns default info without option",
			args: args{
				opts: nil,
			},
			want: want{
				want: &info{
					detail: Detail{
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
			name: "new returns info with 1 option",
			args: args{
				opts: []Option{
					WithServerName("sn"),
				},
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName: "sn",
						Version:    GitCommit,
						GitCommit:  GitCommit,
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: CGOEnabled,
						//StackTrace:        nil,
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
			name: "new returns info with multiple option",
			args: args{
				opts: []Option{
					WithServerName("sn"),
					func(i *info) error {
						i.detail.Version = "ver"
						return nil
					},
				},
			},
			want: want{
				want: &info{
					detail: Detail{
						ServerName: "sn",
						Version:    "ver",
						GitCommit:  GitCommit,
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: CGOEnabled,
						//StackTrace:        nil,
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
			name: "new log the error when invalid option occurred",
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
						Version:    GitCommit,
						GitCommit:  GitCommit,
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						CGOEnabled: CGOEnabled,
						//StackTrace:        nil,
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
			name: "new return error when criical error occurred",
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
			name: "return correct string when stacktrace is initialized",
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
						StackTrace{
							URL:      "url",
							FuncName: "func",
							File:     "file",
							Line:     10,
						},
					},
				},
			},
			want: want{
				want: "\nbuild cpu info flags -> []\nbuild time           -> bt\ncgo enabled          -> true\ngit commit           -> commit\ngo arch              -> goarch\ngo os                -> goos\ngo version           -> 1.1\nngt version          -> 1.2\nserver name          -> srv\nstack trace-0        -> url\tfunc\nvald version         -> \x1b[1m1.0\x1b[22m",
			},
		},
		{
			name: "return valid string when no stacktrace initialized",
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
				want: "\nbuild cpu info flags -> [avx512f avx512dq]\nbuild time           -> bt\ncgo enabled          -> true\ngit commit           -> commit\ngo arch              -> goarch\ngo os                -> goos\ngo version           -> 1.1\nngt version          -> 1.2\nserver name          -> srv\nvald version         -> \x1b[1m1.0\x1b[22m",
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
			name: "return detail object",
			fields: fields{
				rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				},
			},
			want: want{
				want: Detail{
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
			name: "return detail object when stacktrace is initialized",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/vdaas/vald/tree/master",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "info.go",
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
			name: "return detail object when file name has goroot prefix",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), runtime.GOROOT() + "/src/info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/golang/go/blob/" + runtime.Version() + "/src/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     runtime.GOROOT() + "/src/info.go",
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
			name: "return detail object when go mod path is set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/vdaas/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info.go",
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
			name: "return detail object when go mod path with version is set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/vdaas/blob/v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go",
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
			name: "return detail object when go mod path contains pseudo version",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/vdaas/blob/master/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info.go",
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
			name: "return detail object when go src path is set",
			fields: func() fields {
				i := 0
				return fields{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						if i == 0 {
							i++
							return uintptr(0), "/tmp/go/src/github.com/vdaas/vald/internal/info.go", 100, true
						}
						return uintptr(1), "info.go", 100, false
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
						StackTrace{
							URL:      "https://github.com/vdaas/vald/blob/master/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/src/github.com/vdaas/vald/internal/info.go",
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
			name: "set success when all fields are empty",
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
			name: "GitCommit and Version field is not overwritten when GitCommit field is `internal`",
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
			name: "BuildTime field is not overwritten when BuildTime field is `10`",
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
			name: "GoVersion field is not overwritten when GoVersion field is `1.14`",
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
			name: "GoOS field is not overwritten when GoOS field is `linux`",
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
			name: "GoArch fields is not overwritten when GoArch field is `amd`",
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
			name: "CGOEnabled field is not overwritten when CGOEnabled field is `1`",
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
			name: "NGTVersion field is not overwritten when NGTVersion field is `v1.11.5`",
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
			name: "BuildCPUInfoFlags is not overwritten when BuildCPUInfoFlags field is `test`",
			fields: fields{
				detail: Detail{
					BuildCPUInfoFlags: []string{"test"},
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
							"test",
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
	t.Parallel()
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           URL: "",
		           FuncName: "",
		           File: "",
		           Line: 0,
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
		           URL: "",
		           FuncName: "",
		           File: "",
		           Line: 0,
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
