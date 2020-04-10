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
		name string
		args args
		want Server
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

func Test_server_Recover(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Replication_Recovery
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Recover(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Recover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Recover() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Rebalance(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Replication_Rebalance
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Rebalance(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Rebalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Rebalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_AgentInfo(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Empty
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Replication_Agents
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AgentInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.AgentInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.AgentInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
