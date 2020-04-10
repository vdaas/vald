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

	"github.com/vdaas/vald/internal/servers/server"
)

func TestServers_Bind(t *testing.T) {
	tests := []struct {
		name string
		s    *Servers
		want *Servers
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Servers.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServers_GetGRPCStreamConcurrency(t *testing.T) {
	tests := []struct {
		name  string
		s     *Servers
		wantC int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotC := tt.s.GetGRPCStreamConcurrency(); gotC != tt.wantC {
				t.Errorf("Servers.GetGRPCStreamConcurrency() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestHTTP_Bind(t *testing.T) {
	tests := []struct {
		name string
		h    *HTTP
		want *HTTP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTP.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPC_Bind(t *testing.T) {
	tests := []struct {
		name string
		g    *GRPC
		want *GRPC
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPC.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCKeepalive_Bind(t *testing.T) {
	tests := []struct {
		name string
		k    *GRPCKeepalive
		want *GRPCKeepalive
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCKeepalive.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Bind(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Opts(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want []server.Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Opts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Opts() = %v, want %v", got, tt.want)
			}
		})
	}
}
