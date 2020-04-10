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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/net/grpc"
)

func TestWithMetaAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaHostPort(t *testing.T) {
	type args struct {
		host string
		port int
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaHostPort(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaHostPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaClient(t *testing.T) {
	type args struct {
		client grpc.Client
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaClient(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaCacheEnabled(t *testing.T) {
	type args struct {
		flg bool
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaCacheEnabled(tt.args.flg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaCacheEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaCache(t *testing.T) {
	type args struct {
		c cache.Cache
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaCache(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaCacheExpireDuration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaCacheExpireDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaCacheExpireDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetaCacheExpiredCheckDuration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want MetaOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetaCacheExpiredCheckDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetaCacheExpiredCheckDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
