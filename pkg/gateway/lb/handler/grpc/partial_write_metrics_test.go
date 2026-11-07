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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"go.opentelemetry.io/otel/attribute"
)

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
