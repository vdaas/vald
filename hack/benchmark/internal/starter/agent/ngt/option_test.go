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

// Package ngt provides ngt agent starter  functionality
package ngt

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/pkg/agent/ngt/config"
)

func TestWithConfig(t *testing.T) {
	type args struct {
		cfg *config.Data
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
			if got := WithConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDimentaion(t *testing.T) {
	type args struct {
		d int
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
			if got := WithDimentaion(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDimentaion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDistanceType(t *testing.T) {
	type args struct {
		dtype string
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
			if got := WithDistanceType(tt.args.dtype); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDistanceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithObjectType(t *testing.T) {
	type args struct {
		otype string
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
			if got := WithObjectType(tt.args.otype); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithObjectType() = %v, want %v", got, tt.want)
			}
		})
	}
}
