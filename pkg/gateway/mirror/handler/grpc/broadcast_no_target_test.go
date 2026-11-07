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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

// Test_server_write_broadCastReachedNoTarget pins that a user write whose
// BroadCast returned nil without invoking the callback for any target (every
// target skipped because its pool had no usable conn or its circuit breaker was
// open) is reported as Unavailable instead of OK with an empty result.
func Test_server_write_broadCastReachedNoTarget(t *testing.T) {
	t.Parallel()
	vec := &payload.Object_Vector{Id: "uuid-1", Vector: []float32{1, 2, 3}}
	tests := []struct {
		call func(context.Context, *server) (any, error)
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Insert(ctx, &payload.Insert_Request{Vector: vec})
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Update(ctx, &payload.Update_Request{Vector: vec})
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{Vector: vec})
			},
		},
		{
			name: "Remove",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Remove(ctx, &payload.Remove_Request{Id: &payload.Object_ID{Id: vec.GetId()}})
			},
		},
		{
			name: "RemoveByTimestamp",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var calls int
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return "" // request from a user: broadcast to every target
					},
					BroadCastFunc: func(context.Context, func(context.Context, string, service.MirrorClient, ...grpc.CallOption) error) error {
						calls++
						return nil // every target skipped: callback never invoked
					},
				},
			}
			res, err := test.call(t.Context(), s)
			if calls != 1 {
				t.Fatalf("BroadCast calls = %d, want 1", calls)
			}
			if err == nil {
				t.Fatalf("returned OK (%v) for a write that reached no target, want Unavailable", res)
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != codes.Unavailable {
				t.Fatalf("status = %v (ok=%t), want code Unavailable: %v", st, ok, err)
			}
		})
	}
}
