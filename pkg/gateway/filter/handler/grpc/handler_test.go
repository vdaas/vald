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
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	egressapi "github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	ingressapi "github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/comparator"
)

// okGatewayStub is an LB gateway whose search/GetObject calls all succeed.
type okGatewayStub struct {
	client.Client
}

func okSearchResponse() *payload.Search_Response {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: testStrictExistUUID, Distance: 0.1}}}
}

func (*okGatewayStub) Search(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) SearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) GetObject(
	context.Context, *payload.Object_VectorRequest, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return &payload.Object_Vector{Id: testStrictExistUUID, Vector: []float32{0.1}}, nil
}

// noTargetIngressStub and noTargetEgressStub fail Target like the real
// clients do when no address is given (and on any other call, which these
// tests never expect).
type noTargetIngressStub struct {
	ingress.Client
}

func (*noTargetIngressStub) Release(context.Context, ...string) {}

func (*noTargetIngressStub) Target(
	_ context.Context, addrs ...string,
) (ingressapi.FilterClient, error) {
	if len(addrs) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	return nil, errors.Errorf("unexpected ingress targets %v", addrs)
}

type noTargetEgressStub struct {
	egress.Client
}

func (*noTargetEgressStub) Release(context.Context, ...string) {}

func (*noTargetEgressStub) Target(
	_ context.Context, addrs ...string,
) (egressapi.FilterClient, error) {
	if len(addrs) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	return nil, errors.Errorf("unexpected egress targets %v", addrs)
}

// An empty (non-nil) targets list, as decoded from REST `"targets": []`, and
// empty configured filters must skip filtering instead of failing with
// "filter targets [] not found".
func Test_server_emptyFilterTargetsSkipped(t *testing.T) {
	t.Parallel()
	empty := func() *payload.Filter_Config {
		return &payload.Filter_Config{Targets: []*payload.Filter_Target{}}
	}
	searchCfg := func() *payload.Search_Config {
		return &payload.Search_Config{IngressFilters: empty(), EgressFilters: empty()}
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "Search",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Search(ctx, &payload.Search_Request{Vector: []float32{0.1}, Config: searchCfg()})
				return err
			},
		},
		{
			name: "SearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: testStrictExistUUID, Config: searchCfg()})
				return err
			},
		},
		{
			name: "LinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearch(ctx, &payload.Search_Request{Vector: []float32{0.1}, Config: searchCfg()})
				return err
			},
		},
		{
			name: "LinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: testStrictExistUUID, Config: searchCfg()})
				return err
			},
		},
		{
			name: "GetObject",
			call: func(ctx context.Context, s *server) error {
				_, err := s.GetObject(ctx, &payload.Object_VectorRequest{
					Id:      &payload.Object_ID{Id: testStrictExistUUID},
					Filters: empty(),
				})
				return err
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{
				gateway:         &okGatewayStub{},
				ingress:         &noTargetIngressStub{},
				egress:          &noTargetEgressStub{},
				SearchFilters:   []string{},
				DistanceFilters: []string{},
				ObjectFilters:   []string{},
			}
			if err := tc.call(tt.Context(), s); err != nil {
				tt.Errorf("%s() error = %v, want nil", tc.name, err)
			}
		})
	}
}

// Test_filterAddrs pins that filterAddrs lists the configured addrs first and
// formats each request target exactly as fmt.Sprintf("%s:%d") did (an IPv6
// host is not bracketed), with no spare capacity left to regrow into.
func Test_filterAddrs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		configured []string
		targets    []*payload.Filter_Target
	}{
		{name: "nothing"},
		{name: "configured only", configured: []string{"a:1", "b:2"}},
		{
			name: "targets only",
			targets: []*payload.Filter_Target{
				{Host: "filter", Port: 8081},
				{Host: "::1", Port: 0},
				{Host: "", Port: math.MaxUint32},
				nil,
			},
		},
		{
			name:       "configured then targets",
			configured: []string{"configured:8081"},
			targets:    []*payload.Filter_Target{{Host: "10.0.0.1", Port: 80}, {Host: "fe80::1%eth0", Port: 65535}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			want := slices.Clone(tc.configured)
			for _, target := range tc.targets {
				want = append(want, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
			}
			got := filterAddrs(tc.configured, tc.targets)
			if !slices.Equal(got, want) {
				tt.Errorf("filterAddrs() = %q, want %q", got, want)
			}
			if cap(got) != len(want) {
				tt.Errorf("cap(filterAddrs()) = %d, want %d", cap(got), len(want))
			}
		})
	}
}

