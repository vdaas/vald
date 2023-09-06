//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

func TestString(t *testing.T) {
	type want struct {
		want *Detail
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got string) error {
		if got == w.want.String() {
			// check the position of "->"
			var oldIdx int
			for i, str := range strings.Split(strings.TrimPrefix(got, "\n"), "\n") {
				idx := strings.Index(str, "->")
				if i != 0 {
					if oldIdx != idx {
						return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
					}
				}
				oldIdx = idx
			}
			return nil
		}
		return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
	}
	tests := []test{
		{
			name: "return correct string with no stack trace initialized",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				infoProvider, _ = New(WithServerName(""),
					WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
						return uintptr(0), "", 0, false
					}))
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: &Detail{
					Version:           "v0.0.1",
					ServerName:        "",
					GitCommit:         GitCommit,
					BuildTime:         "",
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
					GoRoot:            runtime.GOROOT(),
					CGOEnabled:        cgoUnknown,
					NGTVersion:        "",
					BuildCPUInfoFlags: nil,
					StackTrace:        nil,
				},
			},
		},
		{
			name: "return correct string with no information initialized",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				infoProvider = &info{
					rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
						return uintptr(0), "", 0, false
					},
				}
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: &Detail{
					Version:           "",
					ServerName:        "",
					GitCommit:         GitCommit,
					BuildTime:         "",
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
					GoRoot:            runtime.GOROOT(),
					CGOEnabled:        cgoUnknown,
					NGTVersion:        "",
					BuildCPUInfoFlags: nil,
					StackTrace:        nil,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := String()
			if err := checkFunc(test.want, got); err != nil {
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
			beforeFunc: func(t *testing.T) {
				t.Helper()
				infoProvider, _ = New(WithServerName(""), WithRuntimeCaller(func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}))
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				once = sync.Once{}
				infoProvider = nil
			},
			want: want{
				want: Detail{
					ServerName:        "",
					Version:           Version,
					BuildTime:         "",
					GitCommit:         GitCommit,
					GoVersion:         runtime.Version(),
					GoOS:              runtime.GOOS,
					GoArch:            runtime.GOARCH,
					GoRoot:            runtime.GOROOT(),
					CGOEnabled:        cgoUnknown,
					NGTVersion:        "",
					BuildCPUInfoFlags: []string{""},
					StackTrace:        make([]StackTrace, 0, 10),
				},
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Get()
			if err := checkFunc(test.want, got); err != nil {
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Info) error {
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			// skipcq: VET-V0008
			comparator.Comparer(func(x, y sync.Once) bool {
				// skipcq: VET-V0008
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

	gitCommit := GitCommit
	version := Version
	buildTime := BuildTime
	cgoEnabled := CGOEnabled
	ngtVersion := NGTVersion
	buildCPUInfoFlags := BuildCPUInfoFlags
	tests := []test{
		{
			name: "set success when the server name is not empty",
			args: args{
				name: "gateway",
			},
			want: want{
				want: &info{
					baseURL: "https://github.com/vdaas/vald/tree/gitcommit",
					detail: Detail{
						GitCommit:  "gitcommit",
						ServerName: "gateway",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoTrue,
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
			beforeFunc: func(t *testing.T, _ args) {
				t.Helper()
				GitCommit = "gitcommit"
				Version = ""
				BuildTime = "1s"
				CGOEnabled = cgoTrue
				NGTVersion = "v1.11.6"
				BuildCPUInfoFlags = "\t\tavx512f avx512dq\t"
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				once = sync.Once{}
				infoProvider = nil

				GitCommit = gitCommit
				Version = version
				BuildTime = buildTime
				CGOEnabled = cgoEnabled
				NGTVersion = ngtVersion
				BuildCPUInfoFlags = buildCPUInfoFlags
			},
		},
		{
			name: "set success when the server name is an empty string",
			args: args{
				name: "",
			},
			want: want{
				want: &info{
					baseURL: "https://github.com/vdaas/vald/tree/gitcommit",
					detail: Detail{
						GitCommit:  "gitcommit",
						ServerName: "",
						Version:    "gitcommit",
						BuildTime:  "1s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoTrue,
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
			beforeFunc: func(t *testing.T, _ args) {
				t.Helper()
				GitCommit = "gitcommit"
				Version = ""
				BuildTime = "1s"
				CGOEnabled = cgoTrue
				NGTVersion = "v1.11.6"
				BuildCPUInfoFlags = "\t\tavx512f avx512dq\t"
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				once = sync.Once{}
				infoProvider = nil

				GitCommit = gitCommit
				Version = version
				BuildTime = buildTime
				CGOEnabled = cgoEnabled
				NGTVersion = ngtVersion
				BuildCPUInfoFlags = buildCPUInfoFlags
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			Init(test.args.name)
			if err := checkFunc(test.want, infoProvider); err != nil {
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
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Info, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			// skipcq: VET-V0008
			comparator.Comparer(func(x, y sync.Once) bool {
				// skipcq: VET-V0008
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						ServerName:        "",
						Version:           Version,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						GoRoot:            runtime.GOROOT(),
						CGOEnabled:        cgoUnknown,
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						ServerName:        "gateway",
						Version:           Version,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						GoRoot:            runtime.GOROOT(),
						CGOEnabled:        cgoUnknown,
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						ServerName:        "vald",
						Version:           "v1.0.0",
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						GoRoot:            runtime.GOROOT(),
						CGOEnabled:        cgoUnknown,
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						ServerName:        "",
						Version:           Version,
						GitCommit:         GitCommit,
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						GoRoot:            runtime.GOROOT(),
						CGOEnabled:        cgoUnknown,
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_String(t *testing.T) {
	type fields struct {
		detail      Detail
		rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
		rtFuncForPC func(pc uintptr) *runtime.Func
	}
	type want struct {
		want *Detail
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got string) error {
		if got == w.want.String() {
			// check the position of "->"
			var oldIdx int
			for i, str := range strings.Split(strings.TrimPrefix(got, "\n"), "\n") {
				idx := strings.Index(str, "->")
				if i != 0 {
					if oldIdx != idx {
						return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
					}
				}
				oldIdx = idx
			}
			return nil
		}
		return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
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
					GoRoot:            "/usr/local/go",
					CGOEnabled:        cgoTrue,
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
				want: &Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					GoRoot:            "/usr/local/go",
					CGOEnabled:        cgoTrue,
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
					GoRoot:            "/usr/local/go",
					CGOEnabled:        cgoTrue,
					NGTVersion:        "1.2",
					BuildCPUInfoFlags: nil,
					StackTrace:        []StackTrace{},
				},
				rtCaller: func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				},
			},
			want: want{
				want: &Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					GoRoot:            "/usr/local/go",
					CGOEnabled:        cgoTrue,
					NGTVersion:        "1.2",
					BuildCPUInfoFlags: nil,
					StackTrace:        nil,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			i := info{
				detail:      test.fields.detail,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}

			got := i.String()
			if err := checkFunc(test.want, got); err != nil {
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
		want *Detail
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got string) error {
		if got == w.want.String() {
			// check the position of "->"
			var oldIdx int
			for i, str := range strings.Split(strings.TrimPrefix(got, "\n"), "\n") {
				idx := strings.Index(str, "->")
				if i != 0 {
					if oldIdx != idx {
						return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
					}
				}
				oldIdx = idx
			}
			return nil
		}
		return errors.Errorf("\tgot: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
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
				CGOEnabled:        cgoTrue,
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
				want: &Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					CGOEnabled:        cgoTrue,
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
				CGOEnabled:        cgoTrue,
				NGTVersion:        "1.2",
				BuildCPUInfoFlags: nil,
				StackTrace:        []StackTrace{},
			},
			want: want{
				want: &Detail{
					Version:           "1.0",
					ServerName:        "srv",
					GitCommit:         "commit",
					BuildTime:         "bt",
					GoVersion:         "1.1",
					GoOS:              "goos",
					GoArch:            "goarch",
					CGOEnabled:        cgoTrue,
					NGTVersion:        "1.2",
					BuildCPUInfoFlags: nil,
					StackTrace:        nil,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_Get(t *testing.T) {
	type fields struct {
		detail      Detail
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
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got Detail) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%v\",\n\t\t\t\twant: \"%v\"", got, w.want)
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
					ServerName: "",
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/tree/main",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/golang/go/blob/" + runtime.Version() + "/src/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     runtime.GOROOT() + "/src/info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/blob/v0.0.0-20171023180738-a3a6125de932/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/blob/main/vald/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
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
					Version:    "",
					GitCommit:  GitCommit,
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					GoRoot:     runtime.GOROOT(),
					CGOEnabled: cgoUnknown,
					StackTrace: []StackTrace{
						{
							URL:      "https://github.com/vdaas/vald/blob/main/internal/info_test.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.Test_info_Get",
							File:     "/tmp/go/src/github.com/vdaas/vald/internal/info_test.go",
							Line:     100,
						},
					},
					NGTVersion: NGTVersion,
					BuildTime:  BuildTime,
					BuildCPUInfoFlags: func() []string {
						if len(BuildCPUInfoFlags) == 0 {
							return nil
						}
						return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
					}(),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			i := info{
				detail:      test.fields.detail,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}

			got := i.Get()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_prepare(t *testing.T) {
	type fields struct {
		detail      Detail
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
		afterFunc  func(*testing.T)
	}
	// skipcq: VET-V0008
	defaultCheckFunc := func(got info, w want) error {
		opts := []comparator.Option{
			comparator.AllowUnexported(info{}),
			comparator.IgnoreFields(info{}, "prepOnce"),
		}
		// skipcq: VET-V0008
		if diff := comparator.Diff(w.want, got, opts...); len(diff) != 0 {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []*test{
		{
			name: "set success with all fields are empty",
			want: want{
				want: info{
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/internal",
					detail: Detail{
						GitCommit:  "internal",
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "v1.0.0",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  "10s",
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  "1.14",
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       "linux",
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     "amd",
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
					},
				},
			},
		},
		{
			name: "set success with CGOEnabled set as true",
			fields: fields{
				detail: Detail{
					CGOEnabled: "1",
				},
			},
			want: want{
				want: info{
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoTrue,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
					},
				},
			},
		},
		{
			name: "set success with CGOEnabled set as false",
			fields: fields{
				detail: Detail{
					CGOEnabled: "0",
				},
			},
			want: want{
				want: info{
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoFalse,
						NGTVersion: NGTVersion,
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:  GitCommit,
						Version:    "",
						BuildTime:  BuildTime,
						GoVersion:  runtime.Version(),
						GoOS:       runtime.GOOS,
						GoArch:     runtime.GOARCH,
						GoRoot:     runtime.GOROOT(),
						CGOEnabled: cgoUnknown,
						NGTVersion: "v1.11.5",
						BuildCPUInfoFlags: func() []string {
							if len(BuildCPUInfoFlags) == 0 {
								return nil
							}
							return strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
						}(),
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
					baseURL: "https://github.com/vdaas/vald/tree/main",
					detail: Detail{
						GitCommit:         GitCommit,
						Version:           "",
						BuildTime:         BuildTime,
						GoVersion:         runtime.Version(),
						GoOS:              runtime.GOOS,
						GoArch:            runtime.GOARCH,
						GoRoot:            runtime.GOROOT(),
						CGOEnabled:        cgoUnknown,
						NGTVersion:        NGTVersion,
						BuildCPUInfoFlags: []string{"avx512f"},
					},
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			i := &info{
				detail:      test.fields.detail,
				rtCaller:    test.fields.rtCaller,
				rtFuncForPC: test.fields.rtFuncForPC,
			}
			i.prepare()
			if err := checkFunc(*i, test.want); err != nil {
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
		afterFunc  func(*testing.T)
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			s := StackTrace{
				URL:      test.fields.URL,
				FuncName: test.fields.FuncName,
				File:     test.fields.File,
				Line:     test.fields.Line,
			}

			got := s.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_info_getDetail(t *testing.T) {
// 	type fields struct {
// 		baseURL     string
// 		detail      Detail
// 		rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
// 		rtFuncForPC func(pc uintptr) *runtime.Func
// 	}
// 	type want struct {
// 		want Detail
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Detail) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got Detail) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           baseURL:"",
// 		           detail:Detail{},
// 		           rtCaller:nil,
// 		           rtFuncForPC:nil,
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
// 		           baseURL:"",
// 		           detail:Detail{},
// 		           rtCaller:nil,
// 		           rtFuncForPC:nil,
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
// 			i := info{
// 				baseURL:     test.fields.baseURL,
// 				detail:      test.fields.detail,
// 				rtCaller:    test.fields.rtCaller,
// 				rtFuncForPC: test.fields.rtFuncForPC,
// 			}
//
// 			got := i.getDetail()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
