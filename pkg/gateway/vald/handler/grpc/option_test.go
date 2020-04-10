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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
)

func TestWithGateway(t *testing.T) {
	type args struct {
		g service.Gateway
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithGateway(tt.args.g); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMeta(t *testing.T) {
	type args struct {
		m service.Meta
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMeta(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithBackup(t *testing.T) {
	type args struct {
		b service.Backup
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithBackup(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithBackup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFilters(t *testing.T) {
	type args struct {
		filter service.Filter
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithFilters(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithErrGroup(tt.args.eg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithErrGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReplicationCount(t *testing.T) {
	type args struct {
		rep int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithReplicationCount(tt.args.rep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReplicationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithStreamConcurrency(t *testing.T) {
	type args struct {
		c int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithStreamConcurrency(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithStreamConcurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}
