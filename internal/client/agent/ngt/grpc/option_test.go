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

// Package grpc provides agent ngt gRPC client functions
package grpc

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
)

func TestWithAddr(t *testing.T) {
	type args struct {
		addr string
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
			if got := WithAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithGRPCClientOption(t *testing.T) {
	type args struct {
		opts []grpc.Option
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
			if got := WithGRPCClientOption(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithGRPCClientOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
