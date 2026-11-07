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
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

// Test_wrapTargetError pins the code and message wrapTargetError produces for
// each class of target error, which every proxy and do* handler relies on.
func Test_wrapTargetError(t *testing.T) {
	t.Parallel()
	const rpc = vald.InsertRPCName
	tests := []struct {
		err     error
		name    string
		wantMsg string
		want    codes.Code
	}{
		{name: "canceled", err: context.Canceled, want: codes.Canceled, wantMsg: rpc + " API canceld"},
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: codes.DeadlineExceeded, wantMsg: rpc + " API deadline exceeded"},
		{name: "target not found", err: errors.ErrTargetNotFound, want: codes.Internal, wantMsg: rpc + " API target not found"},
		{name: "connection not found", err: errors.ErrGRPCClientConnNotFound("*"), want: codes.Internal, wantMsg: rpc + " API connection not found"},
		{name: "gRPC status passes through", err: status.Error(codes.NotFound, "nf"), want: codes.NotFound, wantMsg: "nf"},
		{name: "non-status error", err: errors.New("boom"), want: codes.Internal, wantMsg: "failed to parse " + rpc + " gRPC error response"},
	}
	s := &server{name: "mirror", ip: "192.0.2.1", vAddr: "lb:8081"}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := wrapTargetError(nil, rpc, test.err,
				&errdetails.RequestInfo{RequestId: "id"}, s.vAddrResourceInfo(rpc))
			st, ok := status.FromError(err)
			if !ok || st == nil {
				t.Fatalf("wrapTargetError() = %v, want a gRPC status error", err)
			}
			if st.Code() != test.want || st.Message() != test.wantMsg {
				t.Errorf("status = (%v, %q), want (%v, %q)", st.Code(), st.Message(), test.want, test.wantMsg)
			}
		})
	}
}

func Test_server_vAddrResourceInfo(t *testing.T) {
	t.Parallel()
	s := &server{name: "mirror", ip: "192.0.2.1", vAddr: "lb:8081"}
	got := s.vAddrResourceInfo(vald.SearchRPCName)
	wantType := errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName
	wantName := apiName + ": mirror(192.0.2.1) to lb:8081"
	if got.GetResourceType() != wantType || got.GetResourceName() != wantName {
		t.Errorf("vAddrResourceInfo() = (%q, %q), want (%q, %q)",
			got.GetResourceType(), got.GetResourceName(), wantType, wantName)
	}
}

// Test_statusOf pins that statusOf passes a usable status through and falls
// back to Unknown with the rpc-specific message only when there is none.
func Test_statusOf(t *testing.T) {
	t.Parallel()
	const rpc = vald.UpdateRPCName
	tests := []struct {
		err     error
		name    string
		wantMsg string
		want    codes.Code
	}{
		{name: "gRPC status", err: status.Error(codes.NotFound, "nf"), want: codes.NotFound, wantMsg: "nf"},
		{name: "plain error keeps its message", err: errors.New("boom"), want: codes.Unknown, wantMsg: "boom"},
		{name: "empty status message falls back", err: status.Error(codes.Internal, ""), want: codes.Unknown, wantMsg: "failed to parse " + rpc + " gRPC error response"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			st := statusOf(test.err, rpc)
			if st == nil {
				t.Fatal("statusOf() = nil")
			}
			if st.Code() != test.want || st.Message() != test.wantMsg {
				t.Errorf("statusOf() = (%v, %q), want (%v, %q)", st.Code(), st.Message(), test.want, test.wantMsg)
			}
		})
	}
}

// Test_server_write_concurrentBroadCast runs every BroadCast callback of a user
// write on its own goroutine, as RangeConcurrent does, so that -race catches an
// unguarded write to the request-local result map.
func Test_server_write_concurrentBroadCast(t *testing.T) {
	t.Parallel()
	const targets = 8
	vec := &payload.Object_Vector{Id: "id-1", Vector: []float32{1, 2, 3}}
	tests := []struct {
		call func(context.Context, *server) (*payload.Object_Location, error)
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Insert(ctx, &payload.Insert_Request{Vector: vec})
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Update(ctx, &payload.Update_Request{Vector: vec})
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{Vector: vec})
			},
		},
		{
			name: "Remove",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Remove(ctx, &payload.Remove_Request{Id: &payload.Object_ID{Id: vec.GetId()}})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string { return "" },
					BroadCastFunc: func(ctx context.Context, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) error) error {
						var wg sync.WaitGroup
						for i := range targets {
							wg.Go(func() {
								if err := f(ctx, "target-"+strconv.Itoa(i), probeMirrorClient{}); err != nil {
									t.Errorf("callback for target-%d = %v, want nil", i, err)
								}
							})
						}
						wg.Wait()
						return nil
					},
				},
			}
			type result struct {
				loc *payload.Object_Location
				err error
			}
			done := make(chan result, 1)
			go func() {
				loc, err := test.call(t.Context(), s)
				done <- result{loc: loc, err: err}
			}()
			select {
			case got := <-done:
				if got.err != nil {
					t.Fatalf("err = %v, want nil", got.err)
				}
				if got.loc.GetUuid() != vec.GetId() {
					t.Errorf("uuid = %q, want %q", got.loc.GetUuid(), vec.GetId())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timed out waiting for the handler")
			}
		})
	}
}
