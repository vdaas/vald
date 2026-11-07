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
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

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
