// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package malloc

import (
	"context"

	"github.com/vdaas/vald/internal/core/malloc"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	totalFastCountMetricsName        = "total_fast_count"
	totalFastCountMetricsDescription = "Total fast count"

	totalFastSizeMetricsName        = "total_fast_size"
	totalFastSizeMetricsDescription = "Total fast size"

	totalRestCountMetricsName        = "total_rest_count"
	totalRestCountMetricsDescription = "Total rest count"

	totalRestSizeMetricsName        = "total_rest_size"
	totalRestSizeMetricsDescription = "Total rest size"

	systemCurrentSizeMetricsName        = "system_current_size"
	systemCurrentSizeMetricsDescription = "System current size"

	systemMaxSizeMetricsName        = "system_max_size"
	systemMaxSizeMetricsDescription = "System max size"

	aspaceTotalSizeMetricsName        = "aspace_total_size"
	aspaceTotalSizeMetricsDescription = "Aspace total size"

	aspaceMprotectSizeMetricsName        = "aspace_mprotect_size"
	aspaceMprotectSizeMetricsDescription = "Aspace mprotect size"
)

type mallocMetrics struct{}

func New() metrics.Metric {
	return &mallocMetrics{}
}

func (*mallocMetrics) View() ([]*metrics.View, error) {
	totalFastCount, err := view.New(
		view.MatchInstrumentName(totalFastCountMetricsName),
		view.WithSetDescription(totalFastCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	totalFastSize, err := view.New(
		view.MatchInstrumentName(totalFastSizeMetricsName),
		view.WithSetDescription(totalFastSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	totalRestCount, err := view.New(
		view.MatchInstrumentName(totalRestCountMetricsName),
		view.WithSetDescription(totalRestCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	totalRestSize, err := view.New(
		view.MatchInstrumentName(totalRestSizeMetricsName),
		view.WithSetDescription(totalRestSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	systemCurrentSize, err := view.New(
		view.MatchInstrumentName(systemCurrentSizeMetricsName),
		view.WithSetDescription(systemCurrentSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	systemMaxSize, err := view.New(
		view.MatchInstrumentName(systemMaxSizeMetricsName),
		view.WithSetDescription(systemMaxSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	aspaceTotalSize, err := view.New(
		view.MatchInstrumentName(aspaceTotalSizeMetricsName),
		view.WithSetDescription(aspaceTotalSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	aspaceMprotectSize, err := view.New(
		view.MatchInstrumentName(aspaceMprotectSizeMetricsName),
		view.WithSetDescription(aspaceMprotectSizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&totalFastCount,
		&totalFastSize,
		&totalRestCount,
		&totalRestSize,
		&systemCurrentSize,
		&systemMaxSize,
		&aspaceTotalSize,
		&aspaceMprotectSize,
	}, nil
}

func (*mallocMetrics) Register(m metrics.Meter) error {
	totalFastCount, err := m.AsyncInt64().Gauge(
		totalFastCountMetricsName,
		metrics.WithDescription(totalFastCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	totalFastSize, err := m.AsyncInt64().Gauge(
		totalFastSizeMetricsName,
		metrics.WithDescription(totalFastSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	totalRestCount, err := m.AsyncInt64().Gauge(
		totalRestCountMetricsName,
		metrics.WithDescription(totalRestCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	totalRestSize, err := m.AsyncInt64().Gauge(
		totalRestSizeMetricsName,
		metrics.WithDescription(totalRestSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	systemCurrentSize, err := m.AsyncInt64().Gauge(
		systemCurrentSizeMetricsName,
		metrics.WithDescription(systemCurrentSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	systemMaxSize, err := m.AsyncInt64().Gauge(
		systemMaxSizeMetricsName,
		metrics.WithDescription(systemMaxSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	aspaceTotalSize, err := m.AsyncInt64().Gauge(
		aspaceTotalSizeMetricsName,
		metrics.WithDescription(aspaceTotalSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	aspaceMprotectSize, err := m.AsyncInt64().Gauge(
		aspaceMprotectSizeMetricsName,
		metrics.WithDescription(aspaceMprotectSizeMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			totalFastCount,
			totalFastSize,
			totalRestCount,
			totalRestSize,
			systemCurrentSize,
			systemMaxSize,
			aspaceTotalSize,
			aspaceMprotectSize,
		},
		func(ctx context.Context) {
			mallocInfo, err := malloc.GetMallocInfo()
			if err != nil {
				return
			}

			totalFastCount.Observe(ctx, int64(mallocInfo.Total[0].Count))
			totalFastSize.Observe(ctx, int64(mallocInfo.Total[0].Size))
			totalRestCount.Observe(ctx, int64(mallocInfo.Total[1].Count))
			totalRestSize.Observe(ctx, int64(mallocInfo.Total[1].Size))
			systemCurrentSize.Observe(ctx, int64(mallocInfo.System[0].Size))
			systemMaxSize.Observe(ctx, int64(mallocInfo.System[1].Size))
			aspaceTotalSize.Observe(ctx, int64(mallocInfo.Aspace[0].Size))
			aspaceMprotectSize.Observe(ctx, int64(mallocInfo.Aspace[1].Size))
		},
	)
}
