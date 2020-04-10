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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
)

func Test_newGRPCClientConfig(t *testing.T) {
	tests := []struct {
		name string
		want *GRPCClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGRPCClientConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGRPCClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCClient_Bind(t *testing.T) {
	tests := []struct {
		name string
		g    *GRPCClient
		want *GRPCClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCClient.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCClientKeepalive_Bind(t *testing.T) {
	tests := []struct {
		name string
		g    *GRPCClientKeepalive
		want *GRPCClientKeepalive
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCClientKeepalive.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallOption_Bind(t *testing.T) {
	tests := []struct {
		name string
		c    *CallOption
		want *CallOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallOption.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialOption_Bind(t *testing.T) {
	tests := []struct {
		name string
		d    *DialOption
		want *DialOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DialOption.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCClient_Opts(t *testing.T) {
	tests := []struct {
		name string
		g    *GRPCClient
		want []grpc.Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Opts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCClient.Opts() = %v, want %v", got, tt.want)
			}
		})
	}
}
