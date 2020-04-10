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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	"reflect"
	"testing"
)

func TestWithSessionOptions(t *testing.T) {
	type args struct {
		opts *SessionOptions
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
			if got := WithSessionOptions(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSessionOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSessionTarget(t *testing.T) {
	type args struct {
		tgt string
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
			if got := WithSessionTarget(tt.args.tgt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSessionTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSessionConfig(t *testing.T) {
	type args struct {
		cfg []byte
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
			if got := WithSessionConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSessionConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOperations(t *testing.T) {
	type args struct {
		opes []*Operation
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
			if got := WithOperations(tt.args.opes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOperations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithExportPath(t *testing.T) {
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
			if got := WithExportPath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExportPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTags(t *testing.T) {
	type args struct {
		tags []string
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
			if got := WithTags(tt.args.tags...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
