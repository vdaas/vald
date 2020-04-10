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

// Package service
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/internal/net/grpc"
)

func TestNewGateway(t *testing.T) {
	type args struct {
		opts []GWOption
	}
	tests := []struct {
		name    string
		args    args
		wantGw  Gateway
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGw, err := NewGateway(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGateway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGw, tt.wantGw) {
				t.Errorf("NewGateway() = %v, want %v", gotGw, tt.wantGw)
			}
		})
	}
}

func Test_gateway_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		g       *gateway
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("gateway.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gateway.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gateway_BroadCast(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error
	}
	tests := []struct {
		name    string
		g       *gateway
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.BroadCast(tt.args.ctx, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gateway.BroadCast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gateway_Do(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error
	}
	tests := []struct {
		name    string
		g       *gateway
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Do(tt.args.ctx, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gateway.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gateway_DoMulti(t *testing.T) {
	type args struct {
		ctx context.Context
		num int
		f   func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error
	}
	tests := []struct {
		name    string
		g       *gateway
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.DoMulti(tt.args.ctx, tt.args.num, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gateway.DoMulti() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gateway_GetAgentCount(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		g    *gateway
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.GetAgentCount(tt.args.ctx); got != tt.want {
				t.Errorf("gateway.GetAgentCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
