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

func TestWithConfigFilePathKeys(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConfigFilePathKeys(tt.args.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfigFilePathKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConfigFilePathDefault(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConfigFilePathDefault(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfigFilePathDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConfigFileDescription(t *testing.T) {
	type args struct {
		desc string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConfigFileDescription(tt.args.desc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfigFileDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithVersionKeys(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithVersionKeys(tt.args.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithVersionKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithVersionFlagDefault(t *testing.T) {
	type args struct {
		flag bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithVersionFlagDefault(tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithVersionFlagDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithVersionDescription(t *testing.T) {
	type args struct {
		desc string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithVersionDescription(tt.args.desc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithVersionDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
