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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

// failMirrorClient fails every write with a NotFound error naming the id;
// other methods panic through the nil embedded interface.
type failMirrorClient struct {
	service.MirrorClient
}

func failID(id string) error {
	return status.Error(codes.NotFound, "write failed for "+id)
}

func (failMirrorClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Upsert(
	_ context.Context, in *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetId().GetId())
}

// Test_server_multi_joinedErrors pins the error the Multi* write handlers
// return now that per-item failures are collected in a slice and joined once
// after the workers finish: N distinct failures still surface as one status
// carrying the members' code whose message lists each failure exactly once.
func Test_server_multi_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	vec := func(id string) *payload.Object_Vector {
		return &payload.Object_Vector{Id: id, Vector: []float32{1, 2, 3}}
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Insert_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: vec(id)})
				}
				_, err := s.MultiInsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Update_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: vec(id)})
				}
				_, err := s.MultiUpdate(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Upsert_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: vec(id)})
				}
				_, err := s.MultiUpsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Remove_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id}})
				}
				_, err := s.MultiRemove(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			s := &server{
				eg: eg,
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return proxiedFrom // proxied request: forward to the own LB only
					},
					DoFunc: func(ctx context.Context, target string, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
						return f(ctx, target, failMirrorClient{})
					},
				},
			}
			done := make(chan error, 1)
			go func() { done <- test.call(ectx, s) }()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("timed out waiting for the handler")
			}
			if err == nil {
				t.Fatal("returned nil, want the per-item failures")
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != codes.NotFound {
				t.Fatalf("status = %v (ok=%t), want code NotFound: %v", st, ok, err)
			}
			// The joined members are rendered one per line; a flat join of N
			// distinct failures gives exactly N lines.
			if got := strings.Count(st.Message(), "\n") + 1; got != n {
				t.Errorf("message lines = %d, want %d: %q", got, n, st.Message())
			}
			for _, id := range ids {
				if !strings.Contains(st.Message(), "write failed for "+id) {
					t.Errorf("message %q does not mention %q", st.Message(), id)
				}
			}
		})
	}
}
