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

// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

func TestNewRemove(t *testing.T) {
	type args struct {
		opts []RemoveOption
	}
	tests := []struct {
		name string
		args args
		want e2e.Strategy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRemove(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_remove_Run(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		c       client.Client
		dataset assets.Dataset
	}
	tests := []struct {
		name string
		r    *remove
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Run(tt.args.ctx, tt.args.b, tt.args.c, tt.args.dataset)
		})
	}
}

func Test_remove_run(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		c       client.Client
		dataset assets.Dataset
	}
	tests := []struct {
		name string
		r    *remove
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.run(tt.args.ctx, tt.args.b, tt.args.c, tt.args.dataset)
		})
	}
}

func Test_remove_runParallel(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		c       client.Client
		dataset assets.Dataset
	}
	tests := []struct {
		name string
		r    *remove
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.runParallel(tt.args.ctx, tt.args.b, tt.args.c, tt.args.dataset)
		})
	}
}

func Test_remove_do(t *testing.T) {
	type args struct {
		ctx context.Context
		b   *testing.B
		c   client.Client
		id  string
	}
	tests := []struct {
		name string
		r    *remove
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.do(tt.args.ctx, tt.args.b, tt.args.c, tt.args.id)
		})
	}
}
