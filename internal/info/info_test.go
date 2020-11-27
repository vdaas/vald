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
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

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
			name: "string return valid string with no stacktrace initialized",
			beforeFunc: func() {
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
			},
			want: want{
				want: "\nbuild cpu info flags -> []\nbuild time           -> \ncgo enabled          -> \ngit commit           -> master\ngo arch              -> " + runtime.GOARCH + "\ngo os                -> " + runtime.GOOS + "\ngo version           -> " + runtime.Version() + "\nngt version          -> \nserver name          -> \nvald version         -> \x1b[1m\x1b[22m",
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
			name: "get return detail object",
			beforeFunc: func() {
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
			},
			want: want{
				want: Detail{
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					StackTrace: []StackTrace{},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
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
		PrepOnce          sync.Once
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
			name: "string return correct string",
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
					StackTrace{
						URL:      "url",
						FuncName: "func",
						File:     "file",
						Line:     10,
					},
				},
				PrepOnce: sync.Once{},
			},
			want: want{
				want: "\nbuild cpu info flags -> []\nbuild time           -> bt\ncgo enabled          -> true\ngit commit           -> commit\ngo arch              -> goarch\ngo os                -> goos\ngo version           -> 1.1\nngt version          -> 1.2\nserver name          -> srv\nstack trace-0        -> url\tfunc\nvald version         -> \x1b[1m1.0\x1b[22m",
			},
		},
		{
			name: "string return valid string with no stacktrace initialized",
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
				PrepOnce:          sync.Once{},
			},
			beforeFunc: func() {
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
			},
			want: want{
				want: "\nbuild cpu info flags -> []\nbuild time           -> bt\ncgo enabled          -> true\ngit commit           -> commit\ngo arch              -> goarch\ngo os                -> goos\ngo version           -> 1.1\nngt version          -> 1.2\nserver name          -> srv\nvald version         -> \x1b[1m1.0\x1b[22m",
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
				PrepOnce:          test.fields.PrepOnce,
			}

			got := d.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDetail_Get(t *testing.T) {
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
		PrepOnce          sync.Once
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
			return errors.Errorf("\tgot: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "get return detail object",
			beforeFunc: func() {
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					return uintptr(0), "", 0, false
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
			},
			want: want{
				want: Detail{
					GitCommit:  "master",
					GoVersion:  runtime.Version(),
					GoOS:       runtime.GOOS,
					GoArch:     runtime.GOARCH,
					StackTrace: []StackTrace{},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with stacktrace",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), "info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/vdaas/vald/tree/master",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     "info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with fila name has goroot prefix",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), runtime.GOROOT() + "/src/info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/golang/go/blob/go1.15.5/src/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     runtime.GOROOT() + "/src/info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with go mod path",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/vdaas/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas/vald/internal/info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with go mod path with version",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/vdaas/blob/v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932/vald/internal/info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with go mod path with 2 version",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/vdaas/blob/master/vald/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     "/tmp/go/pkg/mod/github.com/vdaas@v0.0.0-20171023180738-a3a6125de932-a843423387/vald/internal/info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
			},
		},
		{
			name: "get return detail object with go src path",
			beforeFunc: func() {
				i := 0
				rtCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
					if i == 0 {
						i++
						return uintptr(0), "/tmp/go/src/github.com/vdaas/vald/internal/info.go", 100, true
					}
					return uintptr(1), "info.go", 100, false
				}
				rtFuncForPC = func(uintptr) *runtime.Func {
					return runtime.FuncForPC(reflect.ValueOf(TestDetail_Get).Pointer())
				}
			},
			afterFunc: func() {
				rtCaller = runtime.Caller
				rtFuncForPC = runtime.FuncForPC
			},
			want: want{
				want: Detail{
					GitCommit: "master",
					GoVersion: runtime.Version(),
					GoOS:      runtime.GOOS,
					GoArch:    runtime.GOARCH,
					StackTrace: []StackTrace{
						StackTrace{
							URL:      "https://github.com/vdaas/vald/blob/master/internal/info.go#L100",
							FuncName: "github.com/vdaas/vald/internal/info.TestDetail_Get",
							File:     "/tmp/go/src/github.com/vdaas/vald/internal/info.go",
							Line:     100,
						},
					},
					PrepOnce: func() sync.Once {
						o := sync.Once{}
						o.Do(func() {})
						return o
					}(),
				},
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
				PrepOnce:          test.fields.PrepOnce,
			}

			got := d.Get()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDetail_prepare(t *testing.T) {
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
		checkFunc  func(want, *Detail) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Detail) error {
		opts := []cmp.Option{
			cmpopts.IgnoreFields(*w.want, "PrepOnce"),
		}
		if diff := cmp.Diff(*w.want, *got, opts...); len(diff) != 0 {
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
				want: &Detail{
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

		{
			name: "GitCommit and Version field is not overwritten when GitCommit field is `internal`",
			fields: fields{
				GitCommit: "internal",
			},
			want: want{
				want: &Detail{
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

		{
			name: "BuildTime field is not overwritten when BuildTime field is `10`",
			fields: fields{
				BuildTime: "10s",
			},
			want: want{
				want: &Detail{
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

		{
			name: "GoVersion field is not overwritten when GoVersion field is `1.14`",
			fields: fields{
				GoVersion: "1.14",
			},
			want: want{
				want: &Detail{
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

		{
			name: "GoOS field is not overwritten when GoOS field is `linux`",
			fields: fields{
				GoOS: "linux",
			},
			want: want{
				want: &Detail{
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

		{
			name: "GoArch fields is not overwritten when GoArch field is `amd`",
			fields: fields{
				GoArch: "amd",
			},
			want: want{
				want: &Detail{
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

		{
			name: "CGOEnabled field is not overwritten when CGOEnabled field is `1`",
			fields: fields{
				CGOEnabled: "1",
			},
			want: want{
				want: &Detail{
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

		{
			name: "NGTVersion field is not overwritten when NGTVersion field is `v1.11.5`",
			fields: fields{
				NGTVersion: "v1.11.5",
			},
			want: want{
				want: &Detail{
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

		{
			name: "BuildCPUInfoFlags is not overwritten when BuildCPUInfoFlags field is `test`",
			fields: fields{
				BuildCPUInfoFlags: []string{"test"},
			},
			want: want{
				want: &Detail{
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
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc()

			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			d := &Detail{
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

			d.prepare()
			if err := test.checkFunc(test.want, d); err != nil {
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
		want *Detail
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Detail) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *Detail) error {
		opts := []cmp.Option{
			cmpopts.IgnoreFields(*w.want, "PrepOnce"),
		}
		if diff := cmp.Diff(*w.want, *got, opts...); len(diff) != 0 {
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
				want: &Detail{
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
			if err := test.checkFunc(test.want, &detail); err != nil {
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
