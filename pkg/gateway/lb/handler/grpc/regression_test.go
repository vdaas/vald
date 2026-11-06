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
	"strings"
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
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Fixture values shared by the handler tests of this package.
const (
	fixtureUUID = "uuid"
	loopbackIP  = "127.0.0.1"
)

// regGateway is a minimal service.Gateway stub. Unset hooks fall through to the
// embedded nil interface, so a test touching an unexpected method fails loudly.
type regGateway struct {
	service.Gateway
	broadCast func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error
	doMulti   func(ctx context.Context, num int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error
	addrs     []string
	agents    int
	// addrsCalls counts Addrs calls: the handlers query it only while building error details.
	addrsCalls atomic.Int64
	// agentCalls counts GetAgentCount calls.
	agentCalls atomic.Int64
}

func (g *regGateway) Addrs(context.Context) []string {
	g.addrsCalls.Add(1)
	return g.addrs
}

func (g *regGateway) GetAgentCount(context.Context) int {
	g.agentCalls.Add(1)
	return g.agents
}

func (g *regGateway) BroadCast(
	ctx context.Context,
	_ service.BroadCastKind,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error,
) error {
	return g.broadCast(ctx, f)
}

func (g *regGateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error,
) error {
	return g.doMulti(ctx, num, f)
}

// regClient is a minimal vald.Client stub.
type regClient struct {
	vald.Client
	update    func(ctx context.Context, in *payload.Update_Request) (*payload.Object_Location, error)
	search    func(ctx context.Context, in *payload.Search_Request) (*payload.Search_Response, error)
	insert    func(ctx context.Context, in *payload.Insert_Request) (*payload.Object_Location, error)
	remove    func(ctx context.Context, in *payload.Remove_Request) (*payload.Object_Location, error)
	exists    func(ctx context.Context, in *payload.Object_ID) (*payload.Object_ID, error)
	getObject func(ctx context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	getTS     func(ctx context.Context, in *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error)
}

func (c *regClient) Insert(
	ctx context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.insert(ctx, in)
}

func (c *regClient) Remove(
	ctx context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.remove(ctx, in)
}

func (c *regClient) Exists(
	ctx context.Context, in *payload.Object_ID, _ ...grpc.CallOption,
) (*payload.Object_ID, error) {
	return c.exists(ctx, in)
}

func (c *regClient) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, _ ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return c.getObject(ctx, in)
}

func (c *regClient) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Timestamp, error) {
	return c.getTS(ctx, in)
}

func (c *regClient) Update(
	ctx context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.update(ctx, in)
}

func (c *regClient) Search(
	ctx context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return c.search(ctx, in)
}

// Upsert with an empty uuid used to dereference a nil err while building the BadRequest detail.
func Test_server_Upsert_emptyUUID(t *testing.T) {
	t.Parallel()
	s := &server{gateway: new(regGateway), replica: 3}

	loc, err := s.Upsert(t.Context(), &payload.Upsert_Request{
		Vector: &payload.Object_Vector{Vector: []float32{1, 2, 3}},
	})
	if loc != nil {
		t.Errorf("unexpected location: %v", loc)
	}
	if err == nil {
		t.Fatal("expected an error for empty uuid, got nil")
	}
	st, ok := status.FromError(err)
	if !ok || st == nil || st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got: %v", err)
	}
}

// The replica top-up failure of Update (disable balanced update) used to be discarded.
func Test_server_Update_shortageDoMultiError(t *testing.T) {
	t.Parallel()
	errDoMulti := errors.New("doMulti failed")
	gw := &regGateway{
		agents: 3,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", &regClient{
				update: func(context.Context, *payload.Update_Request) (*payload.Object_Location, error) {
					return &payload.Object_Location{Name: "agent-0", Ips: []string{loopbackIP}}, nil
				},
			})
		},
		doMulti: func(context.Context, int, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return errDoMulti
		},
	}
	s := &server{gateway: gw, replica: 3}

	loc, err := s.Update(t.Context(), &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
		Config: &payload.Update_Config{DisableBalancedUpdate: true},
	})
	if !errors.Is(err, errDoMulti) {
		t.Errorf("got error: %v, want: %v", err, errDoMulti)
	}
	if loc != nil {
		t.Errorf("unexpected location on error: %v", loc)
	}
}