// lbErrGatewayStub is an LB gateway whose write/GetObject calls all fail with err.
type lbErrGatewayStub struct {
	client.Client
	err error
}

func (g *lbErrGatewayStub) Insert(
	context.Context, *payload.Insert_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) Update(
	context.Context, *payload.Update_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) Upsert(
	context.Context, *payload.Upsert_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) GetObject(
	context.Context, *payload.Object_VectorRequest, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return nil, g.err
}

// passIngressStub is an ingress filter that returns the vector unchanged.
type passIngressStub struct {
	ingress.Client
}

func (*passIngressStub) Release(context.Context, ...string) {}

func (p *passIngressStub) Target(context.Context, ...string) (ingressapi.FilterClient, error) {
	return p, nil
}

func (*passIngressStub) FilterVector(
	_ context.Context, in *payload.Object_Vector, _ ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return in, nil
}

// The filter gateway must preserve the LB gateway's status code after ingress
// filtering (and on GetObject) instead of rewriting it to Internal / NotFound.
func Test_server_lbErrorCodePreserved(t *testing.T) {
	t.Parallel()
	vec := func() *payload.Object_Vector {
		return &payload.Object_Vector{Id: "uuid-1", Vector: []float32{0.1, 0.2, 0.3, 0.4}}
	}
	calls := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Insert(ctx, &payload.Insert_Request{
					Vector: vec(),
					Config: &payload.Insert_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: vec()})
				return err
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{Vector: vec()})
				return err
			},
		},
		{
			name: "GetObject",
			call: func(ctx context.Context, s *server) error {
				_, err := s.GetObject(ctx, &payload.Object_VectorRequest{Id: &payload.Object_ID{Id: "uuid-1"}})
				return err
			},
		},
	}
	errs := []struct {
		lbErr    error
		name     string
		wantCode codes.Code
	}{
		{name: "NotFound", lbErr: status.Error(codes.NotFound, "not found"), wantCode: codes.NotFound},
		{name: "AlreadyExists", lbErr: status.Error(codes.AlreadyExists, "exists"), wantCode: codes.AlreadyExists},
		{name: "Unavailable", lbErr: status.Error(codes.Unavailable, "lb unavailable"), wantCode: codes.Unavailable},
		{name: "non-status error", lbErr: errors.New("connection not found"), wantCode: codes.Internal},
	}
	for _, c := range calls {
		for _, e := range errs {
			t.Run(c.name+"/"+e.name, func(tt *testing.T) {
				tt.Parallel()
				filters := []string{"ingress:8081"}
				s := &server{
					gateway:       &lbErrGatewayStub{err: e.lbErr},
					ingress:       &passIngressStub{},
					InsertFilters: filters,
					UpdateFilters: filters,
					UpsertFilters: filters,
				}
				err := c.call(tt.Context(), s)
				if err == nil {
					tt.Fatal("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok || st == nil {
					tt.Fatalf("error is not a gRPC status: %v", err)
				}
				if st.Code() != e.wantCode {
					tt.Errorf("code = %v, want %v (err: %v)", st.Code(), e.wantCode, err)
				}
			})
		}
	}
}

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

// linearSearchGatewayStub answers LinearSearch with a fixed response.
type linearSearchGatewayStub struct {
	client.Client
	res *payload.Search_Response
}

func (g *linearSearchGatewayStub) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.res, nil
}

// linearSearchStreamStub feeds reqs to Recv (then io.EOF) and records Send.
type linearSearchStreamStub struct {
	vald.Search_StreamLinearSearchServer
	ctx  func() context.Context
	reqs []*payload.Search_Request
	sent []*payload.Search_StreamResponse
	mu   sync.Mutex
}

