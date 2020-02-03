//
// Copyright (C) 2020 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestNew(t *testing.T) {
	type test struct {
		name string
		opts []Option
		want *parser
	}

	tests := []test{
		{
			name: "returns parser instance when opts is empty",
			opts: nil,
			want: &parser{
				filePath: filePath{
					keys: []string{
						"f",
						"file",
						"c",
						"config",
					},
					defaultPath: "/etc/server/config.yaml",
					description: "config file path",
				},
				version: version{
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

		{
			name: "returns parser instance when opts is not empty",
			opts: []Option{
				WithConfigFilePathKeys("p", "path"),
			},
			want: &parser{
				filePath: filePath{
					keys: []string{
						"f",
						"file",
						"c",
						"config",
						"p",
						"path",
					},
					defaultPath: "/etc/server/config.yaml",
					description: "config file path",
				},
				version: version{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	type fields struct {
		filePath filePath
		version  version
	}

	type global struct {
		args []string
	}

	type test struct {
		name      string
		fields    fields
		global    global
		checkFunc func(got, want error) error
		wantData  *Data
		wantFail  bool
		wantError error
	}

	tests := []test{
		{
			name: "returns data and flag and nil when parser is successes",
			fields: fields{
				filePath: filePath{
					keys: []string{
						"conf",
					},
					defaultPath: "vald_config.yml",
					description: "set config file",
				},
				version: version{
					keys: []string{
						"version",
					},
					defaultFlag: false,
					description: "print version",
				},
			},
			global: global{
				args: []string{
					"vald",
					"--conf",
					"./config/vald_config.yml",
					"--version",
					"true",
				},
			},
			checkFunc: func(got, want error) error {
				if got != nil {
					return errors.New("err is not nil")
				}
				return nil
			},
			wantData: &Data{
				configFilePath: "./config/vald_config.yml",
				showVersion:    true,
			},
			wantFail:  false,
			wantError: nil,
		},

		{
			name: "returns data and flag and error when parse error occurs and err is not equal flag.ErrHelp",
			fields: fields{
				filePath: filePath{
					keys: []string{
						"conf",
					},
					defaultPath: "vald_config.yml",
					description: "set config file",
				},
				version: version{
					keys: []string{
						"version",
					},
					defaultFlag: false,
					description: "print version",
				},
			},
			global: global{
				args: []string{
					"vald",
					"--conf",
					"./config/vald_config.yml",
					"-name",
					"true",
				},
			},
			checkFunc: func(got, want error) error {
				if got == nil {
					return errors.New("err is nil")
				} else if !errors.Is(want, got) {
					return errors.Errorf("err is not equal. want: %v, got: %v", want, got)
				}
				return nil
			},
			wantData:  nil,
			wantFail:  false,
			wantError: errors.ErrArgumentParseFailed(fmt.Errorf("flag provided but not defined: -name")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.global.args

			d, fail, err := (&parser{
				filePath: tt.fields.filePath,
				version:  tt.fields.version,
			}).Parse()

			if want, got := tt.wantData, d; !reflect.DeepEqual(want, got) {
				t.Errorf("data is not equal. want: %v, got: %v", want, got)
			}

			if want, got := tt.wantFail, fail; want != got {
				t.Errorf("fail is not equal. want: %v, got: %v", want, got)
			}

			if err := tt.checkFunc(err, tt.wantError); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestData_ConfigFilePath(t *testing.T) {
	type test struct {
		name           string
		configFilePath string
		want           string
	}

	tests := []test{
		{
			name:           "returns config file path",
			configFilePath: "config_file_path",
			want:           "config_file_path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				configFilePath: tt.configFilePath,
			}

			if got, want := d.ConfigFilePath(), tt.want; got != want {
				t.Errorf("not equals. want: %v, got: %v", want, got)
			}
		})
	}
}

func TestData_ShowVersion(t *testing.T) {
	type test struct {
		name        string
		showVersion bool
		want        bool
	}

	tests := []test{
		{
			name:        "returns show version flag",
			showVersion: true,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				showVersion: tt.showVersion,
			}

			if got, want := d.ShowVersion(), tt.want; got != want {
				t.Errorf("not equals. want: %v, got: %v", want, got)
			}
		})
	}
}
