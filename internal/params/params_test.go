//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		p       *parser
		want    Data
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.p.Parse()
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

func Test_data_ConfigFilePath(t *testing.T) {
	tests := []struct {
		name string
		d    *data
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ConfigFilePath(); got != tt.want {
				t.Errorf("data.ConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_ShowVersion(t *testing.T) {
	tests := []struct {
		name string
		d    *data
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ShowVersion(); got != tt.want {
				t.Errorf("data.ShowVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
