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
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want starter.Starter
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

func Test_server_Run(t *testing.T) {
	type args struct {
		ctx context.Context
		tb  testing.TB
	}
	tests := []struct {
		name string
		ns   *server
		args args
		want func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ns.Run(tt.args.ctx, tt.args.tb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_createIndexDir(t *testing.T) {
	tests := []struct {
		name    string
		ns      *server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ns.createIndexDir(); (err != nil) != tt.wantErr {
				t.Errorf("server.createIndexDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_clearIndexDir(t *testing.T) {
	tests := []struct {
		name    string
		ns      *server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ns.clearIndexDir(); (err != nil) != tt.wantErr {
				t.Errorf("server.clearIndexDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