func (s *linearSearchStreamStub) Context() context.Context { return s.ctx() }

func (s *linearSearchStreamStub) Recv() (*payload.Search_Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.reqs) == 0 {
		return nil, io.EOF
	}
	req := s.reqs[0]
	s.reqs = s.reqs[1:]
	return req, nil
}

func (s *linearSearchStreamStub) Send(res *payload.Search_StreamResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sent = append(s.sent, res)
	return nil
}

// TestStreamLinearSearch_nilConfig verifies that a stream message without a
// Config (legal in proto3) is served instead of nil-dereferencing
// req.Config.RequestId, which panicked inside the stream worker and ended the
// whole stream with an error.
func TestStreamLinearSearch_nilConfig(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	want := &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "a", Distance: 0.1}}}
	s := &server{gateway: &linearSearchGatewayStub{res: want}}
	stream := &linearSearchStreamStub{
		ctx:  func() context.Context { return ctx },
		reqs: []*payload.Search_Request{{Vector: []float32{1, 2}}},
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.StreamLinearSearch(stream)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("StreamLinearSearch() error = %v, want nil", err)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for StreamLinearSearch")
	}

	stream.mu.Lock()
	defer stream.mu.Unlock()
	if len(stream.sent) != 1 {
		t.Fatalf("sent %d responses, want 1", len(stream.sent))
	}
	if got := stream.sent[0].GetResponse(); got != want {
		t.Errorf("response = %v, want %v", got, want)
	}
}

// strictExistGatewayStub records what the filter gateway forwards to the LB gateway.
type strictExistGatewayStub struct {
	client.Client
	existsID    *payload.Object_ID
	existsErr   error
	insertReq   *payload.Insert_Request
	updateReq   *payload.Update_Request
	upsertReq   *payload.Upsert_Request
	existsCalls int
}

func (g *strictExistGatewayStub) Exists(
	_ context.Context, _ *payload.Object_ID, _ ...grpc.CallOption,
) (*payload.Object_ID, error) {
	g.existsCalls++
	return g.existsID, g.existsErr
}

func (g *strictExistGatewayStub) Insert(
	_ context.Context, req *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.insertReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

func (g *strictExistGatewayStub) Update(
	_ context.Context, req *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.updateReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

func (g *strictExistGatewayStub) Upsert(
	_ context.Context, req *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.upsertReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

const testStrictExistUUID = "uuid-1"

func testStrictExistVector() *payload.Object_Vector {
	return &payload.Object_Vector{Id: testStrictExistUUID, Vector: []float32{0.1, 0.2, 0.3, 0.4}}
}

func Test_server_Insert_strictExistCheck(t *testing.T) {
	t.Parallel()
	tests := []struct {
		existsErr   error
		existsID    *payload.Object_ID
		name        string
		wantCode    codes.Code
		wantForward bool
	}{
		{
			name:        "not found proceeds to insert",
			existsErr:   status.Error(codes.NotFound, "not found"),
			wantCode:    codes.OK,
			wantForward: true,
		},
		{
			name:     "existing id returns AlreadyExists",
			existsID: &payload.Object_ID{Id: testStrictExistUUID},
			wantCode: codes.AlreadyExists,
		},
		{
			name:      "unavailable Exists error is returned, not bypassed",
			existsErr: status.Error(codes.Unavailable, "lb gateway unavailable"),
			wantCode:  codes.Unavailable,
		},
		{
			name:      "deadline exceeded Exists error is returned, not bypassed",
			existsErr: context.DeadlineExceeded,
			wantCode:  codes.DeadlineExceeded,
		},
		{
			name:      "non-status Exists error is returned as Internal",
			existsErr: errors.New("connection not found"),
			wantCode:  codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			gw := &strictExistGatewayStub{existsID: tc.existsID, existsErr: tc.existsErr}
			s := &server{gateway: gw}
			_, err := s.Insert(tt.Context(), &payload.Insert_Request{Vector: testStrictExistVector()})
			gotCode := codes.OK
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st == nil {
					tt.Fatalf("error is not a gRPC status: %v", err)
				}
				gotCode = st.Code()
			}
			if gotCode != tc.wantCode {
				tt.Errorf("code = %v, want %v (err: %v)", gotCode, tc.wantCode, err)
			}
			if forwarded := gw.insertReq != nil; forwarded != tc.wantForward {
				tt.Errorf("forwarded = %v, want %v", forwarded, tc.wantForward)
			}
		})
	}
}

// Update/Upsert of an existing ID must reach the LB gateway, which performs
// the strict exist check (same-vector comparison) itself.
func Test_server_UpdateUpsert_existingIDForwarded(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call func(*server, context.Context, *strictExistGatewayStub) (forwarded, skip bool, err error)
		name string
	}{
		{
			name: "Update",
			call: func(s *server, ctx context.Context, gw *strictExistGatewayStub) (bool, bool, error) {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: testStrictExistVector()})
				return gw.updateReq != nil, gw.updateReq.GetConfig().GetSkipStrictExistCheck(), err
			},
		},
		{
			name: "Upsert",
			call: func(s *server, ctx context.Context, gw *strictExistGatewayStub) (bool, bool, error) {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{Vector: testStrictExistVector()})
				return gw.upsertReq != nil, gw.upsertReq.GetConfig().GetSkipStrictExistCheck(), err
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			gw := &strictExistGatewayStub{existsID: &payload.Object_ID{Id: testStrictExistUUID}}
			s := &server{gateway: gw}
			forwarded, skip, err := tc.call(s, tt.Context(), gw)
			if err != nil {
				tt.Fatalf("unexpected error: %v", err)
			}
			if !forwarded {
				tt.Fatal("request was not forwarded to the LB gateway")
			}
			if skip {
				tt.Error("SkipStrictExistCheck was forced to true; LB strict check is bypassed")
			}
		})
	}
}

