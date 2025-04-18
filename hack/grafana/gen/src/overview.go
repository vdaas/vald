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
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func addOverviewIndexPanel(builder *dashboard.DashboardBuilder) {
	builder.
		WithPanel(stat.NewPanelBuilder().
			Title("Indices").
			WithTarget(prometheusQuery(
				`sum(agent_core_ngt_index_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdAgentPodName"})`,
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(indiceThresholds),
			).
			Span(widthOneSixth).Height(heightMedium))
}

func addNodeCPUPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Node CPU").
		Span((widthFull - widthOneSixth) / 2).Height(heightMedium).
		WithTarget(prometheusQuery(
			`sum by (instance) (irate(node_cpu_seconds_total{mode!="idle"}[$interval]))`,
		).Format("time_series").LegendFormat("{{instance}}")).
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func addNodeMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Node Memory").
		Span((widthFull - widthOneSixth) / 2).Height(heightMedium).
		WithTarget(prometheusQuery(
			`node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes`,
		).Format("time_series").LegendFormat("{{instance}}")).
		Unit("decbytes").
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func addBackoffPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Backoff Retry Count (Vald LB Gateway)").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			`increase(sum(label_replace(backoff_retry_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdGatewayPodName", backoff_name!~"github.com/.*"}, "rpc", "$1",  "backoff_name", "(.*/.*)/.*")) by (rpc) [$interval:])`,
		).Format("time_series").LegendFormat("{{rpc}}")).
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func addBackoffPerRPCPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Backoff Retry Count / Agent Completed RPCs (Vald LB Gateway)").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			`increase(sum(label_replace(backoff_retry_count{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdGatewayPodName", backoff_name!~"github.com/.*"}, "grpc_server_method", "$1",  "backoff_name", "(.*/.*)/.*")) by (grpc_server_method)[$interval:]) / sum(irate(server_completed_rpcs{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdAgentPodName"}[$interval])) by (grpc_server_method)`,
		).Format("time_series").LegendFormat("{{grpc_server_method}} ({{grpc_server_status}})")).
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func addCircuitBreakerState(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Circuit Breaker State (Vald LB Gateway)").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			`increase(sum(label_replace(circuit_breaker_state{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdGatewayPodName", state="open"}, "rpc", "$1",  "name", "(.*/.*)/.*")) by (rpc, state) [$interval:])`,
		).Format("time_series").LegendFormat("{{rpc}} ({{state}})")).
		WithTarget(prometheusQuery(
			`increase(sum(label_replace(circuit_breaker_state{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$ValdGatewayPodName", state="half-open"}, "rpc", "$1",  "name", "(.*/.*)/.*")) by (rpc, state) [$interval:])`,
		).Format("time_series").LegendFormat("{{rpc}} ({{state}})")).
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func repeatOverview(builder *dashboard.DashboardBuilder) {
	podStatusPanel := timeseries.NewPanelBuilder().
		Title("Pod Status").
		Span(uint32(widthFull - witdhOneThird*2 - widthOneEighth)).Height(heightTall)
	for _, ks := range allKindStatus {
		podStatusPanel.
			WithTarget(prometheusQuery(
				fmt.Sprintf(
					`max(kube_%s_status_%s{namespace="$Namespace", %s="$ReplicaSet"}) and on() count(kube_%s_created{%s=~"$ReplicaSet"}) >= 1`,
					ks.kind,
					ks.status,
					ks.kind,
					ks.kind,
					ks.kind,
				),
			).Format("time_series").LegendFormat(ks.status)).
			FillOpacity(opacity)
	}
	builder.
		WithRow(dashboard.NewRowBuilder("$ReplicaSet").Repeat("ReplicaSet")).
		WithPanel(stat.NewPanelBuilder().
			Title("Pods").
			WithTarget(prometheusQuery(
				`count(kube_pod_info{namespace="$Namespace", pod=~"$ReplicaSet.*"})`,
			).Format("table")).
			Span(widthOneEighth).Height(heightTall)).
		WithPanel(podStatusPanel).
		WithPanel(timeseries.NewPanelBuilder().
			Title("CPU").
			WithTarget(prometheusQuery(
				`sum(irate(container_cpu_usage_seconds_total{namespace="$Namespace", pod=~"$ReplicaSet.*", image!=""}[$interval])) by (pod)`,
			).Format("time_series")).
			Span(witdhOneThird).Height(heightTall).Min(0)).
		WithPanel(timeseries.NewPanelBuilder().
			Title("Memory Working Set").
			WithTarget(prometheusQuery(
				`sum(container_memory_working_set_bytes{namespace="$Namespace", pod=~"$ReplicaSet.*", image!=""}) by (pod)`,
			).Format("time_series")).
			Unit("decbytes").
			Span(witdhOneThird).Height(heightTall).Min(0))
}
