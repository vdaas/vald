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
	"runtime"
	"slices"
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
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"go.opentelemetry.io/otel/attribute"
)

// extClient adds the agent calls okClient does not stub.
type extClient struct {
	*okClient
	noIps bool // answer UpdateTimestamp with a location that carries no ips
}

func (c extClient) UpdateTimestamp(
	context.Context, *payload.Update_TimestampRequest, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	if c.noIps {
		return &payload.Object_Location{}, nil
	}
	return &payload.Object_Location{Uuid: fixtureUUID, Name: "agent-ts", Ips: []string{"10.0.0.9"}}, nil
}

func (c extClient) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "id", Distance: 0.1}}}, nil
}

func newExtClient(vec []float32) extClient {
	c := extClient{okClient: newOKClient(vec)}
	c.search = func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
		return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "id", Distance: 0.1}}}, nil
	}
	return c
}

// The error details (full request JSON, resource names, gateway addresses) are only
// needed when an RPC fails. A successful RPC must not build them: gateway.Addrs is
// queried only while building them, so it must not be called at all.
func Test_server_successDoesNotBuildErrorDetails(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	vec := []float32{1, 2, 3}
	tests := []struct {
		call func(ctx context.Context, s *server) error
		name string
	}{
		{
			name: vald.InsertRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Insert(ctx, &payload.Insert_Request{Vector: &payload.Object_Vector{Id: uuid, Vector: vec}})
				return err
			},
		},
		{
			name: vald.UpsertRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{Id: uuid, Vector: vec},
					Config: &payload.Upsert_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: vald.UpdateRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: &payload.Object_Vector{Id: uuid, Vector: vec}})
				return err
			},
		},
		{
			name: vald.UpdateTimestampRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{Id: uuid, Timestamp: 1})
				return err
			},
		},
		{
			name: vald.RemoveRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Remove(ctx, &payload.Remove_Request{
					Id:     &payload.Object_ID{Id: uuid},
					Config: &payload.Remove_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: vald.SearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: uuid, Config: &payload.Search_Config{Num: 3}})
				return err
			},
		},
		{
			name: vald.LinearSearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: uuid, Config: &payload.Search_Config{Num: 3}})
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one so Update proceeds to remove+insert,
			// and Exists finds nothing so Insert may proceed.
			c := newExtClient([]float32{4, 5, 6})
			c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
				return &payload.Object_ID{}, nil
			}
			gw := syncGateway(c)
			if test.name == "Upsert" || test.name == "Update" {
				// Upsert/Update look the object up first; the stored vector must be found.
				c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
					return &payload.Object_ID{Id: uuid}, nil
				}
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, replica: 1, eg: eg, timeout: time.Second}
			if err := test.call(t.Context(), s); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n := gw.addrsCalls.Load(); n != 0 {
				t.Errorf("a successful %s built the error details: gateway.Addrs called %d times, want 0", test.name, n)
			}
		})
	}
}

// statusSig reduces an error to its status code, message and details. The details are
// sorted because the status builder emits them in map order.
func statusSig(t *testing.T, err error) string {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("not a status error: %v", err)
	}
	ds := make([]string, 0, len(st.Proto().GetDetails()))
	for _, d := range st.Proto().GetDetails() {
		ds = append(ds, d.GetTypeUrl()+"|"+string(d.GetValue()))
	}
	slices.Sort(ds)
	return fmt.Sprintf("%v|%s|%q", st.Code(), st.Message(), ds)
}

