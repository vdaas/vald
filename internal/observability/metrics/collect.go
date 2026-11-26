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

package metrics

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CollectInt64Sums registers m on an isolated manual-reader MeterProvider
// configured with m's views, collects all recorded datapoints, and returns the
// cumulative int64 sum values keyed by instrument name and attribute set key.
func CollectInt64Sums(
	ctx context.Context, m Metric,
) (res map[string]map[attribute.Distinct]int64, err error) {
	views, err := m.View()
	if err != nil {
		return nil, err
	}
	reader := sdkmetric.NewManualReader()
	opts := make([]sdkmetric.Option, 0, len(views)+1)
	opts = append(opts, sdkmetric.WithReader(reader))
	for _, vw := range views {
		opts = append(opts, sdkmetric.WithView(vw))
	}
	mp := sdkmetric.NewMeterProvider(opts...)
	defer func() {
		if serr := mp.Shutdown(context.WithoutCancel(ctx)); serr != nil {
			err = errors.Join(err, serr)
		}
	}()
	if err := m.Register(mp.Meter("test")); err != nil {
		return nil, err
	}
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		return nil, err
	}
	got := make(map[string]map[attribute.Distinct]int64)
	for _, scm := range rm.ScopeMetrics {
		for _, sm := range scm.Metrics {
			s, ok := sm.Data.(metricdata.Sum[int64])
			if !ok {
				return nil, errors.Errorf("%s: data = %T, want metricdata.Sum[int64]", sm.Name, sm.Data)
			}
			if !s.IsMonotonic || s.Temporality != metricdata.CumulativeTemporality {
				return nil, errors.Errorf("%s: monotonic = %v, temporality = %v, want monotonic cumulative", sm.Name, s.IsMonotonic, s.Temporality)
			}
			pts := make(map[attribute.Distinct]int64, len(s.DataPoints))
			for _, dp := range s.DataPoints {
				pts[dp.Attributes.Equivalent()] = dp.Value
			}
			got[sm.Name] = pts
		}
	}
	return got, nil
}
