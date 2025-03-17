//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package params provides implementation of Go API for argument parser
package params

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// TestNew tests the New function for creating a new parser instance.
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
		checkFunc  func(want, Parser) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	// Custom check function: compare only the essential fields.
	defaultCheckFunc := func(w want, got Parser) error {
		p, ok := got.(*parser)
		if !ok {
			return errors.Errorf("got is not *parser")
		}
		// Check filePath fields.
		if !reflect.DeepEqual(p.filePath.keys, w.want.filePath.keys) {
			return errors.Errorf("filePath.keys mismatch: got %v, want %v", p.filePath.keys, w.want.filePath.keys)
		}
		if p.filePath.defaultPath != w.want.filePath.defaultPath {
			return errors.Errorf("filePath.defaultPath mismatch: got %v, want %v", p.filePath.defaultPath, w.want.filePath.defaultPath)
		}
		if p.filePath.description != w.want.filePath.description {
			return errors.Errorf("filePath.description mismatch: got %v, want %v", p.filePath.description, w.want.filePath.description)
		}
		// Check version fields.
		if !reflect.DeepEqual(p.version.keys, w.want.version.keys) {
			return errors.Errorf("version.keys mismatch: got %v, want %v", p.version.keys, w.want.version.keys)
		}
		if p.version.defaultFlag != w.want.version.defaultFlag {
			return errors.Errorf("version.defaultFlag mismatch: got %v, want %v", p.version.defaultFlag, w.want.version.defaultFlag)
		}
		if p.version.description != w.want.version.description {
			return errors.Errorf("version.description mismatch: got %v, want %v", p.version.description, w.want.version.description)
		}
		return nil
	}
	tests := []test{
		{
			name: "should return a default parser when no options are provided",
			want: want{
				want: &parser{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						keys:        []string{"f", "file", "c", "config"},
						defaultPath: "/etc/server/config.yaml",
						description: "config file path",
					},
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						keys:        []string{"v", "ver", "version"},
						defaultFlag: false,
						description: "show server version",
					},
				},
			},
		},
		{
			name: "should return a parser with additional config file keys when options are provided",
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
						keys:        []string{"f", "file", "c", "config", "t", "test"},
						defaultPath: "/etc/server/config.yaml",
						description: "config file path",
					},
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						keys:        []string{"v", "ver", "version"},
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
			check := test.checkFunc
			if check == nil {
				check = defaultCheckFunc
			}
			got := New(test.args.opts...)
			if err := check(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// Test_parser_Parse tests the Parse method of the parser.
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
		want Data  // expected Data (may be nil)
		help bool  // indicates if help option was triggered
		err  error // expected error
	}
	type test struct {
		name       string
		fields     fields
		args       []string // custom os.Args (optional)
		want       want
		checkFunc  func(want, Data, bool, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	// Custom check function: compare only the essential fields of Data.
	defaultCheckFunc := func(w want, got Data, gotHelp bool, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error: %#v, want: %#v", err, w.err)
		}
		if gotHelp != w.help {
			return errors.Errorf("got help flag: %#v, want: %#v", gotHelp, w.help)
		}
		// If no expected data is provided, skip field comparison.
		if w.want == nil {
			return nil
		}
		d, ok := got.(*data)
		if !ok {
			return errors.Errorf("got is not *data")
		}
		expected, ok := w.want.(*data)
		if !ok {
			return errors.Errorf("expected want is not *data")
		}
		// If expected configFilePath is non-empty, compare directly.
		// Otherwise, ensure that got.ConfigFilePath() is not empty.
		if expected.configFilePath != "" {
			if d.configFilePath != expected.configFilePath {
				return errors.Errorf("configFilePath mismatch: got %v, want %v", d.configFilePath, expected.configFilePath)
			}
		} else {
			if d.configFilePath == "" {
				return errors.Errorf("expected non-empty configFilePath, but got empty string")
			}
		}
		if d.showVersion != expected.showVersion {
			return errors.Errorf("showVersion mismatch: got %v, want %v", d.showVersion, expected.showVersion)
		}
		return nil
	}
	tests := []test{
		{
			name: "should successfully parse valid config file and version flag false",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys:        []string{"path", "p"},
					defaultPath: "",
					description: "set file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys:        []string{"version", "v"},
					defaultFlag: false,
					description: "show version flag",
				},
			},
			beforeFunc: func(t *testing.T) {
				// Create a temporary file to ensure config file existence.
				tmpFile, err := os.CreateTemp("", "config-*.yaml")
				if err != nil {
					t.Fatal(err)
				}
				// Ensure the temporary file is removed after the test.
				t.Cleanup(func() { os.Remove(tmpFile.Name()) })
				tmpFile.Close()
				os.Args = []string{"test", "--path=" + tmpFile.Name(), "--version=false"}
			},
			afterFunc: func(t *testing.T) {
				os.Args = nil
			},
			want: want{
				// expected Data: only showVersion is checked directly.
				// For configFilePath, we expect a non-empty string.
				want: &data{
					configFilePath: "", // will be validated as non-empty
					showVersion:    false,
				},
				help: false,
				err:  nil,
			},
			checkFunc: func(w want, got Data, gotHelp bool, err error) error {
				// Use the default check function for essential fields.
				return defaultCheckFunc(w, got, gotHelp, err)
			},
		},
		{
			name: "should parse and return valid data when version flag is true even if file does not exist",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys:        []string{"path", "p"},
					defaultPath: "nonexistent.yaml",
					description: "set file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys:        []string{"version", "v"},
					defaultFlag: false,
					description: "show version flag",
				},
			},
			beforeFunc: func(t *testing.T) {
				os.Args = []string{"test", "--path=nonexistent.yaml", "--version=true"}
			},
			afterFunc: func(t *testing.T) {
				os.Args = nil
			},
			want: want{
				want: &data{
					configFilePath: "nonexistent.yaml",
					showVersion:    true,
				},
				help: false,
				err:  nil,
			},
		},
		{
			name: "should return help when --help flag is provided",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys:        []string{"path", "p"},
					defaultPath: "/etc/server/config.yaml",
					description: "set file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys:        []string{"version", "v"},
					defaultFlag: false,
					description: "show version flag",
				},
			},
			beforeFunc: func(t *testing.T) {
				os.Args = []string{"test", "--help"}
			},
			afterFunc: func(t *testing.T) {
				os.Args = nil
			},
			want: want{
				help: true,
				err:  nil,
			},
		},
		{
			name: "should return parsing error for unknown flag",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys:        []string{"path", "p"},
					defaultPath: "",
					description: "set file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys:        []string{"version", "v"},
					defaultFlag: false,
					description: "show version flag",
				},
			},
			beforeFunc: func(t *testing.T) {
				os.Args = []string{"test", "--unknown"}
			},
			afterFunc: func(t *testing.T) {
				os.Args = nil
			},
			want: want{
				help: false,
				// The error message is wrapped by ErrArgumentParseFailed.
				err: errors.ErrArgumentParseFailed(errors.New("flag provided but not defined: -unknown")),
			},
		},
		{
			name: "should return help when config file path is empty",
			fields: fields{
				filePath: struct {
					keys        []string
					defaultPath string
					description string
				}{
					keys:        []string{"path", "p"},
					defaultPath: "",
					description: "set file path",
				},
				version: struct {
					keys        []string
					defaultFlag bool
					description string
				}{
					keys:        []string{"version", "v"},
					defaultFlag: false,
					description: "show version flag",
				},
			},
			beforeFunc: func(t *testing.T) {
				os.Args = []string{"test", "--path=", "--version=false"}
			},
			afterFunc: func(t *testing.T) {
				os.Args = nil
			},
			want: want{
				help: true,
				err:  errors.New("invalid argument"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			p := &parser{
				filePath: test.fields.filePath,
				version:  test.fields.version,
				f:        flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
			}
			gotData, gotHelp, err := p.Parse()
			if test.checkFunc != nil {
				if err := test.checkFunc(test.want, gotData, gotHelp, err); err != nil {
					tt.Errorf("error = %v", err)
				}
			} else if err := defaultCheckFunc(test.want, gotData, gotHelp, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// Test_data_ConfigFilePath tests the ConfigFilePath getter of the data struct.
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got string) error {
		if got != w.want {
			return errors.Errorf("got: %v, want: %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "should return the provided config file path",
			fields: fields{
				configFilePath: "./path/to/config.yaml",
			},
			want: want{
				want: "./path/to/config.yaml",
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			d := &data{
				configFilePath: test.fields.configFilePath,
			}
			got := d.ConfigFilePath()
			if test.checkFunc != nil {
				if err := test.checkFunc(test.want, got); err != nil {
					tt.Errorf("error = %v", err)
				}
			} else if err := defaultCheckFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// Test_data_ShowVersion tests the ShowVersion getter of the data struct.
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if got != w.want {
			return errors.Errorf("got: %v, want: %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "should return true when showVersion is set to true",
			fields: fields{
				showVersion: true,
			},
			want: want{
				want: true,
			},
		},
		{
			name: "should return false when showVersion is set to false",
			fields: fields{
				showVersion: false,
			},
			want: want{
				want: false,
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			d := &data{
				showVersion: test.fields.showVersion,
			}
			got := d.ShowVersion()
			if test.checkFunc != nil {
				if err := test.checkFunc(test.want, got); err != nil {
					tt.Errorf("error = %v", err)
				}
			} else if err := defaultCheckFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
