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

package vald

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// attemptCtxClient mimics gRPCClient.RoundRobin under backoff.Do: f gets a per-attempt ctx that is cancelled
// as soon as RoundRobin returns.
type attemptCtxClient struct {
	grpc.Client
	conn *grpc.ClientConn
}

func (c *attemptCtxClient) RoundRobin(
	ctx context.Context,
	f func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error),
) (any, error) {
	actx, cancel := context.WithCancel(ctx)
	defer cancel()
	return f(actx, c.conn)
}

// nopClientStream lets the generated stream constructors finish without a network round trip.
type nopClientStream struct {
	ggrpc.ClientStream
}

func (nopClientStream) SendMsg(any) error { return nil }
func (nopClientStream) CloseSend() error  { return nil }

// TestClient_Stream_notBoundToAttemptCtx checks that every Stream* method opens its stream with the method ctx,
// not with the per-attempt ctx RoundRobin passes to its closure, which backoff.Do cancels on return.
func TestClient_Stream_notBoundToAttemptCtx(t *testing.T) {
	t.Parallel()
	tests := []struct {
		open func(context.Context, *client) error
		name string
	}{
		{name: "StreamSearch", open: func(ctx context.Context, c *client) error { _, err := c.StreamSearch(ctx); return err }},
		{name: "StreamSearchByID", open: func(ctx context.Context, c *client) error { _, err := c.StreamSearchByID(ctx); return err }},
		{name: "StreamLinearSearch", open: func(ctx context.Context, c *client) error { _, err := c.StreamLinearSearch(ctx); return err }},
		{name: "StreamLinearSearchByID", open: func(ctx context.Context, c *client) error { _, err := c.StreamLinearSearchByID(ctx); return err }},
		{name: "StreamInsert", open: func(ctx context.Context, c *client) error { _, err := c.StreamInsert(ctx); return err }},
		{name: "StreamUpdate", open: func(ctx context.Context, c *client) error { _, err := c.StreamUpdate(ctx); return err }},
		{name: "StreamUpsert", open: func(ctx context.Context, c *client) error { _, err := c.StreamUpsert(ctx); return err }},
		{name: "StreamRemove", open: func(ctx context.Context, c *client) error { _, err := c.StreamRemove(ctx); return err }},
		{name: "StreamGetObject", open: func(ctx context.Context, c *client) error { _, err := c.StreamGetObject(ctx); return err }},
		{name: "StreamListObject", open: func(ctx context.Context, c *client) error {
			_, err := c.StreamListObject(ctx, &payload.Object_List_Request{})
			return err
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			var streamCtx context.Context
			conn, err := ggrpc.NewClient("passthrough:///unused",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()),
				ggrpc.WithStreamInterceptor(func(
					ctx context.Context, _ *ggrpc.StreamDesc, _ *ggrpc.ClientConn, _ string, _ ggrpc.Streamer, _ ...ggrpc.CallOption,
				) (ggrpc.ClientStream, error) {
					streamCtx = ctx
					return nopClientStream{}, nil
				}),
			)
			if err != nil {
				tt.Fatalf("failed to create client conn: %v", err)
			}
			defer func() {
				if err := conn.Close(); err != nil {
					tt.Errorf("failed to close client conn: %v", err)
				}
			}()
			c := &client{c: &attemptCtxClient{conn: conn}}
			if err := tc.open(tt.Context(), c); err != nil {
				tt.Fatalf("failed to open stream: %v", err)
			}
			if streamCtx == nil {
				tt.Fatal("the stream was not opened")
			}
			if err := streamCtx.Err(); err != nil {
				tt.Errorf("the stream ctx is done after the method returned: %v", err)
			}
		})
	}
}