// vectorizerIngressStub records the addresses passed to Target and vectorizes
// every object to genVec.
type vectorizerIngressStub struct {
	ingress.Client
	targets  [][]string
	released [][]string
	mu       sync.Mutex
}

type vectorizerFilterStub struct {
	ingressapi.FilterClient
}

func genVec() []float32 { return []float32{0.1, 0.2} }

const (
	confVectorizer = "conf:8080"
	testVectorizer = "vectorizer:8080"
)

func (v *vectorizerIngressStub) Target(
	_ context.Context, addrs ...string,
) (ingressapi.FilterClient, error) {
	v.mu.Lock()
	v.targets = append(v.targets, addrs)
	v.mu.Unlock()
	return &vectorizerFilterStub{}, nil
}

func (v *vectorizerIngressStub) Release(_ context.Context, addrs ...string) {
	v.mu.Lock()
	v.released = append(v.released, addrs)
	v.mu.Unlock()
}

func (*vectorizerFilterStub) GenVector(
	context.Context, *payload.Object_Blob, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return &payload.Object_Vector{Vector: genVec()}, nil
}

func Test_server_vectorizerAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		vr         *payload.Filter_Target
		name       string
		configured string
		want       string
	}{
		{name: "request target wins over config", vr: &payload.Filter_Target{Host: "vec", Port: 9000}, configured: confVectorizer, want: "vec:9000"},
		{name: "request target without host defaults to localhost", vr: &payload.Filter_Target{Port: 9000}, want: "localhost:9000"},
		{name: "nil request target falls back to config", configured: confVectorizer, want: confVectorizer},
		{name: "request target without port falls back to config", vr: &payload.Filter_Target{Host: "vec"}, configured: confVectorizer, want: confVectorizer},
		{name: "neither set", vr: &payload.Filter_Target{Host: "portless"}, want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{Vectorizer: tc.configured}
			if got := s.vectorizerAddr(tc.vr); got != tc.want {
				tt.Errorf("vectorizerAddr() = %q, want %q", got, tc.want)
			}
		})
	}
}

