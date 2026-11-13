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
package lb

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// collectSums registers sm on a meter provider built with sm's views and returns
// the cumulative sum per instrument name and attribute set. It fails the test when
// an instrument is not a monotonic cumulative Sum[int64] (the view did not match).
func collectSums(t *testing.T, sm *SearchMetrics) map[string]map[attribute.Distinct]int64 {
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
			if m.Unit != metrics.Dimensionless {
				t.Errorf("%s: unit = %q, want %q", m.Name, m.Unit, metrics.Dimensionless)
			}
			s, ok := m.Data.(metricdata.Sum[int64])
			if !ok {
				t.Errorf("%s: data = %T, want metricdata.Sum[int64]", m.Name, m.Data)
				continue
			}
			if !s.IsMonotonic || s.Temporality != metricdata.CumulativeTemporality {
				t.Errorf("%s: monotonic = %v, temporality = %v, want monotonic cumulative", m.Name, s.IsMonotonic, s.Temporality)
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

func TestSearchMetrics_Register(t *testing.T) {
	t.Parallel()
	failure := func(rpc string, code codes.Code) attribute.Distinct {
		s := attribute.NewSet(attribute.String(rpcKey, rpc), attribute.String(grpcCodeKey, code.String()))
		return s.Equivalent()
	}
	partial := func(rpc string, risk bool) attribute.Distinct {
		s := attribute.NewSet(attribute.String(rpcKey, rpc), attribute.Bool(recallRiskKey, risk))
		return s.Equivalent()
	}
	writePartial := func(rpc string) attribute.Distinct {
		s := attribute.NewSet(attribute.String(rpcKey, rpc))
		return s.Equivalent()
	}
	tests := []struct {
		record func(*SearchMetrics)
		want   map[string]map[attribute.Distinct]int64
		name   string
	}{
		{
			name:   "nothing recorded",
			record: func(*SearchMetrics) {},
			want:   map[string]map[attribute.Distinct]int64{},
		},
		{
			name: "counts are cumulative per attribute set",
			record: func(sm *SearchMetrics) {
				sm.RecordAgentFailure("Search", codes.Unavailable)
				sm.RecordAgentFailure("Search", codes.Unavailable)
				sm.RecordAgentFailure("Search", codes.Aborted)
				sm.RecordAgentFailure("LinearSearch", codes.Unavailable)
				sm.RecordPartialSuccess("Search", false)
				sm.RecordPartialSuccess("Search", false)
				sm.RecordPartialSuccess("Search", true)
				sm.RecordWriteAgentFailure("Insert", codes.Internal)
				sm.RecordWriteAgentFailure("Insert", codes.Internal)
				sm.RecordWriteAgentFailure("Update", codes.Unavailable)
				sm.RecordWritePartialSuccess("Insert")
				sm.RecordWritePartialSuccess("Insert")
				sm.RecordWritePartialSuccess("Update")
			},
			want: map[string]map[attribute.Distinct]int64{
				SearchAgentFailureMetricsName: {
					failure("Search", codes.Unavailable):       2,
					failure("Search", codes.Aborted):           1,
					failure("LinearSearch", codes.Unavailable): 1,
				},
				SearchPartialSuccessMetricsName: {
					partial("Search", false): 2,
					partial("Search", true):  1,
				},
				WriteAgentFailureMetricsName: {
					failure("Insert", codes.Internal):    2,
					failure("Update", codes.Unavailable): 1,
				},
				WritePartialSuccessMetricsName: {
					writePartial("Insert"): 2,
					writePartial("Update"): 1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			sm := New()
			test.record(sm)
			got := collectSums(t, sm)
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

func TestSearchMetrics_nilIsNoop(t *testing.T) {
	t.Parallel()
	var sm *SearchMetrics
	sm.RecordAgentFailure("Search", codes.Unavailable)
	sm.RecordPartialSuccess("Search", true)
	sm.RecordWriteAgentFailure("Insert", codes.Internal)
	sm.RecordWritePartialSuccess("Insert")
}
