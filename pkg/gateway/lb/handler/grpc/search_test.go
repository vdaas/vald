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
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/net/grpc/types"
	"github.com/vdaas/vald/internal/observability/attribute"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/sync"
)

func runTestSearch(
	ctx context.Context,
	s *server,
	rpc string,
	cfg *payload.Search_Config,
	f func(context.Context, *payload.Search_Config, vald.Client, ...grpc.CallOption) (*payload.Search_Response, error),
) (*payload.Search_Response, error) {
	replica, aggr, err := s.prepareSearch(ctx, nil, cfg)
	if err != nil {
		return nil, err
	}
	fcfg := forwardSearchConfig(cfg, aggr.GetFnum())
	res, _, err := s.aggregationSearch(ctx, rpc, replica, aggr, cfg,
		func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return f(ctx, fcfg, vc, copts...)
		},
	)
	return res, err
}

// An agent that rejects the search (wrong dimension: InvalidArgument; indexing or
// flushing: Aborted) used to be treated as a miss, so when no agent had results
// the client got NotFound instead of the agents' code. Results from the other
// agents are still served.
func Test_server_doSearch_agentErrorCodes(t *testing.T) {
	t.Parallel()
	agentA, agentB := new(regClient), new(regClient)
	hit := &payload.Search_Response{
		Results: []*payload.Object_Distance{{Id: "hit", Distance: 0.1}},
	}
	tests := []struct {
		errA     error
		errB     error
		resB     *payload.Search_Response
		name     string
		wantCode codes.Code // codes.OK: a result is expected
	}{
		{
			name:     "all agents InvalidArgument",
			errA:     status.Error(codes.InvalidArgument, "Search API Incompatible Dimension Size detected"),
			errB:     status.Error(codes.InvalidArgument, "Search API Incompatible Dimension Size detected"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "all agents Aborted",
			errA:     status.Error(codes.Aborted, "indexing"),
			errB:     status.Error(codes.Aborted, "flushing"),
			wantCode: codes.Aborted,
		},
		{
			name:     "all agents NotFound stays NotFound",
			errA:     status.Error(codes.NotFound, "nf"),
			errB:     status.Error(codes.NotFound, "nf"),
			wantCode: codes.NotFound,
		},
		{
			name:     "Aborted agent does not hide the other agent's results",
			errA:     status.Error(codes.Aborted, "indexing"),
			resB:     hit,
			wantCode: codes.OK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{
				agents: 2,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					// models service.Gateway.BroadCast: every callback error is returned.
					return errors.Join(f(ctx, "a:8081", agentA), f(ctx, "b:8081", agentB))
				},
			}
			s := &server{gateway: gw, replica: 1, timeout: 10 * time.Second}
			res, err := runTestSearch(t.Context(), s, vald.SearchRPCName, &payload.Search_Config{Num: 1},
				func(_ context.Context, _ *payload.Search_Config, vc vald.Client, _ ...grpc.CallOption) (*payload.Search_Response, error) {
					if vc == agentA {
						return nil, test.errA
					}
					return test.resB, test.errB
				})
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got := len(res.GetResults()); got != 1 {
					t.Fatalf("got %d results, want 1", got)
				}
				return
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != test.wantCode {
				t.Fatalf("got error: %v, want code %v", err, test.wantCode)
			}
		})
	}
}

func Test_skipSearchByIDFallback(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err      error
		name     string
		canceled bool
		want     bool
	}{
		{name: "DeadlineExceeded status skips the fallback", err: status.Error(codes.DeadlineExceeded, "timeout"), want: true},
		{name: "Canceled status skips the fallback", err: status.Error(codes.Canceled, "canceled"), want: true},
		{name: "raw context.DeadlineExceeded skips the fallback", err: context.DeadlineExceeded, want: true},
		{name: "raw context.Canceled skips the fallback", err: context.Canceled, want: true},
		{name: "done request ctx skips the fallback", canceled: true, err: status.Error(codes.Internal, "internal"), want: true},
		{name: "NotFound still falls back", err: status.Error(codes.NotFound, "not found")},
		{name: "Internal still falls back", err: status.Error(codes.Internal, "internal")},
		{name: "Unavailable still falls back", err: status.Error(codes.Unavailable, "unavailable")},
		{name: "non-status error still falls back", err: errors.New("plain")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			if tc.canceled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			if got := skipSearchByIDFallback(ctx, tc.err); got != tc.want {
				t.Errorf("skipSearchByIDFallback() = %v, want %v", got, tc.want)
			}
		})
	}
}