// The details of every error the lazily built RequestInfo/ResourceInfo feed into must be
// exactly what the handlers used to build up front: expected values are assembled here
// from the same inputs the old code used (request JSON, resource type, name, gateway addrs).
func Test_server_errorDetailsUnchanged(t *testing.T) {
	t.Parallel()
	const (
		uuid = "uuid"
		name = "lb-name"
		ip   = "10.0.0.1"
	)
	addrs := []string{"10.0.0.2:8081", "10.0.0.3:8081"}
	vec := []float32{1, 2, 3}
	invalidID := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequestFieldViolation{
			{Field: "invalid id", Description: errors.ErrInvalidMetaDataConfig.Error()},
		},
	}
	insertReq := &payload.Insert_Request{Vector: &payload.Object_Vector{Vector: vec}}
	upsertReq := &payload.Upsert_Request{Vector: &payload.Object_Vector{Vector: vec}}
	updateReq := &payload.Update_Request{Vector: &payload.Object_Vector{Vector: vec}}
	tsReq := &payload.Update_TimestampRequest{Timestamp: 1}
	tsNotFoundReq := &payload.Update_TimestampRequest{Id: uuid, Timestamp: 1}
	searchByIDReq := &payload.Search_IDRequest{Config: &payload.Search_Config{Num: 3}}
	removeReq := &payload.Remove_Request{
		Id:     &payload.Object_ID{Id: uuid},
		Config: &payload.Remove_Config{SkipStrictExistCheck: true},
	}
	tests := []struct {
		call         func(ctx context.Context, s *server) error
		req          any
		wrap         func(msg string, err error, details ...any) error
		cause        error
		name         string
		msg          string
		uuid         string
		resourceType string
		extra        []any
		noIps        bool
	}{
		{
			name: "Insert invalid uuid", req: insertReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Insert(ctx, insertReq); return err },
			msg:  vald.InsertRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.InsertRPCName, extra: []any{invalidID},
		},
		{
			name: "Upsert invalid uuid", req: upsertReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Upsert(ctx, upsertReq); return err },
			msg:  vald.UpsertRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpsertRPCName, extra: []any{invalidID},
		},
		{
			name: "Update invalid uuid", req: updateReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Update(ctx, updateReq); return err },
			msg:  vald.UpdateRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName, extra: []any{invalidID},
		},
		{
			name: "UpdateTimestamp invalid uuid", req: tsReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.UpdateTimestamp(ctx, tsReq); return err },
			msg:  vald.UpdateTimestampRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName, extra: []any{invalidID},
		},
		{
			name: "SearchByID invalid uuid", req: searchByIDReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.SearchByID(ctx, searchByIDReq); return err },
			msg:  vald.SearchByIDRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.SearchByIDRPCName, extra: []any{invalidID},
		},
		{
			name: "LinearSearchByID invalid uuid", req: searchByIDReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, searchByIDReq)
				return err
			},
			msg: vald.LinearSearchByIDRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.LinearSearchByIDRPCName, extra: []any{invalidID},
		},
		{
			name: "Remove target not found", req: removeReq, uuid: uuid, wrap: status.WrapWithNotFound,
			call: func(ctx context.Context, s *server) error { _, err := s.Remove(ctx, removeReq); return err },
			msg:  vald.RemoveRPCName + " API remove target not found", cause: errors.ErrIndexNotFound,
			resourceType: "/vald.v1." + vald.RemoveRPCName,
		},
		{
			name: "UpdateTimestamp target not found", req: tsNotFoundReq, uuid: uuid, wrap: status.WrapWithNotFound, noIps: true,
			call: func(ctx context.Context, s *server) error {
				_, err := s.UpdateTimestamp(ctx, tsNotFoundReq)
				return err
			},
			msg: vald.UpdateTimestampRPCName + " API update target not found", cause: errors.ErrIndexNotFound,
			resourceType: "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newExtClient(vec)
			c.noIps = test.noIps
			c.remove = func(context.Context, *payload.Remove_Request) (*payload.Object_Location, error) {
				return &payload.Object_Location{}, nil
			}
			gw := syncGateway(c)
			gw.addrs = addrs
			s := &server{gateway: gw, name: name, ip: ip, replica: 1, timeout: time.Second}

			err := test.call(t.Context(), s)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			details := append([]any{
				&errdetails.RequestInfo{RequestId: test.uuid, ServingData: errdetails.Serialize(test.req)},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + test.resourceType,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, name, ip, addrs),
				},
			}, test.extra...)
			want := test.wrap(test.msg, test.cause, details...)
			if got, wantSig := statusSig(t, err), statusSig(t, want); got != wantSig {
				t.Errorf("error status changed:\n got: %s\nwant: %s", got, wantSig)
			}
			// The details are built exactly once, on the error path.
			if n := gw.addrsCalls.Load(); n != 1 {
				t.Errorf("gateway.Addrs called %d times, want 1", n)
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

// probeSearchClient answers every search kind with one result.
type probeSearchClient struct {
	*regClient
}

func (*probeSearchClient) SearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

func (*probeSearchClient) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

func (*probeSearchClient) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

// Test_server_Multi_boundedFanOut pins that the Multi* search handlers run at
// most multiConcurrency items at once (every gateway call of an item happens
// inside that item) and still answer every item in request order.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	t.Parallel()
	const (
		n     = 6
		limit = 2
	)
	cfg := func(i int) *payload.Search_Config {
		return &payload.Search_Config{RequestId: "req-" + strconv.Itoa(i), Num: 1, Timeout: int64(time.Minute)}
	}
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: []float32{1, 2, 3}, Config: cfg(i)})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: "id-" + strconv.Itoa(i), Config: cfg(i)})
		}
		return reqs
	}
	tests := []struct {
		call func(context.Context, *server) (*payload.Search_Responses, error)
		name string
	}{
		{
			name: "MultiSearch",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiSearch(ctx, searchReqs())
			},
		},
		{
			name: "MultiSearchByID",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiSearchByID(ctx, idReqs())
			},
		},
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiLinearSearch(ctx, searchReqs())
			},
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiLinearSearchByID(ctx, idReqs())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			p := newFanOutProbe(limit)
			var c vald.Client = &probeSearchClient{regClient: &regClient{
				search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					return probeResponse(), nil
				},
				getObject: func(_ context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
					return &payload.Object_Vector{Id: in.GetId().GetId(), Vector: []float32{1, 2, 3}}, nil
				},
			}}
			gw := syncGateway(c)
			gw.broadCast = func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
				p.hold()
				return f(ctx, "127.0.0.1:8081", c)
			}
			// eg is the shared group the handlers used to fan out on; GetObject still uses it.
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, replica: 1, multiConcurrency: limit, eg: eg, timeout: time.Minute}
			type result struct {
				res *payload.Search_Responses
				err error
			}
			done := make(chan result, 1)
			go func() {
				res, err := test.call(t.Context(), s)
				done <- result{res: res, err: err}
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
			if got.err != nil {
				t.Fatalf("unexpected error: %v", got.err)
			}
			if len(got.res.GetResponses()) != n {
				t.Fatalf("got %d responses, want %d", len(got.res.GetResponses()), n)
			}
			for i, r := range got.res.GetResponses() {
				if want := "req-" + strconv.Itoa(i); r.GetRequestId() != want {
					t.Errorf("responses[%d].RequestId = %q, want %q", i, r.GetRequestId(), want)
				}
			}
		})
	}
}

