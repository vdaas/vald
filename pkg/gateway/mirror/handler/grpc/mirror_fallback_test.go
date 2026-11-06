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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

// Test_server_StreamListObject_noReplayOnGatewayFallback models gateway.Do
// falling back to another mirror after the own cluster's listing failed
// mid-stream: the fallback must not open a second listing on the same stream.
func Test_server_StreamListObject_noReplayOnGatewayFallback(t *testing.T) {
	t.Parallel()
	const n = 3
	// bounds the relay loop even if the stream never terminates.
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	failErr := status.New(codes.Unavailable, "agent stream broken").Err()
	own := &listAgentClient{n: n, failErr: failErr}
	peer := &listAgentClient{n: n, failErr: failErr}
	s := &server{
		streamConcurrency: 1,
		gateway: &gatewayMock{
			DoFunc: func(ctx context.Context, target string,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error),
			) (any, error) {
				res, err := f(ctx, target, own)
				if err == nil {
					return res, nil
				}
				return f(ctx, "peer-mirror", peer)
			},
		},
	}
	ss := &listServerStream{ctx: ctx}

	err := s.StreamListObject(&payload.Object_List_Request{}, ss)
	if ctx.Err() != nil {
		t.Fatalf("StreamListObject did not finish before the test deadline: %v", ctx.Err())
	}
	if !status.Is(err, codes.Aborted) {
		t.Errorf("err = %v, want Aborted", err)
	}
	if got := peer.opens.Load(); got != 0 {
		t.Errorf("fallback mirror stream opened %d times, want 0", got)
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.ids) != n {
		t.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, n)
	}
}

// Test_server_StreamListObject_abortedOnNonAbortedFallbackError models the
// production RoundRobin fallback, which returns the error of the last conn it
// visited (ErrTargetNotFound for the own LB address) instead of the callback's
// Aborted: the caller must still get Aborted after a partial delivery.
func Test_server_StreamListObject_abortedOnNonAbortedFallbackError(t *testing.T) {
	t.Parallel()
	const n = 3
	// bounds the relay loop even if the stream never terminates.
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	failErr := status.New(codes.Unavailable, "agent stream broken").Err()
	own := &listAgentClient{n: n, failErr: failErr}
	peer := &listAgentClient{n: n, failErr: failErr}
	s := &server{
		streamConcurrency: 1,
		gateway: &gatewayMock{
			DoFunc: func(ctx context.Context, target string,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error),
			) (any, error) {
				if res, err := f(ctx, target, own); err == nil {
					return res, nil
				}
				if res, err := f(ctx, "peer-mirror", peer); err == nil {
					return res, nil
				}
				return nil, errors.ErrTargetNotFound
			},
		},
	}
	ss := &listServerStream{ctx: ctx}

	err := s.StreamListObject(&payload.Object_List_Request{}, ss)
	if ctx.Err() != nil {
		t.Fatalf("StreamListObject did not finish before the test deadline: %v", ctx.Err())
	}
	if !status.Is(err, codes.Aborted) {
		t.Errorf("err = %v, want Aborted", err)
	}
	if got := peer.opens.Load(); got != 0 {
		t.Errorf("fallback mirror stream opened %d times, want 0", got)
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.ids) != n {
		t.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, n)
	}
}

// Test_server_Update_repairInsertSkipped verifies that Update does not report
// success when the repair Insert for a NotFound cluster was never attempted
// (DoMulti skipped the target, e.g. no healthy connection).
func Test_server_Update_repairInsertSkipped(t *testing.T) {
	t.Parallel()
	const (
		uuid    = "test"
		agentIP = "127.0.0.1"
		target1 = "vald-01"
		target2 = "vald-02"
	)
	loc := &payload.Object_Location{Uuid: uuid, Ips: []string{agentIP}}
	targets := []string{target1, target2}
	cmap := map[string]service.MirrorClient{
		targets[0]: &clientmock.MirrorClientMock{
			UpdateFunc: func(context.Context, *payload.Update_Request, ...grpc.CallOption) (*payload.Object_Location, error) {
				return loc, nil
			},
		},
		targets[1]: &clientmock.MirrorClientMock{
			UpdateFunc: func(context.Context, *payload.Update_Request, ...grpc.CallOption) (*payload.Object_Location, error) {
				return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
			},
		},
	}
	s := &server{
		gateway: &gatewayMock{
			FromForwardedContextFunc: func(context.Context) string { return "" },
			BroadCastFunc: func(ctx context.Context,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error,
			) error {
				for _, tgt := range targets {
					if err := f(ctx, tgt, cmap[tgt]); err != nil {
						t.Logf("BroadCast to %s: %v", tgt, err)
					}
				}
				return nil
			},
			// Like OrderedRangeConcurrent for a target without a healthy
			// connection: f is never called and no error is returned.
			DoMultiFunc: func(context.Context, []string,
				func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error,
			) error {
				return nil
			},
		},
	}

	gotLoc, err := s.Update(t.Context(), &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
		Config: &payload.Update_Config{SkipStrictExistCheck: true},
	})
	if !status.Is(err, codes.NotFound) {
		t.Errorf("err = %v, want NotFound", err)
	}
	if gotLoc != nil {
		t.Errorf("loc = %#v, want nil", gotLoc)
	}
}
