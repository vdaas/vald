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

// Package e2e provides e2e testing framework functions
package e2e

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client"
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

func TestWithClient(t *testing.T) {
	type args struct {
		c client.Client
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
			if got := WithClient(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithStrategy(t *testing.T) {
	type args struct {
		strategis []Strategy
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
			if got := WithStrategy(tt.args.strategis...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServerStarter(t *testing.T) {
	type args struct {
		f func(context.Context, testing.TB, assets.Dataset) func()
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
			if got := WithServerStarter(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServerStarter() = %v, want %v", got, tt.want)
			}
		})
	}
}