// Concurrent BroadCast callbacks must not write to the caller's request
// (run with -race); each agent gets its own request carrying the forwarding config.
func Test_server_Search_broadCastDoesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const agents = 4
	var (
		mu   sync.Mutex
		seen []*payload.Search_Request
	)
	gw := &regGateway{
		agents: agents,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var wg sync.WaitGroup
			errs := make([]error, agents)
			for i := range agents {
				wg.Go(func() {
					errs[i] = f(ctx, "127.0.0.1:8081", &regClient{
						search: func(_ context.Context, in *payload.Search_Request) (*payload.Search_Response, error) {
							mu.Lock()
							seen = append(seen, in)
							mu.Unlock()
							return &payload.Search_Response{
								Results: []*payload.Object_Distance{{Id: "id", Distance: 0.1}},
							}, nil
						},
					})
				})
			}
			wg.Wait()
			return errors.Join(errs...)
		},
	}
	s := &server{gateway: gw, replica: 1, timeout: time.Second}

	cfg := &payload.Search_Config{Num: 3}
	req := &payload.Search_Request{Vector: []float32{1, 2, 3}, Config: cfg}
	if _, err := s.Search(t.Context(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.GetConfig() != cfg {
		t.Error("Search replaced the caller's req.Config")
	}
	if len(seen) != agents {
		t.Fatalf("got %d agent requests, want %d", len(seen), agents)
	}
	for _, in := range seen {
		if in == req {
			t.Error("the caller's request object was forwarded to an agent")
		}
		if in.GetConfig() == cfg {
			t.Error("the caller's config was forwarded to an agent instead of the forwarding config")
		}
	}
}

// okClient answers every write/read the Insert/Remove/Update handlers issue, recording the
// requests it was forwarded.
type okClient struct {
	*regClient
	insertIn *payload.Insert_Request
	removeIn *payload.Remove_Request
}

func newOKClient(vec []float32) *okClient {
	const uuid = "uuid"
	c := new(okClient)
	loc := &payload.Object_Location{Uuid: uuid, Name: "agent-0", Ips: []string{loopbackIP}}
	c.regClient = &regClient{
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			c.insertIn = in
			return loc, nil
		},
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			c.removeIn = in
			return loc, nil
		},
		exists: func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: uuid}, nil
		},
		getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{Id: uuid, Vector: vec}, nil
		},
	}
	return c
}

// syncGateway runs every callback once, inline, against c.
func syncGateway(c vald.Client) *regGateway {
	return &regGateway{
		agents: 1,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", c)
		},
		doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", c)
		},
	}
}

// Insert used to write the exist-check flag and the timestamp into the caller's request,
// which MultiInsert shares with its other readers. It must forward a new request instead.
func Test_server_Insert_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	filters := &payload.Filter_Config{Targets: []*payload.Filter_Target{{Host: "filter", Port: 1}}}
	tests := []struct {
		cfg        *payload.Insert_Config
		wantFilter *payload.Filter_Config
		name       string
		wantTS     int64 // 0: any non-zero timestamp
	}{
		{name: "nil config (exist check runs)"},
		{name: "skip check, no timestamp", cfg: &payload.Insert_Config{SkipStrictExistCheck: true, Filters: filters}, wantFilter: filters},
		{name: "exist check requested, filters kept", cfg: &payload.Insert_Config{Filters: filters}, wantFilter: filters},
		{name: "timestamp already set", cfg: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 7}, wantTS: 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newOKClient([]float32{4, 5, 6})
			// Exists finds nothing, so Insert may proceed.
			c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
				return &payload.Object_ID{}, nil
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			vector := &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}}
			req := &payload.Insert_Request{Vector: vector, Config: test.cfg}
			var before payload.Insert_Config
			if test.cfg != nil {
				before.SkipStrictExistCheck, before.Filters, before.Timestamp = test.cfg.SkipStrictExistCheck, test.cfg.Filters, test.cfg.Timestamp
			}
			if _, err := s.Insert(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Insert replaced the caller's req.Config")
			}
			if test.cfg != nil && (test.cfg.SkipStrictExistCheck != before.SkipStrictExistCheck ||
				test.cfg.Filters != before.Filters || test.cfg.Timestamp != before.Timestamp) {
				t.Errorf("Insert modified the caller's config: got %v", test.cfg)
			}
			got := c.insertIn
			if got == nil {
				t.Fatal("no request reached the agent")
			}
			if got.GetVector() != vector {
				t.Error("the forwarded request does not carry the caller's vector")
			}
			if !got.GetConfig().GetSkipStrictExistCheck() {
				t.Error("the forwarded request must skip the agent's exist check")
			}
			if ts := got.GetConfig().GetTimestamp(); ts == 0 || (test.wantTS != 0 && ts != test.wantTS) {
				t.Errorf("forwarded timestamp = %d, want %d (0: any non-zero)", ts, test.wantTS)
			}
			if got.GetConfig().GetFilters() != test.wantFilter {
				t.Errorf("forwarded filters = %v, want %v", got.GetConfig().GetFilters(), test.wantFilter)
			}
		})
	}
}

