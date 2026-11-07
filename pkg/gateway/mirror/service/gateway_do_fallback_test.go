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
package service

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/atomic"
	"google.golang.org/grpc/metadata"
)

// fallbackGRPCClient fails every Do with doErr and counts RoundRobin calls.
type fallbackGRPCClient struct {
	grpc.Client
	doErr   error
	rrCalls atomic.Int64
}

func (c *fallbackGRPCClient) Do(
	context.Context, string, func(context.Context, *grpc.ClientConn, ...grpc.CallOption) (any, error),
) (any, error) {
	return nil, c.doErr
}

func (c *fallbackGRPCClient) RoundRobin(
	context.Context, func(context.Context, *grpc.ClientConn, ...grpc.CallOption) (any, error),
) (any, error) {
	c.rrCalls.Add(1)
	return "peer", nil
}

type fallbackMirrorClient struct {
	mirror.Client
	gc grpc.Client
}

func (c *fallbackMirrorClient) GRPCClient() grpc.Client { return c.gc }

func Test_gateway_Do_noFallbackForProxiedRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		forwardedBy string
		wantRR      int64
		wantErr     bool
	}{
		{
			name:        "proxied request returns the own cluster's error without falling back to other mirrors",
			forwardedBy: "mirror-b",
			wantRR:      0,
			wantErr:     true,
		},
		{
			name:    "user request still falls back to other mirrors",
			wantRR:  1,
			wantErr: false,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			gc := &fallbackGRPCClient{doErr: status.Error(codes.NotFound, "not found")}
			g := &gateway{client: &fallbackMirrorClient{gc: gc}, podName: "mirror-a"}
			ctx := tt.Context()
			if test.forwardedBy != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.MD{forwardedContextKey: []string{test.forwardedBy}})
			}
			_, err := g.Do(ctx, "vald-lb", func(context.Context, string, MirrorClient, ...grpc.CallOption) (any, error) {
				return "own", nil
			})
			if (err != nil) != test.wantErr {
				tt.Errorf("err = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr && !status.Is(err, codes.NotFound) {
				tt.Errorf("err = %v, want the own cluster's NotFound", err)
			}
			if got := gc.rrCalls.Load(); got != test.wantRR {
				tt.Errorf("RoundRobin called %d times, want %d", got, test.wantRR)
			}
		})
	}
}
