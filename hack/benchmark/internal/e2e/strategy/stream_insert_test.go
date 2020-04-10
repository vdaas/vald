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

func TestNewStreamInsert(t *testing.T) {
	type args struct {
		opts []StreamInsertOption
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
			if got := NewStreamInsert(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStreamInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_streamInsert_dataProvider(t *testing.T) {
	type args struct {
		total   *uint32
		b       *testing.B
		dataset assets.Dataset
	}
	tests := []struct {
		name  string
		sisrt *streamInsert
		args  args
		want  func() *client.ObjectVector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sisrt.dataProvider(tt.args.total, tt.args.b, tt.args.dataset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("streamInsert.dataProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_streamInsert_Run(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		c       client.Client
		dataset assets.Dataset
	}
	tests := []struct {
		name  string
		sisrt *streamInsert
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sisrt.Run(tt.args.ctx, tt.args.b, tt.args.c, tt.args.dataset)
		})
	}
}