// Remove used to write the exist-check flag and the timestamp into the caller's request.
func Test_server_Remove_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		cfg    *payload.Remove_Config
		name   string
		wantTS int64 // 0: any non-zero timestamp
	}{
		{name: "nil config (exist check runs)"},
		{name: "skip check, no timestamp", cfg: &payload.Remove_Config{SkipStrictExistCheck: true}},
		{name: "exist check requested", cfg: &payload.Remove_Config{}},
		{name: "timestamp already set", cfg: &payload.Remove_Config{SkipStrictExistCheck: true, Timestamp: 7}, wantTS: 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newOKClient(nil)
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			id := &payload.Object_ID{Id: uuid}
			req := &payload.Remove_Request{Id: id, Config: test.cfg}
			var before payload.Remove_Config
			if test.cfg != nil {
				before.SkipStrictExistCheck, before.Timestamp = test.cfg.SkipStrictExistCheck, test.cfg.Timestamp
			}
			if _, err := s.Remove(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Remove replaced the caller's req.Config")
			}
			if test.cfg != nil && (test.cfg.SkipStrictExistCheck != before.SkipStrictExistCheck || test.cfg.Timestamp != before.Timestamp) {
				t.Errorf("Remove modified the caller's config: got %v", test.cfg)
			}
			got := c.removeIn
			if got == nil {
				t.Fatal("no request reached the agent")
			}
			if got.GetId() != id {
				t.Error("the forwarded request does not carry the caller's id")
			}
			if !got.GetConfig().GetSkipStrictExistCheck() {
				t.Error("the forwarded request must skip the agent's exist check")
			}
			if ts := got.GetConfig().GetTimestamp(); ts == 0 || (test.wantTS != 0 && ts != test.wantTS) {
				t.Errorf("forwarded timestamp = %d, want %d (0: any non-zero)", ts, test.wantTS)
			}
		})
	}
}

// Update used to set SkipStrictExistCheck on the caller's request after its own exist
// check, a write nothing read back (it builds fresh Remove/Insert requests).
func Test_server_Update_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		cfg  *payload.Update_Config
		name string
	}{
		{name: "nil config"},
		{name: "exist check requested", cfg: &payload.Update_Config{Timestamp: 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one, so Update proceeds to remove+insert.
			c := newOKClient([]float32{4, 5, 6})
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			req := &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: test.cfg,
			}
			if _, err := s.Update(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Update replaced the caller's req.Config")
			}
			if test.cfg != nil && test.cfg.SkipStrictExistCheck {
				t.Error("Update set SkipStrictExistCheck on the caller's config")
			}
			if c.removeIn == nil || c.insertIn == nil {
				t.Fatalf("remove/insert did not reach the agent: remove=%v insert=%v", c.removeIn, c.insertIn)
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

// exists/getObject used to return a spurious NotFound when the select picked the
// nil broadcast result from ech while the found object was still buffered in the
// result channel. The race is random, so repeat the call and require every one to succeed.
func Test_server_existsAndGetObject_foundObjectWinsOverNilBroadCast(t *testing.T) {
	t.Parallel()
	const (
		uuid  = "uuid"
		calls = 200
	)
	eg, _ := errgroup.New(t.Context())
	// limit 1 runs eg.Go inline, so the found result, the canceled ctx and the broadcast result on ech
	// are all ready before the select runs; the old code then picks the ech case at random and fails.
	eg.SetLimit(1)
	s := &server{gateway: syncGateway(newOKClient([]float32{1, 2, 3})), replica: 1, eg: eg}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			id, err := s.exists(ctx, uuid)
			cancel()
			if err != nil || id.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, id, err, uuid)
			}
		}
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			vec, err := s.getObject(ctx, uuid)
			cancel()
			if err != nil || vec.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, vec, err, uuid)
			}
		}
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		c := newOKClient([]float32{1, 2, 3})
		c.getTS = func(context.Context, *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
			return &payload.Object_Timestamp{Id: uuid, Timestamp: 1}, nil
		}
		eg, _ := errgroup.New(t.Context())
		eg.SetLimit(1)
		ts := &server{gateway: syncGateway(c), replica: 1, eg: eg}
		req := &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}}
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			got, err := ts.GetTimestamp(ctx, req)
			cancel()
			if err != nil || got.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, got, err, uuid)
			}
		}
	})
}

