// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Test_gRPCClient_executeRPC_noResourceExhaustedRetry verifies that only
// ResourceExhausted becomes non-retryable under WithoutResourceExhaustedRetry,
// and that every code keeps its default classification without it.
func Test_gRPCClient_executeRPC_noResourceExhaustedRetry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		code          codes.Code
		marked        bool
		wantRetryable bool
	}{
		{name: "ResourceExhausted is retried by default", code: codes.ResourceExhausted, wantRetryable: true},
		{name: "ResourceExhausted is not retried when marked", code: codes.ResourceExhausted, marked: true},
		{name: "Unavailable is still retried when marked", code: codes.Unavailable, marked: true, wantRetryable: true},
		{name: "Internal is still retried when marked", code: codes.Internal, marked: true, wantRetryable: true},
		{name: "InvalidArgument is never retried", code: codes.InvalidArgument, marked: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g, ok := New("test").(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			conn, err := ggrpc.NewClient("passthrough:///re-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			rpcErr := status.Error(tc.code, "rpc failed")
			ctx := WithQuietErrorLog(t.Context())
			if tc.marked {
				ctx = WithoutResourceExhaustedRetry(ctx)
			}
			_, retryable, err := g.executeRPC(ctx, healthyConnPool{fixedConnPool{conn: conn}}, "10.0.0.1:8081",
				func(context.Context, *ClientConn, ...CallOption) (any, error) { return nil, rpcErr })
			if !errors.Is(err, rpcErr) {
				t.Fatalf("executeRPC() err = %v, want %v", err, rpcErr)
			}
			if retryable != tc.wantRetryable {
				t.Errorf("executeRPC() retryable = %v, want %v", retryable, tc.wantRetryable)
			}
		})
	}
}