// A search that loses some agents still returns the other agents' results; the
// failed agents are counted once each by their last outcome, and the request is
// counted as a partial success labelled by whether the failures reach replica.
//
//nolint:maintidx // table volume is inherent to this table-driven test
func Test_server_doSearch_partialFailureMetrics(t *testing.T) {
	t.Parallel()
	failureKey := func(rpc string, code codes.Code) attribute.Distinct {
		s := attribute.NewSet(attribute.String("rpc", rpc), attribute.String("grpc_code", code.String()))
		return s.Equivalent()
	}
	partialKey := func(rpc string, risk bool) attribute.Distinct {
		s := attribute.NewSet(attribute.String("rpc", rpc), attribute.Bool("recall_risk", risk))
		return s.Equivalent()
	}
	unavailable := status.Error(codes.Unavailable, "agent down")
	internal := status.Error(codes.Internal, "agent broken")
	tests := []struct {
		// replies is each target's reply per call in order; nil means a hit.
		replies map[string][]error
		want    map[string]map[attribute.Distinct]int64
		name    string
		rpc     string
		// calls is the order the fake BroadCast invokes the callback per
		// target; a repeated target models a client backoff retry.
		calls    []string
		wantCode codes.Code // codes.OK: a result is expected
		// waitDeadline makes the fake BroadCast return only once the request
		// deadline has passed, modelling a broadcast cut short by it.
		waitDeadline bool
	}{
		{
			name:  "one failed agent below replica",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {unavailable},
				"b": {nil},
				"c": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.Unavailable): 1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.SearchRPCName, false): 1,
				},
			},
		},
		{
			name:  "replica agents failed: recall at risk",
			rpc:   vald.LinearSearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {unavailable},
				"b": {internal},
				"c": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.LinearSearchRPCName, codes.Unavailable): 1,
					failureKey(vald.LinearSearchRPCName, codes.Internal):    1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.LinearSearchRPCName, true): 1,
				},
			},
		},
		{
			name:  "all agents hit records nothing",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {nil},
				"b": {nil},
				"c": {nil},
			},
			wantCode: codes.OK,
			want:     map[string]map[attribute.Distinct]int64{},
		},
		{
			name:  "DeadlineExceeded and NotFound misses are not failures",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {context.DeadlineExceeded},
				"b": {status.Error(codes.NotFound, "nf")},
				"c": {nil},
			},
			wantCode: codes.OK,
			want:     map[string]map[attribute.Distinct]int64{},
		},
		{
			name:  "all agents failed: failures counted, no partial success",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {unavailable},
				"b": {unavailable},
				"c": {unavailable},
			},
			wantCode: codes.Unavailable,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.Unavailable): 3,
				},
			},
		},
		{
			name:  "all agents reject the request with InvalidArgument: no agent failure",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {status.Error(codes.InvalidArgument, "dim")},
				"b": {status.Error(codes.InvalidArgument, "dim")},
				"c": {status.Error(codes.InvalidArgument, "dim")},
			},
			wantCode: codes.InvalidArgument,
			want:     map[string]map[attribute.Distinct]int64{},
		},
		{
			name:  "InvalidArgument from one agent is not an agent failure",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {status.Error(codes.InvalidArgument, "dim")},
				"b": {nil},
				"c": {nil},
			},
			wantCode: codes.OK,
			want:     map[string]map[attribute.Distinct]int64{},
		},
		{
			name:  "other codes and non-status errors are failures",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b", "c"},
			replies: map[string][]error{
				"a": {status.Error(codes.FailedPrecondition, "fp")},
				"b": {errors.New("boom")},
				"c": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.FailedPrecondition): 1,
					failureKey(vald.SearchRPCName, codes.Unknown):            1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.SearchRPCName, true): 1,
				},
			},
		},
		{
			name:  "agent never invoked is counted as unreached",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b"},
			replies: map[string][]error{
				"a": {nil},
				"b": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.Unavailable): 1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.SearchRPCName, false): 1,
				},
			},
		},
		{
			name:  "unreached agent and failed agent reach replica",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b"},
			replies: map[string][]error{
				"a": {internal},
				"b": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.Internal):    1,
					failureKey(vald.SearchRPCName, codes.Unavailable): 1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.SearchRPCName, true): 1,
				},
			},
		},
		{
			name:  "agent not invoked before the deadline is not counted",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "b"},
			replies: map[string][]error{
				"a": {nil},
				"b": {nil},
			},
			waitDeadline: true,
			wantCode:     codes.OK,
			want:         map[string]map[attribute.Distinct]int64{},
		},
		{
			name:  "retried target counted once by its last outcome",
			rpc:   vald.SearchRPCName,
			calls: []string{"a", "a", "b", "b", "c"},
			replies: map[string][]error{
				"a": {unavailable, unavailable},
				"b": {unavailable, nil},
				"c": {nil},
			},
			wantCode: codes.OK,
			want: map[string]map[attribute.Distinct]int64{
				lbmetrics.SearchAgentFailureMetricsName: {
					failureKey(vald.SearchRPCName, codes.Unavailable): 1,
				},
				lbmetrics.SearchPartialSuccessMetricsName: {
					partialKey(vald.SearchRPCName, false): 1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			clients := map[string]*regClient{"a": new(regClient), "b": new(regClient), "c": new(regClient)}
			targets := make(map[vald.Client]string, len(clients))
			for target, c := range clients {
				targets[c] = target
			}
			// The fake BroadCast calls f sequentially, so the reply cursor needs no lock.
			next := make(map[string]int, len(clients))
			gw := &regGateway{
				agents: len(clients),
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					// models service.Gateway.BroadCast: the last callback error per target is returned.
					last := make(map[string]error, len(clients))
					for _, target := range test.calls {
						last[target] = f(ctx, target+":8081", clients[target])
					}
					if test.waitDeadline {
						// bounded: the request ctx carries the server timeout.
						<-ctx.Done()
						return ctx.Err()
					}
					errs := make([]error, 0, len(last))
					for _, target := range []string{"a", "b", "c"} {
						errs = append(errs, last[target])
					}
					return errors.Join(errs...)
				},
			}
			sm := lbmetrics.New()
			timeout := 10 * time.Second
			if test.waitDeadline {
				timeout = 50 * time.Millisecond
			}
			s := &server{gateway: gw, replica: 2, timeout: timeout, searchMetrics: sm}
			res, err := runTestSearch(t.Context(), s, test.rpc, &payload.Search_Config{Num: 1},
				func(_ context.Context, _ *payload.Search_Config, vc vald.Client, _ ...grpc.CallOption) (*payload.Search_Response, error) {
					target := targets[vc]
					i := next[target]
					next[target]++
					if i >= len(test.replies[target]) {
						t.Errorf("target %s called %d times, want at most %d", target, i+1, len(test.replies[target]))
						return nil, status.Error(codes.Internal, "unexpected call")
					}
					if rerr := test.replies[target][i]; rerr != nil {
						return nil, rerr
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + target, Distance: 0.1}},
					}, nil
				})
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(res.GetResults()) == 0 {
					t.Fatal("got no results, want at least one")
				}
			} else {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
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

// Test_server_Search_quorumKeepsQuorumResults verifies that an early-terminated
// search succeeds with the results of the answered agents and does not wait
// for the straggler.
func Test_server_Search_quorumKeepsQuorumResults(t *testing.T) {
	t.Parallel()
	const num = 10
	s := newStragglerServer(6, 3, time.Millisecond, 5*time.Second)
	WithSearchQuorumGrace("10ms")(s)
	start := time.Now()
	res, err := s.Search(t.Context(), &payload.Search_Request{
		Vector: make([]float32, 8),
		Config: &payload.Search_Config{Num: num, Timeout: int64(10 * time.Second)},
	})
	if err != nil {
		t.Fatalf("Search() err = %v", err)
	}
	if got := len(res.GetResults()); got != num {
		t.Errorf("Search() returned %d results, want %d", got, num)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Search() took %v, want it to stop waiting for the 5s straggler", elapsed)
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
		if in != seen[0] {
			t.Errorf("target got request %p, want shared forwarding request %p", in, seen[0])
		}
	}
}

func Test_server_LinearSearch_sharesRequestAcrossTargets(t *testing.T) {
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
					errs[i] = f(ctx, "127.0.0.1:8081", &recordingLinearClient{
						onLinear: func(in *payload.Search_Request) {
							mu.Lock()
							seen = append(seen, in)
							mu.Unlock()
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
	if _, err := s.LinearSearch(t.Context(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seen) != agents {
		t.Fatalf("got %d agent requests, want %d", len(seen), agents)
	}
	for _, in := range seen {
		if in == req || in.GetConfig() == cfg {
			t.Error("caller's request or config was aliased")
		}
		if in != seen[0] {
			t.Errorf("target got request %p, want shared forwarding request %p", in, seen[0])
		}
	}
}

type recordingLinearClient struct {
	*regClient
	onLinear func(*payload.Search_Request)
}

func (c *recordingLinearClient) LinearSearch(
	_ context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	c.onLinear(in)
	return probeResponse(), nil
}

func Test_server_doSearch_noDoubleQueryOnEmptyResult(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var hitCalls, emptyCalls int
	gw := &regGateway{
		agents: 2,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return errors.Join(
				f(ctx, "empty", new(regClient)),
				f(ctx, "hit", new(regClient)),
			)
		},
	}
	s := &server{gateway: gw, replica: 1, timeout: time.Second}
	res, err := runTestSearch(t.Context(), s, vald.SearchRPCName, &payload.Search_Config{Num: 1},
		func(_ context.Context, _ *payload.Search_Config, _ vald.Client, _ ...grpc.CallOption) (*payload.Search_Response, error) {
			mu.Lock()
			defer mu.Unlock()
			if emptyCalls == 0 {
				emptyCalls++
				return &payload.Search_Response{}, nil
			}
			hitCalls++
			return &payload.Search_Response{
				Results: []*payload.Object_Distance{{Id: "hit", Distance: 0.1}},
			}, nil
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.GetResults()) != 1 {
		t.Fatalf("got %d results, want 1", len(res.GetResults()))
	}
	if emptyCalls != 1 || hitCalls != 1 {
		t.Errorf("calls = (empty:%d, hit:%d), want (1, 1) with no empty-result re-query", emptyCalls, hitCalls)
	}
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
	if _, err := runTestSearch(t.Context(), s, vald.SearchRPCName, &payload.Search_Config{
		Num:   3,
		Ratio: types.Float(0.5),
	}, func(context.Context, *payload.Search_Config, vald.Client, ...grpc.CallOption) (*payload.Search_Response, error) {
		return new(payload.Search_Response), nil
	}); err == nil {
		t.Error("expected an error for an empty search result")
	}
	if got := gw.agentCalls.Load(); got != 1 {
		t.Errorf("GetAgentCount called %d times, want 1", got)
	}
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

func Test_selectAggregator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		checkAggr func(Aggregator) bool
		name      string
		algo      payload.Search_AggregationAlgorithm
	}{
		{
			name: "unknown maps to std (default)",
			algo: payload.Search_Unknown,
			checkAggr: func(a Aggregator) bool {
				_, ok := a.(*valdStdAggr)
				return ok
			},
		},
		{
			name: "concurrent queue maps to std",
			algo: payload.Search_ConcurrentQueue,
			checkAggr: func(a Aggregator) bool {
				_, ok := a.(*valdStdAggr)
				return ok
			},
		},
		{
			name: "sort slice maps to slice",
			algo: payload.Search_SortSlice,
			checkAggr: func(a Aggregator) bool {
				_, ok := a.(*valdSliceAggr)
				return ok
			},
		},
		{
			name: "sort pool slice maps to pool slice",
			algo: payload.Search_SortPoolSlice,
			checkAggr: func(a Aggregator) bool {
				_, ok := a.(*valdPoolSliceAggr)
				return ok
			},
		},
		{
			name: "pairing heap maps to pairing heap",
			algo: payload.Search_PairingHeap,
			checkAggr: func(a Aggregator) bool {
				_, ok := a.(*valdPairingHeapAggr)
				return ok
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := selectAggregator(test.algo, 10, 10, 3)
			if !test.checkAggr(got) {
				tt.Errorf("selectAggregator(%v) returned unexpected type %T", test.algo, got)
			}
		})
	}
}

type fallbackIDClient struct {
	*regClient
	onSearchID func(*payload.Search_IDRequest)
}

func (c *fallbackIDClient) SearchByID(
	_ context.Context, in *payload.Search_IDRequest, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	c.onSearchID(in)
	return probeResponse(), nil
}

func Test_server_Search_subsetRoutingAndQuorum(t *testing.T) {
	t.Parallel()
	addrs := []string{"a0", "a1", "a2", "a3", "a4"}
	vec := []float32{1, 2, 3}

	t.Run("Search with searchQuorumGrace=0 uses 3/5 subset and records 0 metrics", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for _, addr := range addrs {
			clients[addr] = &regClient{
				search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "id-" + addr, Distance: 0.1}},
					}, nil
				},
			}
		}
		gw := newSubsetRegGateway(addrs, clients)
		sm := lbmetrics.New()
		s := &server{gateway: gw, replica: 3, timeout: time.Second, searchMetrics: sm}
		res, err := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 10},
		})
		if err != nil || len(res.GetResults()) != 3 {
			t.Fatalf("Search = (%d, %v), want 3 results", len(res.GetResults()), err)
		}
		gw.mu.Lock()
		if len(gw.waves) != 1 || gw.waves[0].count != 3 || !slices.Equal(gw.waves[0].targets, []string{"a0", "a1", "a2"}) {
			gw.mu.Unlock()
			t.Fatalf("waves = %+v, want single wave of 3", gw.waves)
		}
		gw.mu.Unlock()
		if metrics := collectSearchMetrics(t, sm); len(metrics) != 0 {
			t.Errorf("expected 0 failure metrics, got %v", metrics)
		}
	})

	t.Run("Search 3/5 falls back to remaining agents on error or insufficient MinNum", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for i, addr := range addrs {
			idx := i
			clients[addr] = &regClient{
				search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					if idx < 3 {
						return &payload.Search_Response{
							Results: []*payload.Object_Distance{{Id: "dup", Distance: 0.1}},
						}, nil
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "wave2", Distance: 0.2}},
					}, nil
				},
			}
		}
		gw := newSubsetRegGateway(addrs, clients)
		s := &server{gateway: gw, replica: 3, timeout: time.Second}
		res, err := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 5, MinNum: 2},
		})
		if err != nil || len(res.GetResults()) != 2 {
			t.Fatalf("Search = (%v, %v), want 2 results", res.GetResults(), err)
		}
		gw.mu.Lock()
		defer gw.mu.Unlock()
		if len(gw.waves) != 2 || gw.waves[0].count != 3 || gw.waves[1].count != 2 {
			t.Fatalf("waves = %+v, want [3, 2]", gw.waves)
		}
	})

	t.Run("Search with Ratio uses full broadcast", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for _, addr := range addrs {
			clients[addr] = &regClient{
				search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					return probeResponse(), nil
				},
			}
		}
		gw := newSubsetRegGateway(addrs, clients)
		s := &server{gateway: gw, replica: 3, timeout: time.Second}
		res, err := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 5, Ratio: types.Float(0.5)},
		})
		if err != nil || len(res.GetResults()) == 0 {
			t.Fatalf("Search = (%v, %v)", res, err)
		}
		gw.mu.Lock()
		defer gw.mu.Unlock()
		if len(gw.waves) != 0 || gw.bcCalls != 1 {
			t.Fatalf("waves = %d, bcCalls = %d, want 0 waves and 1 full broadcast", len(gw.waves), gw.bcCalls)
		}
	})

	t.Run("searchQuorumGrace > 0 uses BroadCast + quorum without recording cancellation failures", func(t *testing.T) {
		t.Parallel()
		sm := lbmetrics.New()
		s := newStragglerServer(6, 3, time.Millisecond, 2*time.Second)
		s.searchMetrics = sm
		WithSearchQuorumGrace("5ms")(s)
		res, err := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 10, Timeout: int64(5 * time.Second)},
		})
		if err != nil || len(res.GetResults()) != 10 {
			t.Fatalf("Search = (%v, %v)", res, err)
		}
		if metrics := collectSearchMetrics(t, sm); len(metrics) != 0 {
			t.Errorf("expected 0 failure metrics on quorum cancellation, got %v", metrics)
		}
	})

	t.Run("Wave 2 ErrGRPCClientConnNotFound(*) preserves Wave 1 partial results", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for i, addr := range addrs {
			idx := i
			clients[addr] = &regClient{
				search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					if idx < 3 {
						return &payload.Search_Response{
							Results: []*payload.Object_Distance{{Id: "wave1-" + addr, Distance: 0.1}},
						}, nil
					}
					return nil, errGRPCClientConnNotFoundStar
				},
			}
		}
		gw := newSubsetRegGateway(addrs, clients)
		s := &server{gateway: gw, replica: 3, timeout: time.Second}
		res, err := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			// MinNum=5 forces Wave 2 when Wave 1 returned only 3 unique results.
			Config: &payload.Search_Config{Num: 5, MinNum: 5},
		})
		if st, _ := status.FromError(err); st == nil || st.Code() != codes.NotFound {
			t.Fatalf("err = %v, want NotFound (insufficient results) rather than search API connection not found", err)
		}
		_ = res

		// When Wave 1 fulfills Num/MinNum after a partial Wave 1 trigger (wave 2 no conn),
		// Wave 1 results must be returned instead of Internal.
		clients["a2"] = &regClient{
			search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
				return nil, status.Error(codes.Unavailable, "down")
			},
		}
		res2, err2 := s.Search(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 5},
		})
		if err2 != nil || len(res2.GetResults()) != 2 {
			t.Fatalf("Search = (%v, %v), want 2 Wave 1 results preserved", res2, err2)
		}
	})
}

