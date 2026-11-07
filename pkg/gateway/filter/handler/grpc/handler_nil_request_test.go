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

	ingressapi "github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/comparator"
)

// The REST filter gateway passes a nil *Multi*Request for an empty or `null`
// body; every Multi* handler must treat it as an empty request list.
func Test_server_multi_nilRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call func(context.Context, *server) (int, error)
		name string
	}{
		{name: vald.MultiSearchObjectRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiSearchObject(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiLinearSearchObjectRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiLinearSearchObject(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiInsertObjectRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiInsertObject(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiUpdateObjectRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiUpdateObject(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiUpsertObjectRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiUpsertObject(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiSearchRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiSearch(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiSearchByIDRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiSearchByID(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiLinearSearchRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiLinearSearch(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiLinearSearchByIDRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiLinearSearchByID(ctx, nil)
			return len(r.GetResponses()), err
		}},
		{name: vald.MultiInsertRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiInsert(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiUpdateRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiUpdate(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiUpsertRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiUpsert(ctx, nil)
			return len(r.GetLocations()), err
		}},
		{name: vald.MultiRemoveRPCName, call: func(ctx context.Context, s *server) (int, error) {
			r, err := s.MultiRemove(ctx, nil)
			return len(r.GetLocations()), err
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			defer func() {
				if r := recover(); r != nil {
					tt.Fatalf("%s(nil) panicked: %v", tc.name, r)
				}
			}()
			// No request is dispatched, so no gateway, filter or errgroup is needed.
			n, err := tc.call(tt.Context(), &server{})
			if err != nil {
				tt.Errorf("%s(nil) error = %v, want nil", tc.name, err)
			}
			if n != 0 {
				tt.Errorf("%s(nil) returned %d results, want 0", tc.name, n)
			}
		})
	}
}

// fixedIngressStub resolves any ingress target to a filter that returns
// filteredVec for every FilterVector call.
type fixedIngressStub struct {
	ingress.Client
}

type fixedIngressFilterStub struct {
	ingressapi.FilterClient
}

func filteredVec() []float32 { return []float32{0.5, 0.25} }

func (*fixedIngressStub) Release(context.Context, ...string) {}

func (*fixedIngressStub) Target(context.Context, ...string) (ingressapi.FilterClient, error) {
	return &fixedIngressFilterStub{}, nil
}

func (*fixedIngressFilterStub) FilterVector(
	context.Context, *payload.Object_Vector, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return &payload.Object_Vector{Vector: filteredVec()}, nil
}

// recordSearchGatewayStub records the vectors forwarded to the LB gateway.
type recordSearchGatewayStub struct {
	client.Client
	vecs [][]float32
	mu   sync.Mutex
}

func (g *recordSearchGatewayStub) record(req *payload.Search_Request) *payload.Search_Response {
	g.mu.Lock()
	g.vecs = append(g.vecs, req.GetVector())
	g.mu.Unlock()
	return okSearchResponse()
}

func (g *recordSearchGatewayStub) Search(
	_ context.Context, req *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.record(req), nil
}

func (g *recordSearchGatewayStub) LinearSearch(
	_ context.Context, req *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.record(req), nil
}

// A nil Search_Request (REST `null` body, or a `null` element of a multi
// request) with a configured search ingress filter must not write through the
// nil pointer; the filtered vector is forwarded instead. In the Multi* path the
// old nil write re-panicked inside an errgroup goroutine and killed the process.
func Test_server_search_nilRequestWithIngressFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{name: "Search", call: func(ctx context.Context, s *server) error {
			_, err := s.Search(ctx, nil)
			return err
		}},
		{name: "LinearSearch", call: func(ctx context.Context, s *server) error {
			_, err := s.LinearSearch(ctx, nil)
			return err
		}},
		{name: "MultiSearch null element", call: func(ctx context.Context, s *server) error {
			r, err := s.MultiSearch(ctx, &payload.Search_MultiRequest{Requests: []*payload.Search_Request{nil}})
			if err == nil && (len(r.GetResponses()) != 1 || r.GetResponses()[0] == nil) {
				return errors.Errorf("MultiSearch responses = %v, want one non-nil response", r.GetResponses())
			}
			return err
		}},
		{name: "MultiLinearSearch null element", call: func(ctx context.Context, s *server) error {
			r, err := s.MultiLinearSearch(ctx, &payload.Search_MultiRequest{Requests: []*payload.Search_Request{nil}})
			if err == nil && (len(r.GetResponses()) != 1 || r.GetResponses()[0] == nil) {
				return errors.Errorf("MultiLinearSearch responses = %v, want one non-nil response", r.GetResponses())
			}
			return err
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			gw := &recordSearchGatewayStub{}
			s := &server{
				gateway:       gw,
				ingress:       &fixedIngressStub{},
				eg:            eg,
				SearchFilters: []string{"ingress-filter:8081"},
			}
			done := make(chan error, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						done <- errors.Errorf("panic: %v", r)
					}
				}()
				done <- tc.call(ectx, s)
			}()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				tt.Fatal("timed out waiting for the handler")
			}
			if err != nil {
				tt.Fatalf("%s error = %v, want nil", tc.name, err)
			}
			gw.mu.Lock()
			defer gw.mu.Unlock()
			if diff := comparator.Diff(gw.vecs, [][]float32{filteredVec()}); diff != "" {
				tt.Errorf("forwarded vectors mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
