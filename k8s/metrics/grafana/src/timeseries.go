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
)

func addCPUPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Span(12).Height(8)
	for _, resource := range []string{"statefulset", "deployment", "daemonset"} {
		panel.WithTarget(prometheusQuery(
			fmt.Sprintf(
				`sum(irate(container_cpu_usage_seconds_total{namespace="$Namespace", container="$ReplicaSet", pod=~"$PodName", image!=""}[$interval])) by (pod) and on() count(kube_%s_created{%s="$ReplicaSet"}) >= 1`,
				resource,
				resource,
			),
		).
			Format("time_series").LegendFormat("{{pod}}"))
	}
	dashboard.WithPanel(panel)
}

func addMemoryPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory").
		Span(12).Height(8)
	for _, resource := range []string{"statefulset", "deployment", "daemonset"} {
		panel.WithTarget(prometheusQuery(
			fmt.Sprintf(
				`sum(container_memory_working_set_bytes{namespace="$Namespace", container="$ReplicaSet", pod=~"$PodName", image!=""}) by (pod) and on() count(kube_%s_created{%s="$ReplicaSet"}) >= 1`,
				resource,
				resource,
			),
		).
			Format("time_series").LegendFormat("{{target_pod}}"))
	}
	dashboard.WithPanel(panel)
}

func addLatencyPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Latency").
		Span(12).Height(8)
	for _, quantile := range []float32{0.5, 0.95, 0.99} {
		panel.WithTarget(prometheusQuery(
			fmt.Sprintf(
				`histogram_quantile(%f, sum(rate(server_latency_bucket{exported_kubernetes_namespace="$Namespace", kubernetes_name="$ReplicaSet", target_pod=~"$PodName"}[$interval])) by (le, grpc_server_method))`,
				quantile,
			),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{grpc_server_method}} p%f", quantile)))
	}
	dashboard.WithPanel(panel)
}

func addCompletedRPCPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Completed RPCs").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`sum(irate(server_completed_rpcs{exported_kubernetes_namespace="$Namespace", kubernetes_name="$ReplicaSet", target_pod=~"$PodName"}[$interval])) by (grpc_server_method, grpc_server_status)`,
		).Format("time_series").LegendFormat("{{grpc_server_method}} ({{grpc_server_status}})")).
		WithTarget(prometheusQuery(
			`sum(irate(server_completed_rpcs{exported_kubernetes_namespace="$Namespace", kubernetes_name="$ReplicaSet", target_pod=~"$PodName"}[$interval])) by (grpc_server_status)`,
		).Format("time_series").LegendFormat("Total ({{grpc_server_status}})"))
	dashboard.WithPanel(panel)
}

func addGoroutinePanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Goroutine Count").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`goroutine_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"}`,
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addGCPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("GC Count").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`increase(gc_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_node=~".+"}[$interval])`,
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Total Indices").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`sum(agent_core_ngt_index_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"})`,
		).Format("time_series").LegendFormat("indices"))
	dashboard.WithPanel(panel)
}

func addUncommitedIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Uncommitted Indices").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`sum(agent_core_ngt_uncommitted_index_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"})`,
		).Format("time_series").LegendFormat("uncommitted-indices"))
	dashboard.WithPanel(panel)
}

func addIndexLatencyPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("SaveIndex Latency").
		Span(12).Height(8)
	for _, quantile := range []float32{0.5, 0.95, 0.99} {
		panel.WithTarget(prometheusQuery(
			fmt.Sprintf(
				`histogram_quantile(%f, sum(rate(server_latency_bucket{exported_kubernetes_namespace="$Namespace", kubernetes_name="$ReplicaSet", target_pod=~"$PodName", grpc_server_method=~".*Index$"}[$interval])) by (le, grpc_server_method))`,
				quantile,
			),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{grpc_server_method}} p%f", quantile)))
	}
	dashboard.WithPanel(panel)
}

func addIndexPerPodPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Indices Per Pod").
		Span(12).Height(8).
		WithTarget(prometheusQuery(
			`sum(agent_core_ngt_index_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"}) by (target_pod)`,
		).Format("time_series").LegendFormat("{{hostname}}"))
	dashboard.WithPanel(panel)
}

func addMemstatsPanels(dashboard *dashboard.DashboardBuilder) {
	addOverviewPanel(dashboard, "Alloc", "alloc_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Total Alloc", "alloc_bytes_total", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Sys", "sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Lookups", "lookups_count", nil)
	addOverviewPanel(dashboard, "Mallocs", "mallocs_total", nil)
	addOverviewPanel(dashboard, "Frees", "frees_total", nil)
	addOverviewPanel(dashboard, "HeapAlloc", "heap_alloc_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "HeapSys", "heap_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "HeapIdle", "heap_idle_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "HeapInUse", "heap_inuse_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "HeapReleased", "heap_released_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "HeapObjects", "heap_objects_count", nil)
	addOverviewPanel(dashboard, "StackInUse", "stack_inuse_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "StackSys", "stack_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "MSpanInUse", "mspan_inuse_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "MSpanSys", "mspan_sys_bytes", cog.ToPtr("decbytes"))

	addOverviewPanel(dashboard, "MCacheInUse", "mcache_inuse_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "MCacheSys", "mcache_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "BuckHashSys", "buckhash_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "GCSys", "gc_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "OtherSys", "other_sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "NextGC", "next_gc_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "PauseTotalMS", "pause_ms_total", cog.ToPtr("ms"))
	addOverviewPanel(dashboard, "NumGC", "gc_count", nil)
	addOverviewPanel(dashboard, "NumForcedGC", "forced_gc_count", nil)
	addOverviewPanel(dashboard, "HeapWillReturn", "heap_will_return_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "LiveObjects", "live_objects_count", nil)
}

func addProcStatusPanels(dashboard *dashboard.DashboardBuilder) {
	addOverviewPanel(dashboard, "VMPeak", "vmpeak_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMSize", "vmsize_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMData", "vmdata_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMRSS", "vmrss_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMHWM", "vmhwm_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMStk", "vmstk_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMSwap", "vmswap_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMExe", "vmexe_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMLib", "vmlib_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMLck", "vmlck_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMPin", "vmpin_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "VMPTE", "vmpte_bytes", cog.ToPtr("decbytes"))
}

func addOverviewPanel(dashboard *dashboard.DashboardBuilder, title string, metric string, unit *string) {
	panel := timeseries.NewPanelBuilder().
		Title(title).
		Span(6).Height(8).
		WithTarget(prometheusQuery(
			fmt.Sprintf(`%s{exported_kubernetes_namespace="$Namespace", target_pod=~"$PodName", kubernetes_name=~"$ReplicaSet"}`, metric),
		).Format("time_series").LegendFormat("{{target_pod}}")).Min(0)
	if unit != nil {
		panel.Unit(*unit)
	}
	dashboard.WithPanel(panel)
}
