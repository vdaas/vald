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

// probeNGT answers every search through the probe; the vector whose first
// element is 0 and the id "id-0" panic instead when panics is set.
type probeNGT struct {
	p *fanOutProbe
	stubNGT
	panics bool
}

func (n *probeNGT) search(first float32) (*payload.Search_Response, error) {
	n.p.hold()
	if n.panics && first == 0 {
		panic("probe panic")
	}
	return &payload.Search_Response{}, nil
}

func (n *probeNGT) searchByID(uuid string) ([]float32, *payload.Search_Response, error) {
	n.p.hold()
	if n.panics && uuid == "id-0" {
		panic("probe panic")
	}
	return nil, &payload.Search_Response{}, nil
}

func (n *probeNGT) Search(
	_ context.Context, vec []float32, _ uint32, _, _ float32, _ int32,
) (*payload.Search_Response, error) {
	return n.search(vec[0])
}

func (n *probeNGT) SearchByID(
	_ context.Context, uuid string, _ uint32, _, _ float32, _ int32,
) ([]float32, *payload.Search_Response, error) {
	return n.searchByID(uuid)
}

func (n *probeNGT) LinearSearch(
	_ context.Context, vec []float32, _ uint32,
) (*payload.Search_Response, error) {
	return n.search(vec[0])
}

func (n *probeNGT) LinearSearchByID(
	_ context.Context, uuid string, _ uint32,
) ([]float32, *payload.Search_Response, error) {
	return n.searchByID(uuid)
}

// Test_server_Multi_boundedFanOut pins that the Multi* search handlers run at
// most streamConcurrency items at once, still answer every item in request
// order, and report a panicking item as an error instead of dropping it.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	t.Parallel()
	const (
		n     = 6
		limit = 2
	)
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{
				Vector: []float32{float32(i), 1},
				Config: &payload.Search_Config{RequestId: "req-" + strconv.Itoa(i)},
			})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{
				Id:     "id-" + strconv.Itoa(i),
				Config: &payload.Search_Config{RequestId: "req-" + strconv.Itoa(i)},
			})
		}
		return reqs
	}
	calls := []struct {
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
	for _, c := range calls {
		for _, panics := range []bool{false, true} {
			t.Run(c.name+"/panics="+strconv.FormatBool(panics), func(t *testing.T) {
				t.Parallel()
				p := newFanOutProbe(limit)
				// eg is the shared group the handlers used to fan out on.
				eg, _ := errgroup.New(t.Context())
				s := &server{
					ngt:               &probeNGT{dim: 2, p: p, panics: panics},
					eg:                eg,
					streamConcurrency: limit,
				}
				type result struct {
					res *payload.Search_Responses
					err error
				}
				done := make(chan result, 1)
				go func() {
					res, err := c.call(t.Context(), s)
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
					if got.err == nil || !strings.Contains(got.err.Error(), "probe panic") {
						t.Errorf("err = %v, want the recovered panic", got.err)
					}
					// The panicking item must not make the group skip the queued ones.
					if _, calls := p.stats(); calls != n {
						t.Errorf("searched %d items, want all %d", calls, n)
					}
					return
				}
				if got.err != nil {
					t.Fatalf("unexpected error: %v", got.err)
				}
				if _, calls := p.stats(); calls != n {
					t.Errorf("searched %d items, want %d", calls, n)
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
}
