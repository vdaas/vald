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

// Package discoverer
package discoverer

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantD   Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("New() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func Test_client_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetAddrs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		c         *client
		args      args
		wantAddrs []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAddrs := tt.c.GetAddrs(tt.args.ctx); !reflect.DeepEqual(gotAddrs, tt.wantAddrs) {
				t.Errorf("client.GetAddrs() = %v, want %v", gotAddrs, tt.wantAddrs)
			}
		})
	}
}

func Test_client_GetClient(t *testing.T) {
	tests := []struct {
		name string
		c    *client
		want grpc.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_connect(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.connect(tt.args.ctx, tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("client.connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_dnsDiscovery(t *testing.T) {
	type args struct {
		ctx context.Context
		ech chan<- error
	}
	tests := []struct {
		name      string
		c         *client
		args      args
		wantAddrs []string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddrs, err := tt.c.dnsDiscovery(tt.args.ctx, tt.args.ech)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.dnsDiscovery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAddrs, tt.wantAddrs) {
				t.Errorf("client.dnsDiscovery() = %v, want %v", gotAddrs, tt.wantAddrs)
			}
		})
	}
}

func Test_client_discover(t *testing.T) {
	type args struct {
		ctx context.Context
		ech chan<- error
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.discover(tt.args.ctx, tt.args.ech); (err != nil) != tt.wantErr {
				t.Errorf("client.discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