func Test_server_multiLimit(t *testing.T) {
	t.Parallel()
	cpuCap := runtime.GOMAXPROCS(0) * 10
	tests := []struct {
		name      string
		explicit  int
		replica   int
		agents    int
		fanOutAll bool
		want      int
	}{
		{name: "explicit value wins", explicit: 7, replica: 3, agents: 100, want: 7},
		{name: "write: 6 agents, replica 3 -> ceil(6*64/3)", replica: 3, agents: 6, want: min(cpuCap, 128)},
		{name: "write: 1 agent, replica 3 -> ceil(64/3)", replica: 3, agents: 1, want: min(cpuCap, 22)},
		{name: "write: 300 agents, replica 2 -> ceil(300*64/2)", replica: 2, agents: 300, want: min(cpuCap, 9600)},
		{name: "write: unknown agent count falls back to per-agent target", replica: 3, want: min(cpuCap, autoMultiPerAgent)},
		{name: "search reaches every agent -> per-agent target", replica: 3, agents: 300, fanOutAll: true, want: min(cpuCap, autoMultiPerAgent)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := &server{gateway: &regGateway{agents: tc.agents}, replica: tc.replica, multiConcurrency: tc.explicit}
			if got := s.multiLimit(t.Context(), tc.fanOutAll); got != tc.want {
				t.Errorf("multiLimit() = %d, want %d", got, tc.want)
			}
		})
	}
}

