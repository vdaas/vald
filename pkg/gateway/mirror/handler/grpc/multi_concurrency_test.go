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
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
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

// probeMirrorClient answers every write with a location echoing the uuid; the
// uuid "id-0" panics instead when panics is set. Other methods panic through
// the nil embedded interface.
type probeMirrorClient struct {
	service.MirrorClient
	panics bool
}

func (c probeMirrorClient) location(uuid string) (*payload.Object_Location, error) {
	if c.panics && uuid == "id-0" {
		panic("probe panic")
	}
	return &payload.Object_Location{Uuid: uuid}, nil
}

func (c probeMirrorClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Upsert(
	_ context.Context, in *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetId().GetId())
}

// Test_server_Multi_boundedFanOut pins that the Multi* write handlers run at
// most streamConcurrency items at once, still answer every item in request
// order, and report a panicking item as an error instead of dropping it.
//
// It deliberately does not call t.Parallel: its subtests hold handler
// goroutines open, which the goleak.VerifyNone checks of the parallel
// handler_test.go cases would otherwise report as leaks.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	const (
		n     = 6
		limit = 2
	)
	id := func(i int) string { return "id-" + strconv.Itoa(i) }
	vec := func(i int) *payload.Object_Vector {
		return &payload.Object_Vector{Id: id(i), Vector: []float32{1, 2, 3}}
	}
	tests := []struct {
		call func(context.Context, *server) (*payload.Object_Locations, error)
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Insert_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: vec(i)})
				}
				return s.MultiInsert(ctx, reqs)
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Update_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: vec(i)})
				}
				return s.MultiUpdate(ctx, reqs)
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Upsert_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: vec(i)})
				}
				return s.MultiUpsert(ctx, reqs)
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Remove_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id(i)}})
				}
				return s.MultiRemove(ctx, reqs)
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
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(context.Context) string {
							return proxiedFrom // proxied request: forward to the own LB only
						},
						DoFunc: func(ctx context.Context, target string, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
							p.hold()
							return f(ctx, target, probeMirrorClient{panics: panics})
						},
					},
					streamConcurrency: limit,
				}
				type result struct {
					res *payload.Object_Locations
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
				if panics {
					// The handlers may re-wrap the recovered panic into a status,
					// so only its presence is pinned here.
					if got.err == nil {
						t.Error("err = nil, want the recovered panic")
					}
					// The panicking item 0 must not make the group skip the queued
					// items: every item ran and every other item has its location.
					if _, calls := p.stats(); calls != n {
						t.Errorf("forwarded %d items, want all %d", calls, n)
					}
					locs := got.res.GetLocations()
					if len(locs) != n {
						t.Fatalf("got %d locations, want %d", len(locs), n)
					}
					for i := 1; i < n; i++ {
						if want := id(i); locs[i].GetUuid() != want {
							t.Errorf("locations[%d].Uuid = %q, want %q", i, locs[i].GetUuid(), want)
						}
					}
					return
				}
				if got.err != nil {
					t.Fatalf("unexpected error: %v", got.err)
				}
				if len(got.res.GetLocations()) != n {
					t.Fatalf("got %d locations, want %d", len(got.res.GetLocations()), n)
				}
				for i, l := range got.res.GetLocations() {
					if want := id(i); l.GetUuid() != want {
						t.Errorf("locations[%d].Uuid = %q, want %q", i, l.GetUuid(), want)
					}
				}
			})
		}
	}
}
