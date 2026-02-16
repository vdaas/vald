//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"context"
	"math"
	"time"

	"go.opentelemetry.io/otel/metric"

	obsmetrics "github.com/vdaas/vald/internal/observability/metrics"
)

type otelMetrics struct {
	collector Collector
}

func NewOTELMetrics(c Collector) obsmetrics.Metric {
	return &otelMetrics{
		collector: c,
	}
}

func (o *otelMetrics) View() ([]obsmetrics.View, error) {
	return nil, nil
}

func (o *otelMetrics) Register(m obsmetrics.Meter) error {
	totalRequests, err := m.Int64ObservableGauge(
		"e2e_total_requests",
		obsmetrics.WithDescription("Total number of requests"),
		obsmetrics.WithUnit(obsmetrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	totalErrors, err := m.Int64ObservableGauge(
		"e2e_total_errors",
		obsmetrics.WithDescription("Total number of errors"),
		obsmetrics.WithUnit(obsmetrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	// Latency metrics
	latP50, err := m.Float64ObservableGauge(
		"e2e_latency_p50",
		obsmetrics.WithDescription("Latency P50"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	latP90, err := m.Float64ObservableGauge(
		"e2e_latency_p90",
		obsmetrics.WithDescription("Latency P90"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	latP99, err := m.Float64ObservableGauge(
		"e2e_latency_p99",
		obsmetrics.WithDescription("Latency P99"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	latMax, err := m.Float64ObservableGauge(
		"e2e_latency_max",
		obsmetrics.WithDescription("Latency Max"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	latMin, err := m.Float64ObservableGauge(
		"e2e_latency_min",
		obsmetrics.WithDescription("Latency Min"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	// Queue Wait metrics
	qwP50, err := m.Float64ObservableGauge(
		"e2e_queue_wait_p50",
		obsmetrics.WithDescription("Queue Wait P50"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	qwP90, err := m.Float64ObservableGauge(
		"e2e_queue_wait_p90",
		obsmetrics.WithDescription("Queue Wait P90"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	qwP99, err := m.Float64ObservableGauge(
		"e2e_queue_wait_p99",
		obsmetrics.WithDescription("Queue Wait P99"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	qwMax, err := m.Float64ObservableGauge(
		"e2e_queue_wait_max",
		obsmetrics.WithDescription("Queue Wait Max"),
		obsmetrics.WithUnit(obsmetrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(func(_ context.Context, obs metric.Observer) error {
		snap := o.collector.GlobalSnapshot()
		if snap == nil {
			return nil
		}

		safeInt64 := func(v uint64) int64 {
			if v > math.MaxInt64 {
				return math.MaxInt64
			}
			return int64(v)
		}

		obs.ObserveInt64(totalRequests, safeInt64(snap.Total))
		obs.ObserveInt64(totalErrors, safeInt64(snap.Errors))

		// Latency
		if snap.LatPercentiles != nil {
			obs.ObserveFloat64(latP50, nsToMs(snap.LatPercentiles.Quantile(0.5)))
			obs.ObserveFloat64(latP90, nsToMs(snap.LatPercentiles.Quantile(0.9)))
			obs.ObserveFloat64(latP99, nsToMs(snap.LatPercentiles.Quantile(0.99)))
			obs.ObserveFloat64(latMax, nsToMs(snap.LatPercentiles.Quantile(1.0)))
			obs.ObserveFloat64(latMin, nsToMs(snap.LatPercentiles.Quantile(0.0)))
		}

		// Queue Wait
		if snap.QWPercentiles != nil {
			obs.ObserveFloat64(qwP50, nsToMs(snap.QWPercentiles.Quantile(0.5)))
			obs.ObserveFloat64(qwP90, nsToMs(snap.QWPercentiles.Quantile(0.9)))
			obs.ObserveFloat64(qwP99, nsToMs(snap.QWPercentiles.Quantile(0.99)))
			obs.ObserveFloat64(qwMax, nsToMs(snap.QWPercentiles.Quantile(1.0)))
		}
		return nil
	}, totalRequests, totalErrors, latP50, latP90, latP99, latMax, latMin, qwP50, qwP90, qwP99, qwMax)

	return err
}

func nsToMs(ns float64) float64 {
	return ns / float64(time.Millisecond)
}