// The invalid-argument details of a Multi* batch (the serialized batch and the
// agent addresses) are the same for every invalid request, so they are built
// once per batch rather than once per invalid request.
func Test_server_Multi_buildsInvalidDetailsOncePerBatch(t *testing.T) {
	t.Parallel()
	const invalid = 3
	tests := []struct {
		call func(ctx context.Context, s *server) error
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Insert_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiInsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Update_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiUpdate(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Upsert_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiUpsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Remove_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{}})
				}
				_, err := s.MultiRemove(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{addrs: []string{"127.0.0.1:8081"}}
			s := &server{gateway: gw, name: "lb", ip: "192.0.2.1", multiConcurrency: 1}
			if err := test.call(t.Context(), s); err == nil {
				t.Fatal("expected the invalid requests to be reported")
			}
			if n := gw.addrsCalls.Load(); n != 1 {
				t.Errorf("gateway.Addrs called %d times for %d invalid requests, want 1", n, invalid)
			}
		})
	}
}

// pwCall is one agent call made by a test gateway; agents may be called several
// times under the same target, as the client backoff does.
type pwCall struct {
	agent  *pwAgent
	target string
}

// Partial writes used to be only logged: the failed agents and the partially
// applied requests are now counted, each failed agent once per request with the
// code of its latest failure.
//
//nolint:maintidx // table volume is inherent to this table-driven test
func Test_server_partialWrite_metrics(t *testing.T) {
	t.Parallel()
	internal := status.Error(codes.Internal, "agent internal")
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	removeForUpdate := vald.RemoveRPCName + " for " + vald.UpdateRPCName
	failure := func(rpc string, code codes.Code) attribute.Distinct {
		s := attribute.NewSet(attribute.String("rpc", rpc), attribute.String("grpc_code", code.String()))
		return s.Equivalent()
	}
	partial := func(rpc string) attribute.Distinct {
		s := attribute.NewSet(attribute.String("rpc", rpc))
		return s.Equivalent()
	}
	gateway := func(calls ...pwCall) *regGateway {
		run := func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var errs []error
			for _, c := range calls {
				if err := f(ctx, c.target, c.agent); err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}
		targets := make(map[string]struct{}, len(calls))
		for _, c := range calls {
			targets[c.target] = struct{}{}
		}
		return &regGateway{
			agents:    len(targets),
			broadCast: run,
			doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
				// the failures are only logged by DoMulti, as OrderedRange does.
				if err := run(ctx, f); err != nil {
					t.Logf("doMulti agent failures: %v", err)
				}
				return nil
			},
		}
	}
	update := func(cfg *payload.Update_Config) func(context.Context, *server) (*payload.Object_Location, error) {
		return func(ctx context.Context, s *server) (*payload.Object_Location, error) {
			return s.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
				Config: cfg,
			})
		}
	}
	insert := func(ctx context.Context, s *server) (*payload.Object_Location, error) {
		return s.Insert(ctx, &payload.Insert_Request{
			Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
			Config: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 1},
		})
	}
	updateTimestamp := func(ctx context.Context, s *server) (*payload.Object_Location, error) {
		return s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{Id: fixtureUUID, Timestamp: 10})
	}
	tests := []struct {
		call    func(context.Context, *server) (*payload.Object_Location, error)
		calls   func() []pwCall
		want    map[string]map[attribute.Distinct]int64
		name    string
		replica int
		wantErr bool
		// noCall makes DoMulti skip every agent without calling f, as it does for
		// agents whose circuit breaker is open.
		noCall bool
	}{
		{
			name: "Insert: one agent fails",
			call: insert,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, internal, nil)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
					{target: "c", agent: newPWAgent("c", nil, nil, nil)},
				}
			},
			replica: 2,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.WriteAgentFailureMetricsName:   {failure(vald.InsertRPCName, codes.Internal): 1},
				lbmetrics.WritePartialSuccessMetricsName: {partial(vald.InsertRPCName): 1},
			},
		},
		{
			name: "Insert: retried agent counted once with its latest code, recovered agent not counted",
			call: insert,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, internal, nil)},
					{target: "a", agent: newPWAgent("a", nil, context.Canceled, nil)},
					{target: "b", agent: newPWAgent("b", nil, unavailable, nil)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
					{target: "c", agent: newPWAgent("c", nil, nil, nil)},
				}
			},
			replica: 2,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.WriteAgentFailureMetricsName:   {failure(vald.InsertRPCName, codes.Canceled): 1},
				lbmetrics.WritePartialSuccessMetricsName: {partial(vald.InsertRPCName): 1},
			},
		},
		{
			name: "Insert: every agent fails is not a partial write",
			call: insert,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, internal, nil)},
					{target: "b", agent: newPWAgent("b", nil, internal, nil)},
				}
			},
			replica: 2,
			wantErr: true,
			want:    map[string]map[attribute.Distinct]int64{},
		},
		{
			name: "Insert: no agent called is a failure, not a partial write",
			call: insert,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, nil, nil)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
				}
			},
			noCall:  true,
			replica: 2,
			wantErr: true,
			want:    map[string]map[attribute.Distinct]int64{},
		},
		{
			name: "Insert: no failure records nothing",
			call: insert,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, nil, nil)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
				}
			},
			replica: 2,
			want:    map[string]map[attribute.Distinct]int64{},
		},
		{
			name: "Update with DisableBalancedUpdate: retried agents counted once",
			call: update(&payload.Update_Config{DisableBalancedUpdate: true, Timestamp: 10}),
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, nil, unavailable)},
					{target: "a", agent: newPWAgent("a", nil, nil, internal)},
					{target: "b", agent: newPWAgent("b", nil, nil, unavailable)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
					{target: "c", agent: newPWAgent("c", nil, nil, nil)},
				}
			},
			replica: 1,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.WriteAgentFailureMetricsName:   {failure(vald.UpdateRPCName, codes.Internal): 1},
				lbmetrics.WritePartialSuccessMetricsName: {partial(vald.UpdateRPCName): 1},
			},
		},
		{
			name: "Update with DisableBalancedUpdate: every agent fails is not a partial write",
			call: update(&payload.Update_Config{DisableBalancedUpdate: true, Timestamp: 10}),
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, nil, unavailable)},
					{target: "b", agent: newPWAgent("b", nil, nil, internal)},
				}
			},
			replica: 1,
			wantErr: true,
			want:    map[string]map[attribute.Distinct]int64{},
		},
		{
			name: "UpdateTimestamp: one agent fails",
			call: updateTimestamp,
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", nil, nil, unavailable)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
				}
			},
			replica: 1,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.WriteAgentFailureMetricsName:   {failure(vald.UpdateTimestampRPCName, codes.Unavailable): 1},
				lbmetrics.WritePartialSuccessMetricsName: {partial(vald.UpdateTimestampRPCName): 1},
			},
		},
		{
			name: "Update: remove fails on one agent, insert succeeds everywhere",
			call: update(&payload.Update_Config{SkipStrictExistCheck: true, Timestamp: 10}),
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", internal, nil, nil)},
					{target: "b", agent: newPWAgent("b", nil, nil, nil)},
				}
			},
			replica: 2,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.WriteAgentFailureMetricsName:   {failure(removeForUpdate, codes.Internal): 1},
				lbmetrics.WritePartialSuccessMetricsName: {partial(removeForUpdate): 1},
			},
		},
		{
			name: "Update: remove fails on every agent is not a partial remove",
			call: update(&payload.Update_Config{SkipStrictExistCheck: true, Timestamp: 10}),
			calls: func() []pwCall {
				return []pwCall{
					{target: "a", agent: newPWAgent("a", internal, nil, nil)},
					{target: "b", agent: newPWAgent("b", unavailable, nil, nil)},
				}
			},
			replica: 2,
			want:    map[string]map[attribute.Distinct]int64{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			sm := lbmetrics.New()
			gw := gateway(test.calls()...)
			if test.noCall {
				gw.doMulti = func(context.Context, int, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return nil
				}
			}
			s := &server{gateway: gw, replica: test.replica, searchMetrics: sm}
			loc, err := test.call(t.Context(), s)
			if test.wantErr {
				if err == nil {
					t.Fatalf("returned %v, want an error", loc)
				}
			} else if err != nil {
				t.Fatalf("returned %v, want success", err)
			}
			got := collectSearchMetrics(t, sm)
			if len(got) != len(test.want) {
				t.Errorf("got %d instruments with data, want %d: %v", len(got), len(test.want), got)
			}
			for name, wpts := range test.want {
				gpts := got[name]
				if len(gpts) != len(wpts) {
					t.Errorf("%s: got %d data points, want %d: %v", name, len(gpts), len(wpts), gpts)
				}
				for k, w := range wpts {
					if g, ok := gpts[k]; !ok || g != w {
						t.Errorf("%s: value for %v = %d (present %v), want %d", name, k, g, ok, w)
					}
				}
			}
		})
	}
}

