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
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/promql-builder/go/promql"
	"github.com/vdaas/vald/internal/observability/metrics/agent/core/ngt"
)

func addOverviewIndexPanel(builder *dashboard.DashboardBuilder) {
	builder.
		WithPanel(stat.NewPanelBuilder().
			Title("Indices").
			WithTarget(prometheusQuery(
				promql.Sum(promql.Vector(ngt.IndexCountMetricsName).
					Label(namespaceKey, namespaceVariable).
					LabelMatchRegexp(nameKey, nameVariable).
					LabelMatchRegexp(podKey, "$ValdAgentPodName")).String(),
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
			promql.Sum(promql.Irate(promql.Vector(nodeCPUMetric).
				LabelNeq("mode", "idle").
				Range(intervalVariable))).
				By([]string{instanceKey}).String(),
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
			promql.Increase(promql.Subquery(promql.Sum(promql.LabelReplace(promql.Vector(backoffRetryCount).
				Label(namespaceKey, namespaceVariable).
				LabelMatchRegexp(nameKey, nameVariable).
				LabelMatchRegexp(podKey, "$ValdGatewayPodName").
				LabelNotMatchRegexp("backoff_name", "github.com/.*"),
				"\"rpc\"", "\"$1\"", "\"backoff_name\"", "\"(.*/.*)/.*\"")).
				By([]string{"rpc"})).Range(intervalVariable + ":")).String(),
		).Format("time_series").LegendFormat("{{rpc}}")).
		FillOpacity(opacity)
	builder.WithPanel(panel)
}

func addBackoffPerRPCPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Backoff Retry Count / Agent Completed RPCs (Vald LB Gateway)").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Div(
				promql.Increase(promql.Subquery(promql.Sum(promql.LabelReplace(promql.Vector(backoffRetryCount).
					Label(namespaceKey, namespaceVariable).
					LabelMatchRegexp(nameKey, nameVariable).
					LabelMatchRegexp(podKey, "$ValdGatewayPodName").
					LabelNotMatchRegexp("backoff_name", "github.com/.*"),
					"\"rpc\"", "\"$1\"", "\"backoff_name\"", "\"(.*/.*)/.*\"")).
					By([]string{"rpc"})).Range(intervalVariable+":")),
				promql.Sum(promql.Irate(promql.Vector(serverCompletedRPCs).
					Label(namespaceKey, namespaceVariable).
					LabelMatchRegexp(nameKey, nameVariable).
					LabelMatchRegexp(podKey, "$ValdAgentPodName").
					Range(intervalVariable))).By([]string{"grpc_server_method"})).String(),
		).Format("time_series").LegendFormat("{{grpc_server_method}} ({{grpc_server_status}})"),
		).FillOpacity(opacity)
	builder.WithPanel(panel)
}

func circuitBreakerQuery(state string) string {
	return promql.Increase(promql.Subquery(promql.Sum(promql.LabelReplace(
		promql.Vector(circuitBreakerState).
			Label(namespaceKey, namespaceVariable).
			LabelMatchRegexp(nameKey, nameVariable).
			LabelMatchRegexp(podKey, "$ValdGatewayPodName").
			Label("state", state),
		"\"rpc\"", "\"$1\"", "\"name\"", "\"(.*/.*)/.*\"")).
		By([]string{"rpc", "state"})).Range(intervalVariable + ":")).String()
}

func addCircuitBreakerState(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Circuit Breaker State (Vald LB Gateway)").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			circuitBreakerQuery("open"),
		).Format("time_series").LegendFormat("{{rpc}} ({{state}})")).
		WithTarget(prometheusQuery(
			circuitBreakerQuery("half-open"),
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
				promql.Max(promql.Vector("kube_"+ks.kind+"_status_"+ks.status).
					Label("namespace", namespaceVariable).
					Label(ks.kind, "$ReplicaSet"),
				).String(),
			).Format("time_series").LegendFormat(ks.status)).
			FillOpacity(opacity)
	}
	builder.
		WithRow(dashboard.NewRowBuilder("$ReplicaSet").Repeat("ReplicaSet")).
		WithPanel(stat.NewPanelBuilder().
			Title("Pods").
			WithTarget(prometheusQuery(
				promql.Count(promql.Vector(podInfo).
					Label("namespace", namespaceVariable).
					LabelMatchRegexp("pod", "$ReplicaSet.*")).String(),
			).Format("table")).
			Span(widthOneEighth).Height(heightTall)).
		WithPanel(podStatusPanel).
		WithPanel(timeseries.NewPanelBuilder().
			Title("CPU").
			WithTarget(prometheusQuery(
				promql.Sum(promql.Irate(promql.Vector(cpuMetric).
					Label("namespace", namespaceVariable).
					LabelMatchRegexp("pod", "$ReplicaSet.*").
					LabelNeq("image", "").Range(intervalVariable))).By([]string{"pod"}).String(),
			).Format("time_series")).
			Span(witdhOneThird).Height(heightTall).Min(0)).
		WithPanel(timeseries.NewPanelBuilder().
			Title("Memory Working Set").
			WithTarget(prometheusQuery(
				promql.Subquery(promql.Sum(promql.Vector(memMetric).
					Label("namespace", namespaceVariable).
					LabelMatchRegexp("pod", "$ReplicaSet.*").
					LabelNeq("image", "")).By([]string{"pod"})).String(),
			).Format("time_series")).
			Unit("decbytes").
			Span(witdhOneThird).Height(heightTall).Min(0))
}
