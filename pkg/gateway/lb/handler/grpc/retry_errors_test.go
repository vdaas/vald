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
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// scriptedInsertClient answers the i-th Insert with attempts[i]: a nil entry is a
// success located on host, anything else is returned as the agent's error.
type scriptedInsertClient struct {
	*regClient
	host     string
	attempts []error
	calls    int
}

func (c *scriptedInsertClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	i := c.calls
	c.calls++
	if i < len(c.attempts) && c.attempts[i] != nil {
		return nil, c.attempts[i]
	}
	return &payload.Object_Location{Uuid: in.GetVector().GetId(), Name: c.host, Ips: []string{c.host}}, nil
}

// Insert used to keep an agent's first failure after a retry of the same agent
// (client backoff, or DoMulti's second pass) succeeded, failing a request whose
// replicas were all written.
func Test_server_Insert_retriedTargetError(t *testing.T) {
	t.Parallel()
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	internal := status.Error(codes.Internal, "agent internal")
	exists := status.Error(codes.AlreadyExists, "already exists")
	const (
		uuid     = "uuid"
		wantName = "a,b"
	)
	tests := []struct {
		attempts map[string][]error
		name     string
		wantName string
		order    []string // targets in the order DoMulti calls f
		wantIPs  int
		wantErr  bool
	}{
		{
			name:     "backoff retry of the same target succeeds",
			attempts: map[string][]error{"a": {unavailable, nil}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			name:     "second pass retries a failed target and it succeeds",
			attempts: map[string][]error{"a": {internal, nil}, "b": {nil}},
			order:    []string{"a", "b", "a"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			name:     "retry answers AlreadyExists",
			attempts: map[string][]error{"a": {unavailable, exists}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			// another agent stored the vector, so the request succeeds under-replicated.
			name:     "every attempt of a target fails",
			attempts: map[string][]error{"a": {unavailable, internal}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: "b",
			wantIPs:  1,
		},
		{
			name:     "every attempt of every target fails",
			attempts: map[string][]error{"a": {unavailable, internal}, "b": {internal}},
			order:    []string{"a", "a", "b"},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			clients := make(map[string]vald.Client, len(test.attempts))
			for target, attempts := range test.attempts {
				clients[target] = &scriptedInsertClient{regClient: new(regClient), host: target, attempts: attempts}
			}
			gw := &regGateway{
				doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					// The real DoMulti/OrderedRange only log per-target failures.
					for _, target := range test.order {
						if err := f(ctx, target, clients[target]); err != nil {
							t.Logf("target %s: %v", target, err)
						}
					}
					return nil
				},
			}
			s := &server{gateway: gw, replica: 2}
			loc, err := s.Insert(t.Context(), &payload.Insert_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 1},
			})
			if test.wantErr {
				if err == nil {
					t.Fatalf("Insert returned %v, want the failed target's error", loc)
				}
				if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.Internal {
					t.Errorf("status code = %v (ok=%t), want the last attempt's Internal: %v", st.Code(), ok, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Insert returned %v, want success after the retry", err)
			}
			if loc.GetName() != test.wantName {
				t.Errorf("Name = %q, want %q", loc.GetName(), test.wantName)
			}
			if len(loc.GetIps()) != test.wantIPs {
				t.Errorf("Ips = %v, want %d entries", loc.GetIps(), test.wantIPs)
			}
		})
	}
}

// The averaged fields used to be divided by every entry, invalid ones included,
// while the medians are taken over the valid entries only.
func Test_mergeInfoIndexStatistics_averagesOverValidEntries(t *testing.T) {
	t.Parallel()
	valid := func(v float64) *payload.Info_Index_Statistics {
		return &payload.Info_Index_Statistics{
			Valid:                          true,
			ModeIndegree:                   uint64(v),
			ModeOutdegree:                  uint64(v),
			VarianceOfIndegree:             v,
			VarianceOfOutdegree:            v,
			MeanEdgeLength:                 v,
			MeanEdgeLengthFor10Edges:       v,
			MeanIndegreeDistanceFor10Edges: v,
			MeanNumberOfEdgesPerNode:       v,
			C1Indegree:                     v,
			C5Indegree:                     v,
			C95Outdegree:                   v,
			C99Outdegree:                   v,
		}
	}
	tests := []struct {
		stats map[string]*payload.Info_Index_Statistics
		name  string
		want  float64
	}{
		{
			name:  "one valid one invalid",
			stats: map[string]*payload.Info_Index_Statistics{"a": valid(8), "b": {Valid: false}},
			want:  8,
		},
		{
			name: "two valid two invalid",
			stats: map[string]*payload.Info_Index_Statistics{
				"a": valid(4), "b": valid(8), "c": {Valid: false}, "d": {Valid: false},
			},
			want: 6,
		},
		{
			name:  "all valid",
			stats: map[string]*payload.Info_Index_Statistics{"a": valid(4), "b": valid(8)},
			want:  6,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := mergeInfoIndexStatistics(test.stats)
			for name, v := range map[string]float64{
				"ModeIndegree":                   float64(got.GetModeIndegree()),
				"ModeOutdegree":                  float64(got.GetModeOutdegree()),
				"VarianceOfIndegree":             got.GetVarianceOfIndegree(),
				"VarianceOfOutdegree":            got.GetVarianceOfOutdegree(),
				"MeanEdgeLength":                 got.GetMeanEdgeLength(),
				"MeanEdgeLengthFor10Edges":       got.GetMeanEdgeLengthFor10Edges(),
				"MeanIndegreeDistanceFor10Edges": got.GetMeanIndegreeDistanceFor10Edges(),
				"MeanNumberOfEdgesPerNode":       got.GetMeanNumberOfEdgesPerNode(),
				"C1Indegree":                     got.GetC1Indegree(),
				"C5Indegree":                     got.GetC5Indegree(),
				"C95Outdegree":                   got.GetC95Outdegree(),
				"C99Outdegree":                   got.GetC99Outdegree(),
			} {
				if v != test.want {
					t.Errorf("%s = %v, want %v", name, v, test.want)
				}
			}
		})
	}
}
