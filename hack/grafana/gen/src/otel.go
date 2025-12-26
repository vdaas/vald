//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func addOtelQueuePanel(builder *dashboard.DashboardBuilder) {
	builder.WithPanel(
		timeseries.NewPanelBuilder().
			Title("OpenTelemetry Queue Metrics").
			Description("OpenTelemetry Collector queue size and capacity").
			Unit("short").
			WithTarget(
				prometheusQuery("otelcol_exporter_queue_size{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}").
					LegendFormat("Queue Size - {{pod}}"),
			),
	)
}

func addOtelSpanDropPanel(builder *dashboard.DashboardBuilder) {
	builder.WithPanel(
		timeseries.NewPanelBuilder().
			Title("OpenTelemetry Span Drop Rate").
			Description("Rate of dropped spans").
			Unit("ops").
			WithTarget(
				prometheusQuery("rate(otelcol_exporter_enqueue_failed_spans{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}[5m])").
					LegendFormat("Failed Spans/sec - {{pod}}"),
			),
	)
}

func addOtelMetricDropPanel(builder *dashboard.DashboardBuilder) {
	builder.WithPanel(
		timeseries.NewPanelBuilder().
			Title("OpenTelemetry Metric Drop Rate").
			Description("Rate of dropped metric points").
			Unit("ops").
			WithTarget(
				prometheusQuery("rate(otelcol_exporter_enqueue_failed_metric_points{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}[5m])").
					LegendFormat("Failed Metrics/sec - {{pod}}"),
			),
	)
}

func addOtelNetworkPanel(builder *dashboard.DashboardBuilder) {
	builder.WithPanel(
		timeseries.NewPanelBuilder().
			Title("OpenTelemetry Network I/O").
			Description("Network receive and transmit rates").
			Unit("Bps").
			WithTarget(
				prometheusQuery("rate(container_network_receive_bytes_total{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}[5m])").
					LegendFormat("Receive - {{pod}}"),
			).
			WithTarget(
				prometheusQuery("rate(container_network_transmit_bytes_total{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}[5m])").
					LegendFormat("Transmit - {{pod}}"),
			),
	)
}
