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
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// failGatewayStub fails every forwarded search and remove with the same
// error; the filter handler wraps each one with a per-item message. Other
// methods panic through the nil embedded interface.
type failGatewayStub struct {
	client.Client
}

var errFailGateway = status.Error(codes.Internal, "lb gateway failure")

func (failGatewayStub) Search(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return nil, errFailGateway
}

func (failGatewayStub) SearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return nil, errFailGateway
}

func (failGatewayStub) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return nil, errFailGateway
}

func (failGatewayStub) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return nil, errFailGateway
}

func (failGatewayStub) Remove(
	context.Context, *payload.Remove_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, errFailGateway
}

// Test_server_multi_joinedErrors pins the error the Multi* handlers return now
// that per-item failures are collected in a slice and joined once after the
// workers finish: N distinct per-item failures give one flat joined error with
// N members, whose status code is the members' NotFound.
func Test_server_multi_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: []float32{float32(i), 1, 2}})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: id})
		}
		return reqs
	}
	tests := []struct {
		call     func(context.Context, *server) error
		name     string
		wantMsgs []string
	}{
		{
			name: vald.MultiSearchRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearch(ctx, searchReqs())
				return err
			},
		},
		{
			name: vald.MultiSearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearchByID(ctx, idReqs())
				return err
			},
			wantMsgs: ids,
		},
		{
			name: vald.MultiLinearSearchRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearch(ctx, searchReqs())
				return err
			},
		},
		{
			name: vald.MultiLinearSearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearchByID(ctx, idReqs())
				return err
			},
			wantMsgs: ids,
		},
		{
			name: vald.MultiRemoveRPCName,
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Remove_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id}})
				}
				_, err := s.MultiRemove(ctx, reqs)
				return err
			},
			wantMsgs: ids,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			s := &server{gateway: failGatewayStub{}, eg: eg}
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
			for _, m := range test.wantMsgs {
				// MultiRemove formats the whole Object_ID (id:"<id>"), the others the bare id.
				if msg := err.Error(); !strings.Contains(msg, " "+m+"'s ") && !strings.Contains(msg, `"`+m+`"'s `) {
					t.Errorf("error %q does not mention %q", err.Error(), m)
				}
			}
		})
	}
}
