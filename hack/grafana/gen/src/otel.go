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
			).
			WithTarget(
				prometheusQuery("otelcol_exporter_queue_capacity{namespace=\"$Namespace\", pod=~\"$ReplicaSet.*\"}").
					LegendFormat("Queue Capacity - {{pod}}"),
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