// The invalid element makes MultiInsert/MultiRemove serialize the whole MultiRequest while
// the sibling Insert/Remove goroutines are running; those must not write to their requests
// (meaningful under -race).
func Test_server_Multi_invalidElementDoesNotRaceWithSiblings(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	t.Run("MultiInsert", func(t *testing.T) {
		t.Parallel()
		c := newOKClient(nil)
		s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: 4}
		reqs := &payload.Insert_MultiRequest{Requests: []*payload.Insert_Request{
			{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true},
			},
			{Vector: &payload.Object_Vector{Id: "invalid"}},
		}}
		if _, err := s.MultiInsert(t.Context(), reqs); err == nil {
			t.Error("expected the invalid element to be reported")
		}
		if ts := reqs.GetRequests()[0].GetConfig().GetTimestamp(); ts != 0 {
			t.Errorf("MultiInsert wrote a timestamp (%d) into the caller's request", ts)
		}
	})
	t.Run("MultiRemove", func(t *testing.T) {
		t.Parallel()
		c := newOKClient(nil)
		s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: 4}
		reqs := &payload.Remove_MultiRequest{Requests: []*payload.Remove_Request{
			{
				Id:     &payload.Object_ID{Id: uuid},
				Config: &payload.Remove_Config{SkipStrictExistCheck: true},
			},
			{Id: &payload.Object_ID{}},
		}}
		if _, err := s.MultiRemove(t.Context(), reqs); err == nil {
			t.Error("expected the invalid element to be reported")
		}
		if ts := reqs.GetRequests()[0].GetConfig().GetTimestamp(); ts != 0 {
			t.Errorf("MultiRemove wrote a timestamp (%d) into the caller's request", ts)
		}
	})
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want vald.Server
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, vald.Server) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got vald.Server) error {
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
