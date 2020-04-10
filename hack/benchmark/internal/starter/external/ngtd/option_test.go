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

// Package ngtd provides ngtd starter  functionality
package ngtd

import (
	"reflect"
	"testing"
)

func TestWithDimentaion(t *testing.T) {
	type args struct {
		dim int
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
			if got := WithDimentaion(tt.args.dim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDimentaion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServerType(t *testing.T) {
	type args struct {
		t ServerType
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
			if got := WithServerType(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServerType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIndexDir(t *testing.T) {
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
			if got := WithIndexDir(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIndexDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port int
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
			if got := WithPort(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
