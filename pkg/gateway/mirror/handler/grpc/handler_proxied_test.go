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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

// proxiedFrom is what FromForwardedContext reports for a request another mirror
// gateway proxied, which makes the handler forward it to its own LB only.
const proxiedFrom = "mirror-gateway"

// Test_server_proxiedDoError verifies that a proxied Upsert/Remove surfaces the
// error returned by Gateway.Do even when Do fails before invoking the callback
// (e.g. ErrTargetNotFound), instead of returning (nil, nil).
func Test_server_proxiedDoError(t *testing.T) {
	t.Parallel()
	const uuid = "test"
	type test struct {
		call func(ctx context.Context, s *server) (*payload.Object_Location, error)
		name string
	}
	tests := []test{
		{
			name: "Upsert returns error when Do fails without calling f",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{Id: uuid},
				})
			},
		},
		{
			name: "Remove returns error when Do fails without calling f",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Remove(ctx, &payload.Remove_Request{
					Id: &payload.Object_ID{Id: uuid},
				})
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return proxiedFrom // proxied request
					},
					DoFunc: func(context.Context, string, func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
						return nil, errors.ErrTargetNotFound
					},
				},
			}
			loc, err := test.call(tt.Context(), s)
			if err == nil {
				tt.Fatalf("expected an error, got nil (loc: %v)", loc)
			}
			if loc != nil {
				tt.Errorf("expected nil location, got %v", loc)
			}
			if st, ok := status.FromError(err); !ok || st.Code() != codes.Internal {
				tt.Errorf("expected status code %v, got err: %v", codes.Internal, err)
			}
		})
	}
}
