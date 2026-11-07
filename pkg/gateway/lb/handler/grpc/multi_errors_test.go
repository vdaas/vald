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
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Test_server_MultiInsert_joinedErrors pins the error MultiInsert returns now
// that per-item failures are collected in a slice and joined once after the
// workers finish: distinct failures are all kept, identical messages collapse
// to one, and the status code is still derived from the members.
func Test_server_MultiInsert_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 5
	tests := []struct {
		msg         func(id string) string
		name        string
		wantMembers int // 0: a single (non-joined) error
	}{
		{
			name:        "distinct messages are all kept",
			msg:         func(id string) string { return "agent failure: " + id },
			wantMembers: n,
		},
		{
			name: "identical messages collapse to one",
			msg:  func(string) string { return "agent failure" },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := &regClient{
				insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
					return nil, status.Error(codes.NotFound, test.msg(in.GetVector().GetId()))
				},
			}
			s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: n}
			reqs := &payload.Insert_MultiRequest{Requests: make([]*payload.Insert_Request, 0, n)}
			for i := range n {
				reqs.Requests = append(reqs.Requests, &payload.Insert_Request{
					Vector: &payload.Object_Vector{Id: "id-" + strconv.Itoa(i), Vector: []float32{1, 2, 3}},
					Config: &payload.Insert_Config{SkipStrictExistCheck: true},
				})
			}
			_, err := s.MultiInsert(t.Context(), reqs)
			if err == nil {
				t.Fatal("MultiInsert returned nil, want the per-item failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.NotFound {
				t.Errorf("status code = %v (ok=%t), want NotFound: %v", st.Code(), ok, err)
			}
			var u interface{ Unwrap() []error }
			joined := errors.As(err, &u)
			switch {
			case test.wantMembers == 0 && joined:
				t.Errorf("got joined error with %d members, want a single error: %v", len(u.Unwrap()), err)
			case test.wantMembers != 0 && !joined:
				t.Errorf("got %T, want a joined error: %v", err, err)
			case joined && len(u.Unwrap()) != test.wantMembers:
				t.Errorf("members = %d, want %d: %v", len(u.Unwrap()), test.wantMembers, err)
			}
			for i := range n {
				if m := test.msg("id-" + strconv.Itoa(i)); !strings.Contains(err.Error(), m) {
					t.Errorf("error %q does not contain %q", err.Error(), m)
				}
			}
		})
	}
}

// linearClient serves LinearSearch with one result so that every query asking
// for more (MinNum) fails with its own "not enough" NotFound error.
type linearClient struct {
	*regClient
}

func (*linearClient) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

// LinearSearchByID backs the LinearSearchByID fallback with the same single result.
func (*linearClient) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

// Test_server_MultiLinearSearch_joinedErrors pins the error MultiLinearSearch and
// MultiLinearSearchByID return now that per-query failures are collected in a
// slice and joined once after the workers finish: N distinct failures give one
// flat joined error with N members and the members' status code.
func Test_server_MultiLinearSearch_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	// The i-th query asks for MinNum 2+i, so each failure message differs.
	cfgs := make([]*payload.Search_Config, 0, n)
	for minNum := uint32(2); len(cfgs) < n; minNum++ {
		cfgs = append(cfgs, &payload.Search_Config{Num: 10, MinNum: minNum, Timeout: int64(time.Minute)})
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: []float32{1, 2, 3}, Config: cfgs[i]})
				}
				_, err := s.MultiLinearSearch(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: "id-" + strconv.Itoa(i), Config: cfgs[i]})
				}
				_, err := s.MultiLinearSearchByID(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var c vald.Client = &linearClient{regClient: &regClient{
				getObject: func(_ context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
					return &payload.Object_Vector{Id: in.GetId().GetId(), Vector: []float32{1, 2, 3}}, nil
				},
			}}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: n, eg: eg}
			err := test.call(t.Context(), s)
			if err == nil {
				t.Fatal("returned nil, want the per-query failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.NotFound {
				t.Errorf("status code = %v (ok=%t), want NotFound: %v", st.Code(), ok, err)
			}
			var u interface{ Unwrap() []error }
			if !errors.As(err, &u) {
				t.Fatalf("got %T, want a joined error: %v", err, err)
			}
			if got := len(u.Unwrap()); got != n {
				t.Errorf("members = %d, want %d (flat): %v", got, n, err)
			}
			for _, m := range u.Unwrap() {
				var nested interface{ Unwrap() []error }
				if errors.As(m, &nested) {
					t.Errorf("member %T is itself a joined error, want flat: %v", m, m)
				}
			}
			for i := range n {
				if want := "required: " + strconv.Itoa(2+i) + ","; !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err.Error(), want)
				}
			}
		})
	}
}
