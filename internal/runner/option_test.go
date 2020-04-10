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

// Package runner provides implementation of process runner
package runner

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
)

func TestWithName(t *testing.T) {
	type args struct {
		name string
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
			if got := WithName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithVersion(t *testing.T) {
	type args struct {
		ver string
		max string
		min string
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
			if got := WithVersion(tt.args.ver, tt.args.max, tt.args.min); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConfigLoader(t *testing.T) {
	type args struct {
		f func(string) (interface{}, *config.GlobalConfig, error)
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
			if got := WithConfigLoader(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfigLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDaemonInitializer(t *testing.T) {
	type args struct {
		f func(interface{}) (Runner, error)
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
			if got := WithDaemonInitializer(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDaemonInitializer() = %v, want %v", got, tt.want)
			}
		})
	}
}
