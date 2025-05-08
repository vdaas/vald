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
package main

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/promql-builder/go/promql"
	"github.com/vdaas/vald/internal/observability/metrics/agent/core/ngt/public_const"
)

func addCPUPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(promql.Vector(cpuMetric).
				Label("namespace", namespaceVariable).
				LabelMatchRegexp("container", "$ReplicaSet").
				LabelMatchRegexp("pod", "$PodName").
				LabelNeq("image", "").
				Range(intervalVariable))).By([]string{"pod"}).String(),
		).Format("time_series").LegendFormat("{{pod}}"))
	dashboard.WithPanel(panel.Min(0))
}

func addMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory").
		Span(widthHalf).Height(heightTall)
	panel.WithTarget(prometheusQuery(
		promql.Sum(promql.Vector(memMetric).
			Label("namespace", namespaceVariable).
			Label("container", nameVariable).
			LabelMatchRegexp("pod", podVariable).
			LabelNeq("image", "")).By([]string{"pod"}).String()).
		Format("time_series").LegendFormat("{{target_pod}}"))
	builder.WithPanel(panel.Min(0))
}

func addJobCPUPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(promql.Vector(cpuMetric).
				Label("namespace", namespaceVariable).
				Label("container", nameVariable).
				LabelMatchRegexp("pod", podVariable).
				LabelNeq("image", "").
				Range(intervalVariable))).By([]string{"pod"}).String()).
			Format("time_series").LegendFormat("{{pod}}"))
	dashboard.WithPanel(panel.Min(0))
}

func addJobMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Vector(memMetric).
				Label("namespace", namespaceVariable).
				Label("container", nameVariable).
				LabelMatchRegexp("pod", podVariable).
				LabelNeq("image", "")).By([]string{"pod"}).String()).
			Format("time_series").LegendFormat("{{pod}}"))
	builder.WithPanel(panel.Min(0))
}

func addLatencyPanel(
	builder *dashboard.DashboardBuilder,
	subTitle string,
	kubernetesName string,
	targetPod string,
	container string,
	method string,
	match bool,
) {
	panel := timeseries.NewPanelBuilder().
		Title(fmt.Sprintf("Latency (%s%s)", subTitle, targetPod)).
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(quantile, promql.Sum(promql.Rate(
				addGRPCMatch(
					addBasicLabel(promql.Vector(serverLatencyBucket)).
						LabelMatchRegexp("container", container),
					method, match,
				).Range(intervalVariable))).By([]string{"le", grpcServerMethod})).String(),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{grpc_server_method}} p%d", int(quantile*100)))).Min(0)
	}
	builder.WithPanel(panel)
}

func addCompletedRPCPanel(
	dashboard *dashboard.DashboardBuilder, subTitle string, kubernetesName string, targetPod string,
) {
	panel := timeseries.NewPanelBuilder().
		Title(fmt.Sprintf("Completed RPCs (%s%s)", subTitle, targetPod)).
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(
				addBasicLabel(promql.Vector(serverCompletedRPCs)).
					Range(intervalVariable))).By([]string{grpcServerMethod, grpcServerStatus}).String(),
		).Format("time_series").LegendFormat("{{grpc_server_method}} ({{grpc_server_status}})")).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(
				addBasicLabel(promql.Vector(serverCompletedRPCs)).
					Range(intervalVariable))).By([]string{grpcServerStatus}).String(),
		).Format("time_series").LegendFormat("Total ({{grpc_server_status}})"))
	dashboard.WithPanel(panel)
}

func addGoroutinePanel(
	dashboard *dashboard.DashboardBuilder, kubernetesName string, targetPod string,
) {
	panel := timeseries.NewPanelBuilder().
		Title("Goroutine Count").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			addBasicLabel(promql.Vector("goroutine_count")).String(),
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addGCPanel(dashboard *dashboard.DashboardBuilder, kubernetesName string) {
	panel := timeseries.NewPanelBuilder().
		Title("GC Count").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Increase(promql.Vector("gc_count").
				Label(namespaceKey, namespaceVariable).
				LabelMatchRegexp(nameKey, kubernetesName).
				LabelMatchRegexp("target_node", ".+").
				Range(intervalVariable)).String(),
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Total Indices").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(
				promql.Vector(public_const.IndexCountMetricsName))).String(),
		).Format("time_series").LegendFormat("indices"))
	dashboard.WithPanel(panel)
}

func addUncommitedIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Uncommitted Indices").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(promql.Vector(public_const.UncommittedIndexCountMetricsName))).String(),
		).Format("time_series").LegendFormat("uncommitted-indices"))
	dashboard.WithPanel(panel)
}

func addIndexLatencyPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("SaveIndex Latency").
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(quantile, promql.Sum(promql.Rate(addBasicLabel(
				promql.Vector(serverLatencyBucket),
			).LabelMatchRegexp(grpcServerMethod, ".*Index$").
				Range(intervalVariable))).By([]string{"le", grpcServerMethod})).String(),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{grpc_server_method}} p%d", int(quantile*100))))
	}
	dashboard.WithPanel(panel)
}

func addIndexPerPodPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Indices Per Pod").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(promql.Vector(
				public_const.IndexCountMetricsName))).
				By([]string{"target_pod"}).String(),
		).Format("time_series").LegendFormat("{{hostname}}"))
	dashboard.WithPanel(panel)
}

func addMemstatsPanels(dashboard *dashboard.DashboardBuilder) {
	addMetricPanel(dashboard, "Alloc", "alloc_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "Total Alloc", "alloc_bytes_total", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "Sys", "sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "Lookups", "lookups_count", nil)
	addMetricPanel(dashboard, "Mallocs", "mallocs_total", nil)
	addMetricPanel(dashboard, "Frees", "frees_total", nil)
	addMetricPanel(dashboard, "HeapAlloc", "heap_alloc_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "HeapSys", "heap_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "HeapIdle", "heap_idle_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "HeapInUse", "heap_inuse_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "HeapReleased", "heap_released_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "HeapObjects", "heap_objects_count", nil)
	addMetricPanel(dashboard, "StackInUse", "stack_inuse_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "StackSys", "stack_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "MSpanInUse", "mspan_inuse_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "MSpanSys", "mspan_sys_bytes", cog.ToPtr("decbytes"))

	addMetricPanel(dashboard, "MCacheInUse", "mcache_inuse_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "MCacheSys", "mcache_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "BuckHashSys", "buckhash_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "GCSys", "gc_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "OtherSys", "other_sys_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "NextGC", "next_gc_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "PauseTotalMS", "pause_ms_total", cog.ToPtr("ms"))
	addMetricPanel(dashboard, "NumGC", "gc_count", nil)
	addMetricPanel(dashboard, "NumForcedGC", "forced_gc_count", nil)
	addMetricPanel(dashboard, "HeapWillReturn", "heap_will_return_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "LiveObjects", "live_objects_count", nil)
}

func addProcStatusPanels(dashboard *dashboard.DashboardBuilder) {
	addMetricPanel(dashboard, "VMPeak", "vmpeak_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMSize", "vmsize_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMData", "vmdata_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMRSS", "vmrss_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMHWM", "vmhwm_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMStk", "vmstk_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMSwap", "vmswap_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMExe", "vmexe_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMLib", "vmlib_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMLck", "vmlck_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMPin", "vmpin_bytes", cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, "VMPTE", "vmpte_bytes", cog.ToPtr("decbytes"))
}

func addMetricPanel(
	dashboard *dashboard.DashboardBuilder, title string, metric string, unit *string,
) {
	panel := timeseries.NewPanelBuilder().
		Title(title).
		Span(widthQuarter).Height(heightTall).
		WithTarget(prometheusQuery(
			addBasicLabel(promql.Vector(metric)).String(),
		).Format("time_series").LegendFormat("{{target_pod}}")).Min(0)
	if unit != nil {
		panel.Unit(*unit)
	}
	dashboard.WithPanel(panel)
}