// When nothing was found, exists/getObject/GetTimestamp used to pick the closed result channel
// or the broadcast error at random, turning a connection error into NotFound about half the time.
func Test_server_existsAndGetObject_broadCastErrorWinsOverClosedResult(t *testing.T) {
	t.Parallel()
	const (
		uuid  = "uuid"
		calls = 200
	)
	newServer := func() *server {
		eg, _ := errgroup.New(t.Context())
		// limit 1 runs eg.Go inline, so the closed result channel and the buffered broadcast
		// error are both ready before the select runs.
		eg.SetLimit(1)
		return &server{
			gateway: &regGateway{
				broadCast: func(context.Context, func(context.Context, string, vald.Client, ...grpc.CallOption) error) error {
					return errors.ErrGRPCClientConnNotFound("*")
				},
			},
			replica: 1,
			eg:      eg,
		}
	}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		for i := range calls {
			if _, err := s.exists(t.Context(), uuid); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				t.Fatalf("call %d: got %v, want connection not found", i, err)
			}
		}
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		for i := range calls {
			if _, err := s.getObject(t.Context(), uuid); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				t.Fatalf("call %d: got %v, want connection not found", i, err)
			}
		}
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		req := &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}}
		for i := range calls {
			_, err := s.GetTimestamp(t.Context(), req)
			if st, _ := status.FromError(err); st.Code() != codes.Internal {
				t.Fatalf("call %d: got %v, want code Internal", i, err)
			}
		}
	})
}

// doSearch used to query the agent count twice when a ratio was set (once itself, once in calculateNum).
func Test_server_doSearch_queriesAgentCountOnce(t *testing.T) {
	t.Parallel()
	gw := &regGateway{
		agents: 3,
		broadCast: func(context.Context, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return nil
		},
	}
	s := &server{gateway: gw, replica: 1, timeout: time.Second}

	// The broadcast finds nothing, so the search itself must fail with an empty-result error.
	if _, _, err := s.doSearch(t.Context(), vald.SearchRPCName, &payload.Search_Config{
		Num:   3,
		Ratio: wrapperspb.Float(0.5),
	}, func(context.Context, *payload.Search_Config, vald.Client, ...grpc.CallOption) (*payload.Search_Response, error) {
		return new(payload.Search_Response), nil
	}); err == nil {
		t.Error("expected an error for an empty search result")
	}
	if got := gw.agentCalls.Load(); got != 1 {
		t.Errorf("GetAgentCount called %d times, want 1", got)
	}
}

// A configured concurrency of 1 is valid and must not be dropped in favor of the default.
func Test_WithConcurrencyOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		c    int
		want int
	}{
		{name: "one is accepted", c: 1, want: 1},
		{name: "positive is accepted", c: 8, want: 8},
		{name: "zero keeps the current value", c: 0, want: 5},
		{name: "negative keeps the current value", c: -1, want: 5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s := &server{streamConcurrency: 5, multiConcurrency: 5}
			WithStreamConcurrency(test.c)(s)
			WithMultiConcurrency(test.c)(s)
			if s.streamConcurrency != test.want {
				t.Errorf("streamConcurrency = %d, want %d", s.streamConcurrency, test.want)
			}
			if s.multiConcurrency != test.want {
				t.Errorf("multiConcurrency = %d, want %d", s.multiConcurrency, test.want)
			}
		})
	}
}