// The *Object RPCs must use the configured ingress_filter.vectorizer when the
// request names no vectorizer, instead of rejecting the request.
func Test_server_searchObject_configuredVectorizer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call       func(context.Context, *server, *payload.Search_ObjectRequest) error
		name       string
		configured string
		wantTarget [][]string
		wantCode   codes.Code
	}{
		{
			name: "SearchObject uses the configured vectorizer",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.SearchObject(ctx, req)
				return err
			},
			configured: testVectorizer,
			wantCode:   codes.OK,
			wantTarget: [][]string{{testVectorizer}},
		},
		{
			name: "LinearSearchObject uses the configured vectorizer",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.LinearSearchObject(ctx, req)
				return err
			},
			configured: testVectorizer,
			wantCode:   codes.OK,
			wantTarget: [][]string{{testVectorizer}},
		},
		{
			name: "SearchObject without any vectorizer is rejected",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.SearchObject(ctx, req)
				return err
			},
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			ing := &vectorizerIngressStub{}
			gw := &recordSearchGatewayStub{}
			s := &server{
				gateway:    gw,
				ingress:    ing,
				Vectorizer: tc.configured,
			}
			done := make(chan error, 1)
			go func() {
				done <- tc.call(ctx, s, &payload.Search_ObjectRequest{Object: []byte("obj")})
			}()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				tt.Fatal("timed out waiting for the handler")
			}
			code := codes.OK
			if err != nil {
				st, _ := status.FromError(err)
				code = st.Code()
			}
			if code != tc.wantCode {
				tt.Fatalf("code = %v, want %v (err = %v)", code, tc.wantCode, err)
			}
			ing.mu.Lock()
			defer ing.mu.Unlock()
			if diff := comparator.Diff(ing.targets, tc.wantTarget); diff != "" {
				tt.Errorf("Target addrs mismatch (-got +want):\n%s", diff)
			}
			// Every target resolved for the request is released when it ends.
			if diff := comparator.Diff(ing.released, tc.wantTarget); diff != "" {
				tt.Errorf("Release addrs mismatch (-got +want):\n%s", diff)
			}
			if tc.wantCode != codes.OK {
				return
			}
			gw.mu.Lock()
			defer gw.mu.Unlock()
			if diff := comparator.Diff(gw.vecs, [][]float32{genVec()}); diff != "" {
				tt.Errorf("forwarded vectors mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

// fanOutProbe records how many calls to hold overlap. A call stays in flight
// until more than limit calls overlap (then every call is released at once, so
// an unbounded fan-out is seen without waiting) or a short grace period ends
// (so a bounded fan-out still makes progress).
type fanOutProbe struct {
	over   chan struct{}
	mu     sync.Mutex
	limit  int
	cur    int
	peak   int
	calls  int
	closed bool
}

func newFanOutProbe(limit int) *fanOutProbe {
	return &fanOutProbe{over: make(chan struct{}), limit: limit}
}

func (p *fanOutProbe) hold() {
	p.mu.Lock()
	p.cur++
	p.calls++
	p.peak = max(p.peak, p.cur)
	if p.cur > p.limit && !p.closed {
		p.closed = true
		close(p.over)
	}
	p.mu.Unlock()
	timer := time.NewTimer(20 * time.Millisecond)
	select {
	case <-p.over:
	case <-timer.C:
	}
	timer.Stop()
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
}

func (p *fanOutProbe) stats() (peak, calls int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak, p.calls
}

// probeGateway answers every forwarded call through the probe, echoing the
// item's request id (search) or uuid (writes); the item "req-0" / "id-0"
// panics instead when panics is set. Other methods panic through the nil
// embedded interface.
type probeGateway struct {
	client.Client
	p      *fanOutProbe
	panics bool
}

func (g *probeGateway) answer(key string) {
	g.p.hold()
	if g.panics && (key == "req-0" || key == "id-0") {
		panic("probe panic")
	}
}

func (g *probeGateway) search(cfg *payload.Search_Config) (*payload.Search_Response, error) {
	g.answer(cfg.GetRequestId())
	return &payload.Search_Response{RequestId: cfg.GetRequestId()}, nil
}

func (g *probeGateway) location(uuid string) (*payload.Object_Location, error) {
	g.answer(uuid)
	return &payload.Object_Location{Uuid: uuid}, nil
}

func (g *probeGateway) Search(
	_ context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.search(in.GetConfig())
}

func (g *probeGateway) SearchByID(
	_ context.Context, in *payload.Search_IDRequest, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.search(in.GetConfig())
}

func (g *probeGateway) LinearSearch(
	_ context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.search(in.GetConfig())
}

func (g *probeGateway) LinearSearchByID(
	_ context.Context, in *payload.Search_IDRequest, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.search(in.GetConfig())
}

func (g *probeGateway) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return g.location(in.GetId().GetId())
}

func (g *probeGateway) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return g.location(in.GetVector().GetId())
}

