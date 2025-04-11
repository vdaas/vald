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

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
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

func addAgentMemoryPanels(dashboard *dashboard.DashboardBuilder) {
	addOverviewPanel(dashboard, "Alloc", "alloc_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Total Alloc", "alloc_bytes_total", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Sys", "sys_bytes", cog.ToPtr("decbytes"))
	addOverviewPanel(dashboard, "Lookups", "lookups_count", cog.ToPtr("decbytes"))
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