// The ResourceType detail of an index broadcast error names the "BroadCast" stage (was "BroadCase").
func Test_server_handleIndexBroadCastError_resourceType(t *testing.T) {
	t.Parallel()
	const target = "127.0.0.1:8081"
	s := &server{name: "gateway", ip: loopbackIP}

	err := s.handleIndexBroadCastError(nil, target, vald.IndexInfoRPCName, status.Error(codes.Internal, "boom"))
	if err == nil {
		t.Fatal("expected the Internal error to be propagated")
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("expected a gRPC status, got: %v", err)
	}
	for _, d := range st.Details() {
		ri, ok := d.(*errdetails.ResourceInfo)
		if !ok {
			continue
		}
		want := "." + "BroadCast/" + target
		if !strings.HasSuffix(ri.GetResourceType(), want) {
			t.Errorf("ResourceType = %q, want suffix %q", ri.GetResourceType(), want)
		}
		return
	}
	t.Errorf("no ResourceInfo detail in %v", st.Details())
}

// A NotFound answer from an agent is tolerated by Remove and ParseError never attaches details to
// it, so the (costly) error details must not be built for it: Addrs is queried only while
// building them. A real agent failure must still carry them.
func Test_server_Remove_buildsErrorDetailsOnlyForReturnedErrors(t *testing.T) {
	t.Parallel()
	const agents = 4
	tests := []struct {
		agentErr      error
		name          string
		wantAddrsCall bool
		wantErr       bool
	}{
		{name: "NotFound agents are tolerated without details", agentErr: status.Error(codes.NotFound, "not found")},
		{name: "Internal agent error carries details", agentErr: status.Error(codes.Internal, "boom"), wantAddrsCall: true, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ok := newOKClient(nil)
			gw := &regGateway{
				agents: agents,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					errs := make([]error, 0, agents)
					for i := range agents {
						if i == agents-1 {
							// the last agent holds the object
							errs = append(errs, f(ctx, "127.0.0.1:8081", ok))
							continue
						}
						errs = append(errs, f(ctx, "127.0.0.1:8081", &regClient{
							remove: func(context.Context, *payload.Remove_Request) (*payload.Object_Location, error) {
								return nil, test.agentErr
							},
						}))
					}
					return errors.Join(errs...)
				},
			}
			s := &server{gateway: gw, replica: agents}

			_, err := s.Remove(t.Context(), &payload.Remove_Request{
				Id:     &payload.Object_ID{Id: fixtureUUID},
				Config: &payload.Remove_Config{SkipStrictExistCheck: true},
			})
			if (err != nil) != test.wantErr {
				t.Fatalf("got error: %v, want error: %t", err, test.wantErr)
			}
			if got := gw.addrsCalls.Load(); (got != 0) != test.wantAddrsCall {
				t.Errorf("Addrs called %d times, want called: %t", got, test.wantAddrsCall)
			}
		})
	}
}

// When every entry is invalid there is nothing to merge; this used to index into an empty slice.
func Test_mergeInfoIndexStatistics_noValidEntry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		stats     map[string]*payload.Info_Index_Statistics
		name      string
		wantValid bool
	}{
		{name: "nil", stats: nil},
		{name: "all invalid", stats: map[string]*payload.Info_Index_Statistics{
			"a": {Valid: false, NumberOfNodes: 3},
			"b": {Valid: false},
		}},
		{name: "one valid one invalid", stats: map[string]*payload.Info_Index_Statistics{
			"a": {Valid: true, NumberOfNodes: 3},
			"b": {Valid: false},
		}, wantValid: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := mergeInfoIndexStatistics(test.stats)
			if got == nil {
				t.Fatal("got nil statistics")
			}
			if got.GetValid() != test.wantValid {
				t.Errorf("Valid = %t, want %t", got.GetValid(), test.wantValid)
			}
			if !test.wantValid && got.GetNumberOfNodes() != 0 {
				t.Errorf("NumberOfNodes = %d, want 0 for nothing to merge", got.GetNumberOfNodes())
			}
		})
	}
}

