//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *parser
	}{
		{
			name: "return parser success",
			args: args{
				opts: []Option{
					WithVersionKey("dummyVersionKey"),
				},
			},
			want: &parser{
				filePath: filePath{
					key:         "f",
					defaultPath: "/etc/server/config.yaml",
					description: "config file path",
				},
				version: version{
					key:         "dummyVersionKey",
					defaultFlag: false,
					description: "show server version",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	type fields struct {
		filePath struct {
			key         string
			defaultPath string
			description string
		}
		version struct {
			key         string
			defaultFlag bool
			description string
		}
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Data
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				filePath: tt.fields.filePath,
				version:  tt.fields.version,
			}
			got, got1, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.Parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parser.Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestData_ConfigFilePath(t *testing.T) {
	type fields struct {
		configFilePath string
		showVersion    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "return configFilePath success",
			fields: fields{
				configFilePath: "dummy_path",
			},
			want: "dummy_path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				configFilePath: tt.fields.configFilePath,
				showVersion:    tt.fields.showVersion,
			}
			if got := d.ConfigFilePath(); got != tt.want {
				t.Errorf("Data.ConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_ShowVersion(t *testing.T) {
	type fields struct {
		configFilePath string
		showVersion    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "return showVersion success",
			fields: fields{
				showVersion: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				configFilePath: tt.fields.configFilePath,
				showVersion:    tt.fields.showVersion,
			}
			if got := d.ShowVersion(); got != tt.want {
				t.Errorf("Data.ShowVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
