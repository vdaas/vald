// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
)

const (
	// Metric names
	TotalRequests = "e2e_total_requests"
	TotalErrors   = "e2e_total_errors"

	// Latency metrics
	LatencyP50 = "e2e_latency_p50"
	LatencyP90 = "e2e_latency_p90"
	LatencyP99 = "e2e_latency_p99"
	LatencyMax = "e2e_latency_max"
	LatencyMin = "e2e_latency_min"

	// Queue Wait metrics
	QueueWaitP50 = "e2e_queue_wait_p50"
	QueueWaitP90 = "e2e_queue_wait_p90"
	QueueWaitP99 = "e2e_queue_wait_p99"
	QueueWaitMax = "e2e_queue_wait_max"
)

type otelMetrics struct {
	collector Collector
}

// NewOTELMetrics creates a new OTLP metrics exporter that reads from the given Collector.
func NewOTELMetrics(c Collector) metrics.Metric {
	return &otelMetrics{
		collector: c,
	}
}

func (m *otelMetrics) View() ([]metrics.View, error) {
	// We return empty views here because we are using Observable Gauges which don't strictly require Views for basic aggregation.
	// If we wanted to change aggregation (e.g., Histogram buckets), we would define it here.
	// Since we are exporting pre-calculated percentiles as Gauges, LastValue aggregation is default and correct.
	return nil, nil
}

func (m *otelMetrics) Register(meter metrics.Meter) error {
	totalRequests, err := meter.Int64ObservableGauge(
		TotalRequests,
		metrics.WithDescription("Total number of requests executed in E2E test"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	totalErrors, err := meter.Int64ObservableGauge(
		TotalErrors,
		metrics.WithDescription("Total number of failed requests in E2E test"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	latencyP50, err := meter.Float64ObservableGauge(
		LatencyP50,
		metrics.WithDescription("E2E Request Latency P50"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	latencyP90, err := meter.Float64ObservableGauge(
		LatencyP90,
		metrics.WithDescription("E2E Request Latency P90"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	latencyP99, err := meter.Float64ObservableGauge(
		LatencyP99,
		metrics.WithDescription("E2E Request Latency P99"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	latencyMax, err := meter.Float64ObservableGauge(
		LatencyMax,
		metrics.WithDescription("E2E Request Latency Max"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	latencyMin, err := meter.Float64ObservableGauge(
		LatencyMin,
		metrics.WithDescription("E2E Request Latency Min"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	qwP50, err := meter.Float64ObservableGauge(
		QueueWaitP50,
		metrics.WithDescription("E2E Queue Wait P50"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	qwP90, err := meter.Float64ObservableGauge(
		QueueWaitP90,
		metrics.WithDescription("E2E Queue Wait P90"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	qwP99, err := meter.Float64ObservableGauge(
		QueueWaitP99,
		metrics.WithDescription("E2E Queue Wait P99"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}
	qwMax, err := meter.Float64ObservableGauge(
		QueueWaitMax,
		metrics.WithDescription("E2E Queue Wait Max"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			snap := m.collector.GlobalSnapshot()
			if snap == nil {
				return nil
			}

			o.ObserveInt64(totalRequests, int64(snap.Total))
			o.ObserveInt64(totalErrors, int64(snap.Errors))

			if snap.LatPercentiles != nil {
				o.ObserveFloat64(latencyP50, toMillis(snap.LatPercentiles.Quantile(0.5)))
				o.ObserveFloat64(latencyP90, toMillis(snap.LatPercentiles.Quantile(0.9)))
				o.ObserveFloat64(latencyP99, toMillis(snap.LatPercentiles.Quantile(0.99)))
				o.ObserveFloat64(latencyMax, toMillis(snap.LatPercentiles.Quantile(1.0)))
				o.ObserveFloat64(latencyMin, toMillis(snap.LatPercentiles.Quantile(0.0)))
			}

			if snap.QWPercentiles != nil {
				o.ObserveFloat64(qwP50, toMillis(snap.QWPercentiles.Quantile(0.5)))
				o.ObserveFloat64(qwP90, toMillis(snap.QWPercentiles.Quantile(0.9)))
				o.ObserveFloat64(qwP99, toMillis(snap.QWPercentiles.Quantile(0.99)))
				o.ObserveFloat64(qwMax, toMillis(snap.QWPercentiles.Quantile(1.0)))
			}

			return nil
		},
		totalRequests,
		totalErrors,
		latencyP50,
		latencyP90,
		latencyP99,
		latencyMax,
		latencyMin,
		qwP50,
		qwP90,
		qwP99,
		qwMax,
	)

	return err
}

func toMillis(ns float64) float64 {
	return ns / 1e6
}
