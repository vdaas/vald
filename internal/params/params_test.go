//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package params provides implementation of Go API for argument parser
package params

import (
	stderrs "errors"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *parser
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *parser) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *parser) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns *parser when opts is nil",
			want: want{
				want: &parser{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						keys: []string{
							"f",
							"file",
							"c",
							"config",
						},
						defaultPath: "/etc/server/config.yaml",
						description: "config file path",
					},
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						keys: []string{
							"v",
							"ver",
							"version",
						},
						defaultFlag: false,
						description: "show server version",
					},
				},
			},
		},

		{
			name: "returns *parser when opts is not nil",
			args: args{
				opts: []Option{
					WithConfigFilePathKeys("t", "test"),
				},
			},
			want: want{
				want: &parser{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						keys: []string{
							"f",
							"file",
							"c",
							"config",
							"t",
							"test",
						},
						defaultPath: "/etc/server/config.yaml",
						description: "config file path",
					},
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						keys: []string{
							"v",
							"ver",
							"version",
						},
						defaultFlag: false,
						description: "show server version",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got := New(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	type fields struct {
		filePath struct {
			keys        []string
			defaultPath string
			description string
		}
		version struct {
			keys        []string
			defaultFlag bool
			description string
		}
	}
	type want struct {
		want  Data
		want1 bool
		err   error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, Data, bool, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Data, got1 bool, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (d, false, nil) when parse succeed",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys: []string{
						"path", "p",
					},
					defaultPath: "./params.go",
					description: "sets file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys: []string{
						"version", "v",
					},
					defaultFlag: true,
					description: "show version",
				},
			},
			beforeFunc: func() {
				os.Args = []string{
					"test", "--path=./params.go", "--version=false",
				}
			},
			afterFunc: func() { os.Args = nil },
			want: want{
				want: &data{
					configFilePath: "./params.go",
					showVersion:    false,
				},
			},
		},

		{
			name: "returns (nil, true, nil) When parse fails but the help option is set",
			beforeFunc: func() {
				os.Args = []string{
					"test", "--help",
				}
			},
			afterFunc: func() { os.Args = nil },
			want: want{
				want1: true,
			},
		},

		{
			name: "returns (nil, true, nil) When parse fails but the help option is not set",
			beforeFunc: func() {
				os.Args = []string{
					"test", "--name",
				}
			},
			afterFunc: func() { os.Args = nil },
			want: want{
				want1: false,
				err:   errors.ErrArgumentParseFailed(stderrs.New("flag provided but not defined: -name")),
			},
		},

		{
			name: "returns (nil, true, error) When the configFilePath option is set but file dose not exist",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys: []string{
						"path", "p",
					},
					description: "sets file path",
				},
			},
			beforeFunc: func() {
				os.Args = []string{
					"test", "--path=config.yml",
				}
			},
			afterFunc: func() { os.Args = nil },
			want: want{
				want1: true,
				err: &os.PathError{
					Op:   "stat",
					Path: "config.yml",
					Err:  syscall.Errno(0x2),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			p := &parser{
				filePath: test.fields.filePath,
				version:  test.fields.version,
			}

			got, got1, err := p.Parse()
			if err := checkFunc(test.want, got, got1, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_data_ConfigFilePath(t *testing.T) {
	type fields struct {
		configFilePath string
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
			name: "returns `./path` when d.configFilePath is `./path`",
			fields: fields{
				configFilePath: "./path",
			},
			want: want{
				want: "./path",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			d := &data{
				configFilePath: test.fields.configFilePath,
			}

			got := d.ConfigFilePath()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_data_ShowVersion(t *testing.T) {
	type fields struct {
		showVersion bool
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns true when d.showVersion is true",
			fields: fields{
				showVersion: true,
			},
			want: want{
				want: true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			d := &data{
				showVersion: test.fields.showVersion,
			}

			got := d.ShowVersion()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