func Test_server_SearchByID_fallbackSharesIDRequest(t *testing.T) {
	t.Parallel()
	const agents = 4
	var (
		mu   sync.Mutex
		seen []*payload.Search_IDRequest
	)
	gw := &regGateway{
		agents: agents,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			errs := make([]error, 0, agents)
			for range agents {
				errs = append(errs, f(ctx, "127.0.0.1:8081", &fallbackIDClient{
					regClient: &regClient{
						getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
							return nil, status.Error(codes.NotFound, "miss")
						},
					},
					onSearchID: func(in *payload.Search_IDRequest) {
						mu.Lock()
						seen = append(seen, in)
						mu.Unlock()
					},
				}))
			}
			return errors.Join(errs...)
		},
	}
	s := &server{gateway: gw, replica: 1, timeout: time.Second}
	cfg := &payload.Search_Config{Num: 3, MinNum: 1}
	res, err := s.SearchByID(t.Context(), &payload.Search_IDRequest{Id: "uuid-fb", Config: cfg})
	if err != nil || len(res.GetResults()) == 0 {
		t.Fatalf("SearchByID = (%v, %v)", res, err)
	}
	if len(seen) != agents {
		t.Fatalf("got %d SearchByID calls, want %d", len(seen), agents)
	}
	for _, in := range seen {
		if in != seen[0] {
			t.Errorf("target got %p, want shared hoisted Search_IDRequest %p", in, seen[0])
		}
		if in.GetConfig() == cfg || in.GetConfig().GetMinNum() != 0 {
			t.Errorf("forwarding config not cloned/sanitized: %+v", in.GetConfig())
		}
	}
}
