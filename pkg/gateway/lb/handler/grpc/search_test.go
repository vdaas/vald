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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/sync"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
			res, _, err := s.doSearch(t.Context(), vald.SearchRPCName, &payload.Search_Config{Num: 1},
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
			res, _, err := s.doSearch(t.Context(), test.rpc, &payload.Search_Config{Num: 1},
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

// NOT IMPLEMENTED BELOW
//
// func Test_server_Search(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
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
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
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
// func Test_server_calculateNum(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		num   uint32
// 		ratio float32
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantN uint32
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN uint32) error {
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
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
// 		           num:0,
// 		           ratio:0,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           num:0,
// 		           ratio:0,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotN := s.calculateNum(test.args.ctx, test.args.num, test.args.ratio)
// 			if err := checkFunc(test.want, gotN); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doSearch(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		cfg *payload.Search_Config
// 		f   func(ctx context.Context, cfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes   *payload.Search_Response
// 		wantAttrs []attribute.KeyValue
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, []attribute.KeyValue, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, gotAttrs []attribute.KeyValue, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		if !reflect.DeepEqual(gotAttrs, w.wantAttrs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAttrs, w.wantAttrs)
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
// 		           cfg:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           cfg:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, gotAttrs, err := s.doSearch(test.args.ctx, test.args.cfg, test.args.f)
// 			if err := checkFunc(test.want, gotRes, gotAttrs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_selectAggregator(t *testing.T) {
// 	type args struct {
// 		algo    payload.Search_AggregationAlgorithm
// 		num     int
// 		fnum    int
// 		replica int
// 	}
// 	type want struct {
// 		want Aggregator
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Aggregator) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Aggregator) error {
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
// 		           algo:nil,
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
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
// 		           algo:nil,
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
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
// 			got := selectAggregator(test.args.algo, test.args.num, test.args.fnum, test.args.replica)
// 			if err := checkFunc(test.want, got); err != nil {
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
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
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
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
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
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
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
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
