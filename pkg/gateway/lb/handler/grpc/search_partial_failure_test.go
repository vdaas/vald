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
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func Test_searchRecallAtRisk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		failures int
		replica  int
		want     bool
	}{
		{name: "no failure", failures: 0, replica: 3, want: false},
		{name: "fewer failures than replica", failures: 2, replica: 3, want: false},
		{name: "failures equal replica", failures: 3, replica: 3, want: true},
		{name: "failures above replica", failures: 4, replica: 3, want: true},
		{name: "single replica, one failure", failures: 1, replica: 1, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := searchRecallAtRisk(test.failures, test.replica); got != test.want {
				t.Errorf("searchRecallAtRisk(%d, %d) = %v, want %v", test.failures, test.replica, got, test.want)
			}
		})
	}
}

func Test_searchAgentFailure(t *testing.T) {
	t.Parallel()
	const target = "a:8081"
	tests := []struct {
		err      error
		name     string
		wantCode codes.Code
		want     bool
	}{
		{name: "nil", err: nil},
		{name: "context canceled", err: context.Canceled},
		{name: "context deadline", err: context.DeadlineExceeded},
		{name: "rpc call canceled", err: errors.ErrRPCCallFailed(target, context.Canceled)},
		{name: "rpc call deadline", err: errors.ErrRPCCallFailed(target, context.DeadlineExceeded)},
		{name: "status Canceled", err: status.Error(codes.Canceled, "c")},
		{name: "status DeadlineExceeded", err: status.Error(codes.DeadlineExceeded, "d")},
		{name: "status NotFound", err: status.Error(codes.NotFound, "nf")},
		{name: "status InvalidArgument is client-caused", err: status.Error(codes.InvalidArgument, "dim")},
		{name: "status Unavailable", err: status.Error(codes.Unavailable, "u"), wantCode: codes.Unavailable, want: true},
		{name: "status FailedPrecondition", err: status.Error(codes.FailedPrecondition, "fp"), wantCode: codes.FailedPrecondition, want: true},
		{name: "status Unimplemented", err: status.Error(codes.Unimplemented, "ui"), wantCode: codes.Unimplemented, want: true},
		{name: "non-status error", err: errors.New("boom"), wantCode: codes.Unknown, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			code, ok := searchAgentFailure(target, test.err)
			if ok != test.want || (ok && code != test.wantCode) {
				t.Errorf("searchAgentFailure(%v) = (%v, %v), want (%v, %v)", test.err, code, ok, test.wantCode, test.want)
			}
		})
	}
}

// The per-agent errors are logged at Debug; this per-request level decides
// whether the agent failures are a Warn.
func Test_partialSearchLogLevel(t *testing.T) {
	t.Parallel()
	bcErr := status.Error(codes.Unavailable, "agent down")
	tests := []struct {
		bcErr    error
		reqErr   error
		name     string
		failures int
		want     level.Level
	}{
		{name: "nothing failed", want: level.Unknown},
		{name: "failures below replica", failures: 1, bcErr: bcErr, want: level.DEBUG},
		{name: "failures below replica without broadcast error", failures: 1, want: level.DEBUG},
		{name: "failures reach replica", failures: 2, bcErr: bcErr, want: level.WARN},
		{name: "broadcast error without agent failure", bcErr: context.Canceled, want: level.WARN},
		{name: "failed request with agent failures", failures: 1, reqErr: bcErr, want: level.WARN},
		{name: "failed request without agent failures", reqErr: bcErr, want: level.Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := partialSearchLogLevel(test.failures, 2, test.bcErr, test.reqErr); got != test.want {
				t.Errorf("partialSearchLogLevel(%d, 2, %v, %v) = %v, want %v", test.failures, test.bcErr, test.reqErr, got, test.want)
			}
		})
	}
}

// A failed request logs exactly one Warn line: when the deferred summary logs
// the request error, aggregationSearch must not log it (or bcErr) again.
func Test_searchSummaryLogsRequestErr(t *testing.T) {
	t.Parallel()
	reqErr := status.Error(codes.Unavailable, "agent down")
	tests := []struct {
		reqErr   error
		name     string
		failures int
		want     bool
	}{
		{name: "successful request", failures: 2},
		{name: "failed request without agent failures", reqErr: reqErr},
		{name: "failed request with agent failures", failures: 1, reqErr: reqErr, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := searchSummaryLogsRequestErr(test.failures, test.reqErr)
			if got != test.want {
				t.Errorf("searchSummaryLogsRequestErr(%d, %v) = %v, want %v", test.failures, test.reqErr, got, test.want)
			}
			// the guard must agree with the summary's own decision, for any bcErr
			for _, bcErr := range []error{nil, reqErr} {
				summary := partialSearchLogLevel(test.failures, 2, bcErr, test.reqErr) == level.WARN && test.reqErr != nil
				if summary != got {
					t.Errorf("summary Warns the request error = %v, guard = %v (bcErr %v)", summary, got, bcErr)
				}
			}
		})
	}
}

// collectSearchMetrics returns the cumulative value per instrument name and
// attribute set currently reported by sm.
func collectSearchMetrics(
	t *testing.T, sm *lbmetrics.SearchMetrics,
) map[string]map[attribute.Distinct]int64 {
	t.Helper()
	views, err := sm.View()
	if err != nil {
		t.Fatalf("View() error = %v", err)
	}
	reader := sdkmetric.NewManualReader()
	opts := make([]sdkmetric.Option, 0, len(views)+1)
	opts = append(opts, sdkmetric.WithReader(reader))
	for _, vw := range views {
		opts = append(opts, sdkmetric.WithView(vw))
	}
	mp := sdkmetric.NewMeterProvider(opts...)
	t.Cleanup(func() {
		if err := mp.Shutdown(context.WithoutCancel(t.Context())); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})
	if err := sm.Register(mp.Meter("test")); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(t.Context(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	got := make(map[string]map[attribute.Distinct]int64)
	for _, scm := range rm.ScopeMetrics {
		for _, m := range scm.Metrics {
			s, ok := m.Data.(metricdata.Sum[int64])
			if !ok {
				t.Errorf("%s: data = %T, want metricdata.Sum[int64]", m.Name, m.Data)
				continue
			}
			pts := make(map[attribute.Distinct]int64, len(s.DataPoints))
			for _, dp := range s.DataPoints {
				pts[dp.Attributes.Equivalent()] = dp.Value
			}
			got[m.Name] = pts
		}
	}
	return got
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
