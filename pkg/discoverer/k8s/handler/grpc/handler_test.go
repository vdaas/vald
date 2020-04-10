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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantDs  DiscovererServer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDs, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDs, tt.wantDs) {
				t.Errorf("New() = %v, want %v", gotDs, tt.wantDs)
			}
		})
	}
}

func Test_server_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		s    *server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Start(tt.args.ctx)
		})
	}
}

func Test_server_Pods(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Discoverer_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Info_Pods
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Pods(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Pods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Pods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Nodes(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Discoverer_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Info_Nodes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Nodes(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Nodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Nodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_singleflightKey(t *testing.T) {
	type args struct {
		pref string
		req  *payload.Discoverer_Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := singleflightKey(tt.args.pref, tt.args.req); got != tt.want {
				t.Errorf("singleflightKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
