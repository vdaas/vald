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

// probeSearchClient answers every search kind with one result.
type probeSearchClient struct {
	*regClient
}

func probeResponse() *payload.Search_Response {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "found", Distance: 0.1}}}
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
