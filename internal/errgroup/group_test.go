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

// Package errgroup provides server global wait group for graceful kill all goroutine
package errgroup

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  Group
		want1 context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := New(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("New() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantEgctx context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEgctx := Init(tt.args.ctx); !reflect.DeepEqual(gotEgctx, tt.wantEgctx) {
				t.Errorf("Init() = %v, want %v", gotEgctx, tt.wantEgctx)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGo(t *testing.T) {
	type args struct {
		f func() error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Go(tt.args.f)
		})
	}
}

func Test_group_Limitation(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name string
		g    *group
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Limitation(tt.args.limit)
		})
	}
}

func Test_group_Go(t *testing.T) {
	type args struct {
		f func() error
	}
	tests := []struct {
		name string
		g    *group
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Go(tt.args.f)
		})
	}
}

func Test_group_doCancel(t *testing.T) {
	tests := []struct {
		name string
		g    *group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.doCancel()
		})
	}
}

func TestWait(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wait(); (err != nil) != tt.wantErr {
				t.Errorf("Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_group_Wait(t *testing.T) {
	tests := []struct {
		name    string
		g       *group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Wait(); (err != nil) != tt.wantErr {
				t.Errorf("group.Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