// propClient answers IndexProperty with a canned result.
type propClient struct {
	*regClient
	err error
}

func (c propClient) IndexProperty(
	context.Context, *payload.Empty, ...grpc.CallOption,
) (*payload.Info_Index_PropertyDetail, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &payload.Info_Index_PropertyDetail{
		Details: map[string]*payload.Info_Index_Property{"p": {Dimension: 3}},
	}, nil
}

// IndexProperty shares handleIndexBroadCastError with the other Index* handlers: a failing agent is
// propagated, a NotFound/canceled one is tolerated, a healthy one is merged under "<target>-<key>".
func Test_server_IndexProperty_broadCastErrorHandling(t *testing.T) {
	t.Parallel()
	const target = "127.0.0.1:8081"
	tests := []struct {
		agentErr error
		name     string
		wantCode codes.Code // codes.OK: no error expected
		wantKeys int
	}{
		{name: "healthy agent", wantKeys: 1},
		{name: "NotFound is tolerated", agentErr: status.Error(codes.NotFound, "nf")},
		{name: "Canceled is tolerated", agentErr: context.Canceled},
		{name: "Internal is propagated", agentErr: status.Error(codes.Internal, "boom"), wantCode: codes.Internal},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{
				agents: 1,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return f(ctx, target, propClient{regClient: new(regClient), err: test.agentErr})
				},
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, eg: eg, name: "gateway", ip: loopbackIP}

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			detail, err := s.IndexProperty(ctx, new(payload.Empty))
			if test.wantCode != codes.OK {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := len(detail.GetDetails()); got != test.wantKeys {
				t.Errorf("got %d property entries, want %d", got, test.wantKeys)
			}
			if test.wantKeys > 0 && detail.GetDetails()[target+"-p"] == nil {
				t.Errorf("missing %q in %v", target+"-p", detail.GetDetails())
			}
		})
	}
}

// exists/getObject/GetTimestamp derive a cancelable context that was only canceled once an object was found;
// every other return path left it registered with the parent. The ctx handed to BroadCast must be
// canceled by the time they return.
func Test_server_existsAndGetObject_cancelDerivedContext(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	run := func(t *testing.T, call func(s *server, ctx context.Context) error) {
		t.Helper()
		var bctx context.Context
		c := newOKClient(nil)
		// nothing is found, so the call returns through the broadcast-result path
		c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{}, nil
		}
		c.getObject = func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{}, nil
		}
		c.getTS = func(context.Context, *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
			return &payload.Object_Timestamp{}, nil
		}
		gw := &regGateway{
			agents: 1,
			broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
				bctx = ctx
				return f(ctx, "127.0.0.1:8081", c)
			},
		}
		eg, _ := errgroup.New(t.Context())
		s := &server{gateway: gw, replica: 1, eg: eg}

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		if err := call(s, ctx); err == nil {
			t.Fatal("expected a not found error")
		}
		if bctx == nil {
			t.Fatal("BroadCast was not called")
		}
		if bctx.Err() == nil {
			t.Error("the derived context was not canceled when the call returned")
		}
	}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.exists(ctx, uuid)
			return err
		})
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.getObject(ctx, uuid)
			return err
		})
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.GetTimestamp(ctx, &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}})
			return err
		})
	})
}

// calculateNum used to return n-1 with n == 0 (num == 0 with a ratio), wrapping to MaxUint32.
func Test_server_calculateNum_zeroN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		num    uint32
		ratio  float32
		agents int
		want   uint32
	}{
		{name: "num 0 with ratio", num: 0, ratio: 0.5, agents: 3, want: 0},
		{name: "regular ratio unchanged", num: 10, ratio: 0.5, agents: 2, want: 7},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := &server{replica: 1}
			if got := s.calculateNum(tc.num, tc.ratio, tc.agents); got != tc.want {
				t.Errorf("calculateNum(%d, %v, %d) = %d, want %d", tc.num, tc.ratio, tc.agents, got, tc.want)
			}
		})
	}
}