func (g *probeGateway) Upsert(
	_ context.Context, in *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return g.location(in.GetVector().GetId())
}

// Test_server_Multi_boundedFanOut pins that the Multi* handlers run at most
// streamConcurrency items at once, still answer every item in request order,
// and report a panicking item as an error instead of dropping it.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	t.Parallel()
	const (
		n     = 6
		limit = 2
	)
	reqID := func(i int) string { return "req-" + strconv.Itoa(i) }
	id := func(i int) string { return "id-" + strconv.Itoa(i) }
	cfg := func(i int) *payload.Search_Config { return &payload.Search_Config{RequestId: reqID(i)} }
	vec := func(i int) *payload.Object_Vector {
		return &payload.Object_Vector{Id: id(i), Vector: []float32{1, 2, 3}}
	}
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: []float32{1, 2, 3}, Config: cfg(i)})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: id(i), Config: cfg(i)})
		}
		return reqs
	}
	searchKeys := func(res *payload.Search_Responses) (keys []string) {
		for _, r := range res.GetResponses() {
			keys = append(keys, r.GetRequestId())
		}
		return keys
	}
	locKeys := func(res *payload.Object_Locations) (keys []string) {
		for _, l := range res.GetLocations() {
			keys = append(keys, l.GetUuid())
		}
		return keys
	}
	tests := []struct {
		call func(context.Context, *server) ([]string, error)
		want func(int) string
		name string
	}{
		{
			name: vald.MultiSearchRPCName,
			want: reqID,
			call: func(ctx context.Context, s *server) ([]string, error) {
				res, err := s.MultiSearch(ctx, searchReqs())
				return searchKeys(res), err
			},
		},
		{
			name: vald.MultiSearchByIDRPCName,
			want: reqID,
			call: func(ctx context.Context, s *server) ([]string, error) {
				res, err := s.MultiSearchByID(ctx, idReqs())
				return searchKeys(res), err
			},
		},
		{
			name: vald.MultiLinearSearchRPCName,
			want: reqID,
			call: func(ctx context.Context, s *server) ([]string, error) {
				res, err := s.MultiLinearSearch(ctx, searchReqs())
				return searchKeys(res), err
			},
		},
		{
			name: vald.MultiLinearSearchByIDRPCName,
			want: reqID,
			call: func(ctx context.Context, s *server) ([]string, error) {
				res, err := s.MultiLinearSearchByID(ctx, idReqs())
				return searchKeys(res), err
			},
		},
		{
			name: vald.MultiRemoveRPCName,
			want: id,
			call: func(ctx context.Context, s *server) ([]string, error) {
				reqs := &payload.Remove_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id(i)}})
				}
				res, err := s.MultiRemove(ctx, reqs)
				return locKeys(res), err
			},
		},
		{
			name: vald.MultiUpdateRPCName,
			want: id,
			call: func(ctx context.Context, s *server) ([]string, error) {
				reqs := &payload.Update_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: vec(i)})
				}
				res, err := s.MultiUpdate(ctx, reqs)
				return locKeys(res), err
			},
		},
		{
			name: vald.MultiUpsertRPCName,
			want: id,
			call: func(ctx context.Context, s *server) ([]string, error) {
				reqs := &payload.Upsert_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: vec(i)})
				}
				res, err := s.MultiUpsert(ctx, reqs)
				return locKeys(res), err
			},
		},
	}
	for _, test := range tests {
		for _, panics := range []bool{false, true} {
			t.Run(test.name+"/panics="+strconv.FormatBool(panics), func(t *testing.T) {
				t.Parallel()
				p := newFanOutProbe(limit)
				// eg is the shared group the handlers used to fan out on.
				eg, _ := errgroup.New(t.Context())
				s := &server{
					gateway:           &probeGateway{p: p, panics: panics},
					eg:                eg,
					streamConcurrency: limit,
				}
				type result struct {
					err  error
					keys []string
				}
				done := make(chan result, 1)
				go func() {
					keys, err := test.call(t.Context(), s)
					done <- result{keys: keys, err: err}
				}()
				var got result
				select {
				case got = <-done:
				case <-time.After(10 * time.Second):
					t.Fatal("timed out waiting for the handler")
				}
				if peak, _ := p.stats(); peak > limit {
					t.Errorf("peak concurrency = %d, want <= %d", peak, limit)
				}
				if panics {
					// The handlers may re-wrap the recovered panic into a status,
					// so only its presence is pinned here.
					if got.err == nil {
						t.Error("err = nil, want the recovered panic")
					}
					// The panicking item 0 must not make the group skip the queued
					// items: every item ran and every other item has its result.
					if _, calls := p.stats(); calls != n {
						t.Errorf("forwarded %d items, want all %d", calls, n)
					}
					if len(got.keys) != n {
						t.Fatalf("got %d results, want %d", len(got.keys), n)
					}
					for i := 1; i < n; i++ {
						if want := test.want(i); got.keys[i] != want {
							t.Errorf("results[%d] = %q, want %q", i, got.keys[i], want)
						}
					}
					return
				}
				if got.err != nil {
					t.Fatalf("unexpected error: %v", got.err)
				}
				if len(got.keys) != n {
					t.Fatalf("got %d results, want %d", len(got.keys), n)
				}
				for i, k := range got.keys {
					if want := test.want(i); k != want {
						t.Errorf("results[%d] = %q, want %q", i, k, want)
					}
				}
			})
		}
	}
}

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

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want vald.ServerWithFilter
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, vald.ServerWithFilter) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got vald.ServerWithFilter) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := New(test.args.opts...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_SearchObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_ObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.SearchObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearchObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearchObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamSearchObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamSearchObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_ObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Search_Response
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.LinearSearchObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamSearchObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamLinearSearchObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_InsertObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_ObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.InsertObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamInsertObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamInsertObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamInsertObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiInsertObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiInsertObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_UpdateObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_ObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.UpdateObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpdateObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamUpdateObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamUpdateObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpdateObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiUpdateObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_UpsertObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_ObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.UpsertObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpsertObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamUpsertObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamUpsertObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpsertObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiObjectRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiUpsertObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Exists(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		meta *payload.Object_ID
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Object_ID
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_ID, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           meta:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           meta:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Search(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.Search(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_SearchByID(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamLinearSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamLinearSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Insert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.Insert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamInsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Insert_StreamInsertServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamInsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiInsert(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiInsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Update(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.Update(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamUpdate(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpdate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Upsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Upsert_StreamUpsertServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamUpsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpsert(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Remove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.Remove(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamRemove(t *testing.T) {
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamRemove(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiRemove(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Remove_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.MultiRemove(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Flush(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Flush_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantLoc *payload.Info_Index_Count
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Count, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Info_Index_Count, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.Flush(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_RemoveByTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_TimestampRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		want *payload.Object_Locations
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.RemoveByTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_GetObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_VectorRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		wantVec *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			gotVec, err := s.GetObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamGetObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Object_StreamGetObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		eg                                errgroup.Group
// 		Vectorizer                        string
// 		defaultVectorizer                 string
// 		name                              string
// 		ip                                string
// 		copts                             []grpc.CallOption
// 		DistanceFilters                   []string
// 		ObjectFilters                     []string
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		defaultFilters                    []string
// 		streamConcurrency                 int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServerWithFilter:nil,
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           eg:nil,
// 		           Vectorizer:"",
// 		           defaultVectorizer:"",
// 		           name:"",
// 		           ip:"",
// 		           copts:nil,
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           defaultFilters:nil,
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				eg:                                test.fields.eg,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				copts:                             test.fields.copts,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamGetObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
